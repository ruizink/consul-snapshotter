package azure

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"

	"github.com/ruizink/consul-snapshotter/logger"
)

type AzureConfig struct {
	ContainerName    string
	ContainerPath    string
	Filename         string
	StorageAccount   string
	CloudDomain      string
	StorageAccessKey string
	StorageSASToken  string
	CreateContainer  bool
	BlockSize        int64
	Parallelism      uint16
	Emulated         bool
	EmulatorUrl      string
}

type Azure struct {
	client *azblob.Client
	config *AzureConfig
}

func NewAzure(config *AzureConfig) (*Azure, error) {
	var (
		azclient *azblob.Client
		azURL    *url.URL
		err      error
	)

	if config.StorageAccount == "" {
		return nil, fmt.Errorf("Azure Account Name not provided")
	}
	if config.StorageAccessKey == "" && config.StorageSASToken == "" {
		return nil, fmt.Errorf("Azure Account Access Key or SAS Token must be provided")
	}

	// create azure client
	if config.StorageSASToken != "" {
		if config.Emulated {
			azURL, _ = url.Parse(fmt.Sprintf("%s/%s/?%s", config.EmulatorUrl, config.StorageAccount, config.StorageSASToken))
		} else {
			azURL, _ = url.Parse(fmt.Sprintf("https://%s.%s/?%s", config.StorageAccount, config.CloudDomain, config.StorageSASToken))
		}
		logger.Debug("Using Azure Blob URL: ", azURL.String())
		azclient, err = azblob.NewClientWithNoCredential(azURL.String(), nil)
	} else {
		if config.Emulated {
			azURL, _ = url.Parse(fmt.Sprintf("%s/%s/", config.EmulatorUrl, config.StorageAccount))
		} else {
			azURL, _ = url.Parse(fmt.Sprintf("https://%s.%s/", config.StorageAccount, config.CloudDomain))
		}
		cred, cerr := azblob.NewSharedKeyCredential(config.StorageAccount, config.StorageAccessKey)
		if cerr != nil {
			return nil, cerr
		}
		logger.Debug("Using Azure Blob URL: ", azURL.String())
		azclient, err = azblob.NewClientWithSharedKeyCredential(azURL.String(), cred, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating azure client: %s", err)
	}

	return &Azure{client: azclient, config: config}, nil
}

func (az *Azure) ListBlobsOlderThan(period time.Duration) ([]*container.BlobItem, error) {
	var results = make([]*container.BlobItem, 0)

	// blob listings are returned across multiple pages
	pager := az.client.NewListBlobsFlatPager(az.config.ContainerName, nil)

	// continue fetching pages until no more remain
	for pager.More() {
		// advance to the next page
		logger.Debug("Getting next page of blobs...")
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, blob := range page.Segment.BlobItems {
			if time.Since(*blob.Properties.LastModified) > period {
				results = append(results, blob)
			}
		}
	}

	return results, nil
}

func (az *Azure) DeleteBlob(blob *container.BlobItem) error {
	logger.Debug("Deleting blob: ", *blob.Name)
	_, err := az.client.DeleteBlob(context.Background(), az.config.ContainerName, *blob.Name, nil)
	if err != nil {
		return err
	}

	return nil
}

func (az *Azure) UploadBlob(srcFile string) error {
	// Create the container if it doesn't exist
	if az.config.CreateContainer {
		logger.Debug("Creating container: ", az.config.ContainerName)
		// _, err := azclient.CreateContainer(context.Background(), az.ContainerName, &azblob.CreateContainerOptions{
		// 	Access: azblob.PublicAccessNone,
		// })
		_, err := az.client.CreateContainer(context.Background(), az.config.ContainerName, nil)
		az.client.URL()

		var respErr *azcore.ResponseError
		if err != nil {
			if !(errors.As(err, &respErr) && respErr.ErrorCode == "ContainerAlreadyExists") {
				return fmt.Errorf("error creating container: %s", err)
			} else {
				logger.Debug("Got ContainerAlreadyExists, ignoring...")
			}
		}
	}

	// Upload the blob
	logger.Info(fmt.Sprintf("Uploading the file (BlockSize: %v, Parallelism: %v)", az.config.BlockSize, az.config.Parallelism))

	file, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}

	destFile := path.Join(az.config.ContainerPath, az.config.Filename)

	_, err = az.client.UploadFile(context.Background(), az.config.ContainerName, destFile, file, &azblob.UploadFileOptions{
		BlockSize:   az.config.BlockSize,
		Concurrency: az.config.Parallelism,
	})
	if err != nil {
		return fmt.Errorf("error uploading file: %s", err)
	}

	return nil
}
