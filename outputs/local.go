package outputs

import (
	"log"
	"os"
	"path"
	"time"
)

type LocalOutput struct {
	DestinationPath          string
	Filename                 string
	RetentionKeep            int
	RetentionDeleteOlderThan time.Duration
}

func (o *LocalOutput) Save(snap string) {
	dstFile := path.Join(o.DestinationPath, o.Filename)
	if err := os.Rename(snap, dstFile); err != nil {
		log.Println("Error writing snapshot file: ", err)
		return
	}
	log.Println("Saved snapshot to:", dstFile)
}
