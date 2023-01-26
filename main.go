package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/robfig/cron"

	"github.com/ruizink/consul-snapshotter/azure"
	"github.com/ruizink/consul-snapshotter/consul"
	"github.com/ruizink/consul-snapshotter/logger"
	"github.com/ruizink/consul-snapshotter/outputs"
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
		// fmt.Fprintf(os.Stderr, "%s\n", err)
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

		// Cleanup: Release the lock
		defer func() {
			if err := consulWorker.ReleaseLock(); err != nil {
				logger.Error("Could not release lock: ", err)
			}
			logger.Debug("Released lock for session ID: ", consulWorker.SessionID)
		}()

		// Start renewing the session until doneChan is closed
		doneChan := make(chan struct{})
		go consulWorker.RenewSession(doneChan)

		// Cleanup: Close the channel used for session renewal
		defer close(doneChan)

		// Get consul snapshot
		snap, err := consulWorker.GetSnapshot()
		if err != nil {
			logger.Error("Could not perform snapshot: ", err)
			return err
		}

		// Cleanup: Remove the temporary snapshot
		defer os.Remove(snap)

		// Export the snapshot to all the configured outputs
		if err := processOutputs(snap, c); err != nil {
			return err
		}

		return nil
	}

	runSnapshotterCron := func() {
		_ = runSnapshotter()
	}

	if c.Cron != "" {
		logger.Info("Starting with cron expression: ", c.Cron)

		cron := cron.New()
		cron.AddFunc(c.Cron, runSnapshotterCron)
		cron.Start()

		for {
			select {
			case <-ctx.Done():
				return nil
			}
		}
	} else {
		logger.Info("Starting a single execution...")
		return runSnapshotter()
	}
}

func processOutputs(snap string, c *config) error {

	var errors error

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

			if err := o.ApplyRetentionPolicy(); err != nil {
				logger.Error(err)
				errors = multierror.Append(errors, err)
				continue
			}
		}
	}
	return errors
}
