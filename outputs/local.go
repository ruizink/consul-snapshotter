package outputs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type LocalOutput struct {
	DestinationPath string
	Filename        string
	RetentionKeep   int
	RetentionPeriod time.Duration
}

func (o *LocalOutput) Save(snap string) error {
	dstFile := path.Join(o.DestinationPath, o.Filename)
	if err := os.Rename(snap, dstFile); err != nil {
		return err
	}
	log.Println("Saved snapshot to:", dstFile)
	return nil
}

func (o *LocalOutput) ApplyRetentionPolicy() error {
	if o.RetentionPeriod > 0 {
		log.Println(fmt.Sprintf("Applying retention policy (remove files older than %v) in: %s", o.RetentionPeriod, o.DestinationPath))
		files, err := findFilesOlderThan(o.DestinationPath, o.RetentionPeriod)
		if err != nil {
			return err
		}

		if len(files) > 0 {
			log.Println("List of files to remove:", strings.Join(files, ", "))
			for _, file := range files {
				os.Remove(file)
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
