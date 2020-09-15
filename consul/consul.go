package consul

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/snapshot"
	"github.com/rboyer/safeio"
)

type worker struct {
	client         *api.Client
	key            string
	SessionID      string
	sessionTimeout string
}

func GetSnapshot(w *worker) (string, error) {

	var buf bytes.Buffer

	// Take the snapshot
	snap, metadata, err := w.client.Snapshot().Save(&api.QueryOptions{})
	if err != nil {
		return "", fmt.Errorf("Error requesting the snapshot: %v", err)
	}
	defer snap.Close()
	log.Println(fmt.Sprintf("Performed snapshot (up to index=%d)", metadata.LastIndex))

	tee := io.TeeReader(snap, &buf)

	// Verify the snapshot
	if _, err := snapshot.Verify(tee); err != nil {
		return "", fmt.Errorf("Error verifying snapshot: %v", err)
	}

	// Save the verified snapshot to a temporary location
	snapFile, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf("Error creating temp file: %v", err)
	}
	snapFileName := snapFile.Name()
	log.Println("Saving snapshot to temporary file: ", snapFileName)

	if _, err := safeio.WriteToFile(&buf, snapFileName, 0644); err != nil {
		return "", fmt.Errorf("Error writing snapshot file: %v", err)
	}

	return snapFileName, nil
}

func NewWorker(consulURL, consulToken, key string, sessionTimeout time.Duration) (*worker, error) {

	// Create the HTTP client
	conf := api.DefaultConfig()
	conf.Address = consulURL
	conf.Token = consulToken
	client, err := api.NewClient(conf)
	if err != nil {
		log.Println("Could not create a new client:", err)
		return nil, err
	}

	// create new session for this worker
	w := &worker{
		client:         client,
		key:            key,
		sessionTimeout: sessionTimeout.String(),
	}

	return w, nil
}

func AcquireLock(w *worker) error {
	if err := w.createSession(); err != nil {
		return err
	}

	r, err := w.acquireLock()
	if err != nil {
		return err
	}
	if !r {
		return fmt.Errorf("Lock is acquired by another resource")
	}
	return nil
}

func ReleaseLock(w *worker) error {
	_, err := w.releaseLock()
	if err != nil {
		return err
	}
	if err := w.destroySession(); err != nil {
		return err
	}
	return nil
}

func (w *worker) RenewSession(doneChan <-chan struct{}) error {
	err := w.client.Session().RenewPeriodic(w.sessionTimeout, w.SessionID, nil, doneChan)
	if err != nil {
		return err
	}
	return nil
}

func (w *worker) createSession() error {
	sessionConf := &api.SessionEntry{
		TTL:      w.sessionTimeout,
		Behavior: "delete",
	}

	sessionID, _, err := w.client.Session().Create(sessionConf, nil)
	if err != nil {
		return err
	}

	w.SessionID = sessionID
	return nil
}

func (w *worker) acquireLock() (bool, error) {
	KVPair := &api.KVPair{
		Key:     w.key,
		Value:   []byte(w.SessionID),
		Session: w.SessionID,
	}

	acquired, _, err := w.client.KV().Acquire(KVPair, nil)
	return acquired, err
}

func (w *worker) releaseLock() (bool, error) {
	KVPair := &api.KVPair{
		Key:     w.key,
		Value:   []byte(w.SessionID),
		Session: w.SessionID,
	}

	released, _, err := w.client.KV().Release(KVPair, nil)
	return released, err
}

func (w *worker) destroySession() error {
	_, err := w.client.Session().Destroy(w.SessionID, nil)
	if err != nil {
		return fmt.Errorf("Could not destroy session: %v", err)
	}

	return nil
}
