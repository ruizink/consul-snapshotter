package outputs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/ruizink/consul-snapshotter/logger"
)

type LocalOutput struct {
	DestinationPath   string
	Filename          string
	CreateDestination bool
	RetentionPeriod   time.Duration
}

func (o *LocalOutput) Save(snap string) error {
	// create destination dir if it doesn't exist
	if o.CreateDestination {
		if _, err := os.Stat(o.DestinationPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(o.DestinationPath, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	dstFile := path.Join(o.DestinationPath, o.Filename)

	//Read all the contents of the  original file
	bytesRead, err := os.ReadFile(snap)
	if err != nil {
		return err
	}

	//Copy all the contents to the desitination file
	err = os.WriteFile(dstFile, bytesRead, 0644)
	if err != nil {
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

	return errors
}

func findFilesOlderThan(dir string, period time.Duration) (fileList []string, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		file, _ := entry.Info()
		if file.Mode().IsRegular() {
			if time.Since(file.ModTime()) > period {
				fileList = append(fileList, path.Join(dir, file.Name()))
			}
		}
	}
	return fileList, nil
}
