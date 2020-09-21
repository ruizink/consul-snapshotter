package azure

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
)

type config struct {
	accountName string
	accountKey  string
	sasToken    string
}

func AzureConfig(accountName, accountKey, sasToken string) (*config, error) {
	if accountName == "" {
		return nil, fmt.Errorf("Azure Account Name not provided")
	}
	if accountKey == "" && sasToken == "" {
		return nil, fmt.Errorf("Azure Account Access Key or SAS Token must be provided")
	}
	c := &config{
		accountName: accountName,
		accountKey:  accountKey,
		sasToken:    sasToken,
	}
	return c, nil
}

func AuthenticateAccountKey(containerName string, c *config) (pipeline.Pipeline, error) {
	// Create a default request pipeline using your storage account name and account key
	credential, err := azblob.NewSharedKeyCredential(c.accountName, c.accountKey)
	if err != nil {
		return nil, fmt.Errorf("Invalid credentials: %s", err)
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	return p, err
}

func AuthenticateSASToken(containerName string, c *config) pipeline.Pipeline {
	return azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})
}

func UploadBlob(srcFile, destFile, containerName string, c *config) (int, error) {
	var p pipeline.Pipeline
	var queryParameters string

	if c.sasToken != "" {
		p = AuthenticateSASToken(containerName, c)
		queryParameters = "?" + c.sasToken
	} else {
		var err error
		p, err = AuthenticateAccountKey(containerName, c)
		if err != nil {
			return 0, err
		}
	}

	// TODO: Allow the URL to be a parameter
	// Setup the blob service URL endpoint
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s%s", c.accountName, containerName, queryParameters))

	// Create a ContainerURL object using the container URL and a request pipeline
	containerURL := azblob.NewContainerURL(*URL, p)

	ctx := context.Background()

	// Create a BlobURL object using the ContainerURL
	blobURL := containerURL.NewBlockBlobURL(destFile)
	file, err := os.Open(srcFile)
	if err != nil {
		return 0, fmt.Errorf("Error opening file: %s", err)
	}

	// TODO: Allow the Parallelism to be a parameter
	// Upload the blob
	uBlockSize := int64(4 * 1024 * 1024)
	uParallelism := uint16(16)
	log.Println(fmt.Sprintf("Uploading the file (BlockSize: %v, Parallelism: %v)", uBlockSize, uParallelism))
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   uBlockSize,
		Parallelism: uParallelism})
	if err != nil {
		return 0, fmt.Errorf("Error uploading file: %s", err)
	}

	return 1, nil
}
