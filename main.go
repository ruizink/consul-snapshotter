package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/robfig/cron"
	"github.com/ruizink/consul-snapshot/azure"
	"github.com/ruizink/consul-snapshot/consul"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	c := &config{}

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					log.Printf("Caught SIGHUP. Triggering a config reload.")
					c.loadConfig()
				case os.Interrupt:
					cancel()
					os.Exit(1)
				}
			case <-ctx.Done():
				// log.Printf("Terminated.")
				os.Exit(0)
			}
		}
	}()

	if err := run(ctx, c, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, c *config, stdout io.Writer) error {
	c.loadConfig()
	log.SetOutput(os.Stdout)

	cron := cron.New()

	log.Println("Starting with cron expression:", c.Cron)

	cron.AddFunc(c.Cron, func() {
		log.Println("####################################################################################")
		log.Println("Starting snapshot backup procedure...")
		defer log.Println("####################################################################################")
		// create new worker
		worker, err := consul.NewWorker(c.ConsulConfig.URL, c.ConsulConfig.Token, c.ConsulConfig.LockKey, c.ConsulConfig.LockTimeout)
		if err != nil {
			log.Println("Could not create a worker:", err)
			return
		}

		// acquire lock
		lock, err := consul.AcquireLock(worker)
		if err != nil || !lock {
			log.Println("Could not acquire lock. Skipping...")
			return
		}
		log.Println("Acquired lock for session ID", worker.SessionID)

		// Start renewing the session until doneChan is closed
		doneChan := make(chan struct{})
		go worker.RenewSession(doneChan)

		// Close the channel used for session renewal
		defer close(doneChan)

		// Get consul snapshot
		snap, err := consul.GetSnapshot(worker)
		if err != nil {
			log.Println("Could not perform snapshot:", err)
			return
		}

		// Export the snapshot to all the configured outputs
		processOutputs(snap, c)

		// Remove the temporary snapshot
		os.Remove(snap)

		// Release the lock
		if err := consul.ReleaseLock(worker); err != nil {
			log.Println("Could not release lock:", err)
		}
		log.Println("Released lock for session ID", worker.SessionID)
	})
	cron.Start()

	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func processOutputs(snap string, c *config) {

	outputFileName := fmt.Sprintf("%s%v%s", c.FilenamePrefix, time.Now().UnixNano(), c.FileExtension)

	for _, output := range c.Outputs {
		switch output {
		case "local":
			// Rename temp file if local output is defined (make sure to run this output as last one)
			defer func() {
				dstFile := path.Join(c.LocalOutputConfig.DestinationPath, outputFileName)
				log.Println("Processing output: local")
				if err := os.Rename(snap, dstFile); err != nil {
					log.Println("Error writing snapshot file: ", err)
					return
				}
				log.Println("Saved snapshot to:", dstFile)
			}()
		case "azure_blob":
			// Upload to Azure
			func() {
				log.Println("Processing output: azure_blob")
				destFile := path.Join(c.AzureOutputConfig.ContainerPath, outputFileName)
				config, err := azure.AzureConfig(c.AzureOutputConfig.StorageAccount, c.AzureOutputConfig.StoraceAccessKey)
				if err != nil {
					log.Println("Invalid Azure config:", err)
					return
				}
				_, err = azure.UploadBlob(snap, destFile, c.AzureOutputConfig.ContainerName, config)
				if err != nil {
					log.Println("Error uploading snapshot file:", err)
					return
				}
				log.Println("Uploaded snapshot to:", destFile)
			}()
		}
	}
}
