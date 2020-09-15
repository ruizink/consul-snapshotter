package azure

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type config struct {
	accountName string
	accountKey  string
}

func AzureConfig(accountName, accountKey string) (*config, error) {
	if accountName == "" {
		return nil, fmt.Errorf("Azure Account Name not provided")
	}
	if accountKey == "" {
		return nil, fmt.Errorf("Azure Account Access Key not provided")
	}
	c := &config{
		accountName: accountName,
		accountKey:  accountKey,
	}
	return c, nil
}

func UploadBlob(srcFile, destFile, containerName string, c *config) (int, error) {

	// Create a default request pipeline using your storage account name and account key
	credential, err := azblob.NewSharedKeyCredential(c.accountName, c.accountKey)
	if err != nil {
		return 0, fmt.Errorf("Invalid credentials: %s", err)
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// TODO: Allow the URL to be a parameter
	// Setup the blob service URL endpoint
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

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
	log.Println(fmt.Sprintf("Uploading the file with blob name: %s", destFile))
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	if err != nil {
		return 0, fmt.Errorf("Error uploading file: %s", err)
	}

	return 1, nil
}
