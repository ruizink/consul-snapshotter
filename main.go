package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ruizink/consul-snapshotter/outputs"

	"github.com/robfig/cron"
	"github.com/ruizink/consul-snapshotter/consul"
	"github.com/ruizink/consul-snapshotter/logger"
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
					logger.Info("Caught SIGHUP. Triggering a config reload.")
					c.loadConfig()
				case os.Interrupt:
					cancel()
					os.Exit(1)
				}
			case <-ctx.Done():
				// logger.Info("Terminated.")
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
	logger.SetLevel(c.LogLevel)

	runSnapshotter := func() error {
		logger.Info("####################################################################################")
		logger.Info("===> Performing Consul snapshot backup procedure...")
		defer logger.Info("####################################################################################")
		// create new consul client
		consulWorker, err := consul.NewConsul(c.ConsulConfig.URL, c.ConsulConfig.Token, c.ConsulConfig.LockKey, c.ConsulConfig.LockTimeout)
		if err != nil {
			logger.Error("Could not create a consul client: ", err)
			return err
		}

		// acquire lock
		if err := consulWorker.AcquireLock(); err != nil {
			logger.Error("Could not acquire lock: ", err)
			return err
		}
		logger.Debug("Acquired lock for session ID: ", consulWorker.SessionID)

		// Start renewing the session until doneChan is closed
		doneChan := make(chan struct{})
		go worker.RenewSession(doneChan)

		// Close the channel used for session renewal
		defer close(doneChan)

		// Get consul snapshot
		snap, err := consul.GetSnapshot(worker)
		if err != nil {
			logger.Error("Could not perform snapshot: ", err)
			return err
		}

		// Export the snapshot to all the configured outputs
		processOutputs(snap, c)

		// Remove the temporary snapshot
		os.Remove(snap)
		logger.Info("Starting with cron expression: ", c.Cron)

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
		logger.Info("Starting a single execution...")
	}
}

func processOutputs(snap string, c *config) {

	outputFileName := fmt.Sprintf("%s%v%s", c.FilenamePrefix, time.Now().UnixNano(), c.FileExtension)

	for _, output := range c.Outputs {
		switch output {
		case "local":
			logger.Info("===> Processing output: local")

			o := &outputs.LocalOutput{
				DestinationPath:   c.LocalOutputConfig.DestinationPath,
				Filename:          outputFileName,
				CreateDestination: c.LocalOutputConfig.CreateDestination,
				RetentionPeriod:   c.LocalOutputConfig.RetentionPeriod,
			}
			if err := o.Save(snap); err != nil {
				logger.Error(err)
				errors = multierror.Append(errors, err)
				continue
			}
			if err := o.ApplyRetentionPolicy(); err != nil {
				logger.Error(err)
				errors = multierror.Append(errors, err)
				continue
			}
		case "azure_blob":
			// Upload to Azure
			logger.Info("===> Processing output: azure_blob")

			o := &outputs.AzureBlobOutput{
				AzureConfig: &azure.AzureConfig{
					ContainerName:    c.AzureOutputConfig.ContainerName,
					ContainerPath:    c.AzureOutputConfig.ContainerPath,
					Filename:         outputFileName,
					StorageAccount:   c.AzureOutputConfig.StorageAccount,
					StorageAccessKey: c.AzureOutputConfig.StorageAccessKey,
					StorageSASToken:  c.AzureOutputConfig.StorageSASToken,
					CreateContainer:  c.AzureOutputConfig.CreateContainer,
					BlockSize:        c.AzureOutputConfig.BlockSize,
					Parallelism:      c.AzureOutputConfig.Parallelism,
					Emulated:         c.AzureOutputConfig.Emulated,
					EmulatorUrl:      c.AzureOutputConfig.EmulatorUrl,
				},
				RetentionPeriod: c.AzureOutputConfig.RetentionPeriod,
			}

			if err := o.Save(snap); err != nil {
				logger.Error(err)
				errors = multierror.Append(errors, err)
				continue
			}
			o.Save(snap)

			if err := o.ApplyRetentionPolicy(); err != nil {
				logger.Error(err)
				errors = multierror.Append(errors, err)
				continue
			}
		}
	}
}
