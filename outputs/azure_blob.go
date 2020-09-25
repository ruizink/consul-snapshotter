package outputs

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/ruizink/consul-snapshotter/azure"
)

type AzureBlobOutput struct {
	ContainerName    string
	ContainerPath    string
	Filename         string
	StorageAccount   string
	StorageAccessKey string
	StorageSASToken  string
	RetentionPeriod  time.Duration
}

func (o *AzureBlobOutput) Save(snap string) {
	destFile := path.Join(o.ContainerPath, o.Filename)
	config, err := azure.AzureConfig(o.StorageAccount, o.StorageAccessKey, o.StorageSASToken)
	if err != nil {
		log.Println("Invalid Azure config:", err)
		return
	}
	_, err = azure.UploadBlob(snap, destFile, o.ContainerName, config)
	if err != nil {
		log.Println("Error uploading snapshot file:", err)
		return
	}
	log.Println("Uploaded snapshot to:", destFile)
}

func (o *AzureBlobOutput) ApplyRetentionPolicy() error {
	log.Println(fmt.Sprintf("Azure Blob Storage retention: %v", o.RetentionPeriod))
	if o.RetentionPeriod <= 0 {
		return nil
	}

	log.Println(fmt.Sprintf("Applying retention policy (remove files older than %v) in Azure Blob Storage", o.RetentionPeriod))

	config, err := azure.AzureConfig(o.StorageAccount, o.StorageAccessKey, o.StorageSASToken)
	if err != nil {
		return err
	}

	blobList, err := azure.ListBlobs(o.ContainerName, config)

	for _, blob := range blobList {
		if time.Now().Sub(blob.Properties.LastModified) <= o.RetentionPeriod {
			continue
		}

		log.Println("Removing from Azure Blob Storage: " + blob.Name)
		err := azure.DeleteBlob(o.ContainerName, blob, config)
		if err != nil {
			return err
		}
	}

	return nil
}
