package outputs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ruizink/consul-snapshotter/logger"
)

type LocalOutput struct {
	DestinationPath string
	Filename        string
	RetentionPeriod time.Duration
}

func (o *LocalOutput) Save(snap string) error {
	dstFile := path.Join(o.DestinationPath, o.Filename)
	if err := os.Rename(snap, dstFile); err != nil {
		return err
	}

	logger.Info("Saved snapshot to: ", dstFile)
	return nil
}

func (o *LocalOutput) ApplyRetentionPolicy() error {
	var errors error

	if o.RetentionPeriod <= 0 {
		return nil
	}

	logger.Info(fmt.Sprintf("Applying local retention policy (remove files older than %v)", o.RetentionPeriod))
	files, err := findFilesOlderThan(o.DestinationPath, o.RetentionPeriod)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		logger.Info("List of files to remove: ")
		for _, file := range files {
			logger.Info(file)
			if err := os.Remove(file); err != nil {
				errors = multierror.Append(errors, err)
			}
		}
	}
	return nil
}

func findFilesOlderThan(dir string, period time.Duration) (fileList []string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if time.Now().Sub(file.ModTime()) > period {
				fileList = append(fileList, path.Join(dir, file.Name()))
			}
		}
	}
	return
}
