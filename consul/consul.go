package consul

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/snapshot"
	"github.com/rboyer/safeio"

	"github.com/ruizink/consul-snapshotter/logger"
)

type Worker struct {
	client         *api.Client
	key            string
	SessionID      string
	sessionTimeout string
}

func NewConsul(consulURL, consulToken, key string, sessionTimeout time.Duration) (*Worker, error) {

	// Create the HTTP client
	conf := api.DefaultConfig()
	conf.Address = consulURL
	conf.Token = consulToken
	client, err := api.NewClient(conf)
	if err != nil {
		return nil, fmt.Errorf("could not create a new client: %v", err)
	}

	// create new session for this Worker
	w := &Worker{
		client:         client,
		key:            key,
		sessionTimeout: sessionTimeout.String(),
	}

	return w, nil
}

func (w *Worker) GetSnapshot() (string, error) {

	var buf bytes.Buffer

	// Take the snapshot
	snap, metadata, err := w.client.Snapshot().Save(&api.QueryOptions{})
	if err != nil {
		return "", fmt.Errorf("error requesting the snapshot: %v", err)
	}
	defer snap.Close()
	logger.Info(fmt.Sprintf("Performed snapshot (up to index=%d)", metadata.LastIndex))

	tee := io.TeeReader(snap, &buf)

	// Verify the snapshot
	if _, err := snapshot.Verify(tee); err != nil {
		return "", fmt.Errorf("error verifying snapshot: %v", err)
	}

	// Save the verified snapshot to a temporary location
	snapFile, err := os.CreateTemp("", "")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %v", err)
	}
	snapFileName := snapFile.Name()
	logger.Debug("Saving snapshot to temporary file: ", snapFileName)

	if _, err := safeio.WriteToFile(&buf, snapFileName, 0644); err != nil {
		return "", fmt.Errorf("error writing snapshot file: %v", err)
	}

	return snapFileName, nil
}

func (w *Worker) AcquireLock() error {
	// create session
	sessionConf := &api.SessionEntry{
		TTL:      w.sessionTimeout,
		Behavior: "delete",
	}

	sessionID, _, err := w.client.Session().Create(sessionConf, nil)
	if err != nil {
		return err
	}

	w.SessionID = sessionID

	// acquire lock
	KVPair := &api.KVPair{
		Key:     w.key,
		Value:   []byte(w.SessionID),
		Session: w.SessionID,
	}

	r, _, err := w.client.KV().Acquire(KVPair, nil)
	if err != nil {
		return err
	}
	if !r {
		return fmt.Errorf("lock is acquired by another resource")
	}
	return nil
}

func (w *Worker) ReleaseLock() error {
	KVPair := &api.KVPair{
		Key:     w.key,
		Value:   []byte(w.SessionID),
		Session: w.SessionID,
	}

	// release lock
	if _, _, err := w.client.KV().Release(KVPair, nil); err != nil {
		return err
	}

	// destroy session
	if _, err := w.client.Session().Destroy(w.SessionID, nil); err != nil {
		return err
	}

	return nil
}

func (w *Worker) RenewSession(doneChan <-chan struct{}) error {
	err := w.client.Session().RenewPeriodic(w.sessionTimeout, w.SessionID, nil, doneChan)
	if err != nil {
		return err
	}
	return nil
}
