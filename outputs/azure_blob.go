package outputs

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/ruizink/consul-snapshotter/azure"
	"github.com/ruizink/consul-snapshotter/logger"
)

type AzureBlobOutput struct {
	AzureConfig     *azure.AzureConfig
	RetentionPeriod time.Duration
}

func (o *AzureBlobOutput) Save(snap string) error {
	az, err := azure.NewAzure(o.AzureConfig)
	if err != nil {
		return fmt.Errorf("invalid azure config: %v", err)
	}
	err = az.UploadBlob(snap)
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

	azure, err := azure.NewAzure(o.AzureConfig)
	if err != nil {
		return err
	}

	blobList, _ := azure.ListBlobsOlderThan(o.RetentionPeriod)

	if len(blobList) > 0 {
		logger.Info("List of Azure Blobs to remove:")
		for _, blob := range blobList {
			logger.Info(*blob.Name)
			if err := azure.DeleteBlob(blob); err != nil {
				errors = multierror.Append(errors, err)
			}
		}
	}

	return errors
}
