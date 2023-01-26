package outputs

import (
	"fmt"
	"path"
	"time"

	"github.com/ruizink/consul-snapshotter/azure"
	"github.com/ruizink/consul-snapshotter/logger"
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
		return fmt.Errorf("invalid azure config: %v", err)
	}
	_, err = azure.UploadBlob(snap, destFile, o.ContainerName, config)
	if err != nil {
		return fmt.Errorf("error uploading snapshot file: %v", err)
	}
	// logger.Info("Uploaded snapshot to: ", destFile)
	return nil
}

func (o *AzureBlobOutput) ApplyRetentionPolicy() error {
	var errors error

	if o.RetentionPeriod <= 0 {
		return nil
	}

	logger.Info(fmt.Sprintf("Applying Azure Blob Storage retention policy (remove blobs older than %v)", o.RetentionPeriod))

	config, err := azure.AzureConfig(o.StorageAccount, o.StorageAccessKey, o.StorageSASToken)
	if err != nil {
		return err
	}

	blobList, err := azure.ListBlobs(o.ContainerName, config)

	for _, blob := range blobList {
		if time.Now().Sub(blob.Properties.LastModified) <= o.RetentionPeriod {
			continue
		}

	if len(blobList) > 0 {
		logger.Info("List of Azure Blobs to remove:")
		for _, blob := range blobList {
			logger.Info(*blob.Name)
			if err := azure.DeleteBlob(blob); err != nil {
				errors = multierror.Append(errors, err)
			}
		}
	}

	return nil
}
