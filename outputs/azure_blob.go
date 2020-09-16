package outputs

import (
	"log"
	"path"

	"github.com/ruizink/consul-snapshot/azure"
)

type AzureBlobOutput struct {
	ContainerName    string
	ContainerPath    string
	Filename         string
	StorageAccount   string
	StoraceAccessKey string
}

func (o *AzureBlobOutput) Save(snap string) {
	destFile := path.Join(o.ContainerPath, o.Filename)
	config, err := azure.AzureConfig(o.StorageAccount, o.StoraceAccessKey)
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
