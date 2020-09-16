package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type consulConfig struct {
	URL         string        `json:"url"`
	Token       string        `json:"token"`
	LockKey     string        `json:"lock-key"`
	LockTimeout time.Duration `json:"lock-timeout"`
}

type localOutputConfig struct {
	DestinationPath string `json:"destination-path"`
}

type azureOutputConfig struct {
	ContainerName    string `json:"container-name"`
	ContainerPath    string `json:"container-path"`
	StorageAccount   string `json:"azure-storage-account"`
	StoraceAccessKey string `json:"azure-storage-access-key"`
}

type config struct {
	Cron              string            `json:"cron"`
	Outputs           []string          `json:"outputs"`
	ConsulConfig      consulConfig      `json:"consul"`
	AzureOutputConfig azureOutputConfig `json:"azure-blob"`
	LocalOutputConfig localOutputConfig `json:"local"`
	FilenamePrefix    string            `json:"filename-prefix"`
	FileExtension     string            `json:"file-extension"`
}

func regFlagString(flag string, value string, usage string) {
	if pflag.Lookup(flag) == nil {
		pflag.String(flag, value, usage)
	}
}

func regFlagStringSliceP(flag, shorthand string, value []string, usage string) {
	if pflag.Lookup(flag) == nil {
		pflag.StringSliceP(flag, shorthand, value, usage)
	}
}

func regFlagDuration(flag string, value time.Duration, usage string) {
	if pflag.Lookup(flag) == nil {
		pflag.Duration(flag, value, usage)
	}
}

func regFlagBoolP(flag, shorthand string, value bool, usage string) {
	if pflag.Lookup(flag) == nil {
		pflag.BoolP(flag, shorthand, value, usage)
	}
}

func (c *config) loadConfig() error {

	// set defaults
	viper.SetDefault("cron", "@every 1h")
	viper.SetDefault("filename-prefix", "consul-snapshot-")
	viper.SetDefault("file-extension", ".snap")
	viper.SetDefault("consul.url", "http://127.0.0.1:8500")
	viper.SetDefault("consul.lock-key", "consul-snapshot/.lock")
	viper.SetDefault("consul.lock-timeout", 10*time.Minute)
	viper.SetDefault("outputs", []string{"local"})
	viper.SetDefault("local.destination-path", ".")
	// bind env vars
	viper.BindEnv("consul.url", "CONSUL_HTTP_ADDR")
	viper.BindEnv("consul.token", "CONSUL_HTTP_TOKEN")
	viper.BindEnv("azure-blob.storage-account", "AZURE_STORAGE_ACCOUNT")
	viper.BindEnv("azure-blob.storage-access-key", "AZURE_STORAGE_ACCESS_KEY")

	// read command flags
	regFlagString("configdir", ".", "The path to look for the configuration file")
	regFlagString("cron", viper.GetString("cron"), "The cron expression to define when to run")
	regFlagString("filename-prefix", viper.GetString("filename-prefix"), "The prefix to use in the snapshot name")
	regFlagString("file-extension", viper.GetString("file-extension"), "The file extension to use in the snapshot name")
	regFlagString("consul.url", viper.GetString("consul.url"), "The Consul Agent URL")
	regFlagString("consul.token", viper.GetString("consul.token"), "The Consul Agent auth token")
	regFlagString("consul.lock-key", viper.GetString("consul.lock-key"), "The Key to use in the KV lock")
	regFlagDuration("consul.lock-timeout", viper.GetDuration("consul.lock-timeout"), "The timeout for the session lock")
	regFlagStringSliceP("outputs", "o", viper.GetStringSlice("outputs"), "The list of outputs to push the snapshot to")
	regFlagString("azure-blob.container-name", "", "The name of the Azure Blob container to use")
	regFlagString("azure-blob.container-path", "", "The path to use inside the Azure Blob container")
	regFlagString("azure-blob.storage-account", "", "The Azure Blob storage account to use")
	regFlagString("azure-blob.storage-access-key", "", "The Azure Blob storage access key to use")
	regFlagString("local.destination-path", viper.GetString("local.destination-path"), "The local path where to save the snapshots")
	regFlagBoolP("help", "h", false, "Prints this help message")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return err
	}

	// print usage if --help or -h
	if viper.GetBool("help") {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		pflag.PrintDefaults()
		os.Exit(0)
	}

	// load config from file
	viper.SetConfigName("config")
	viper.AddConfigPath(viper.GetString("configdir"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("[WARN] Could not load config file: %s \n", err)
	}

	// Consul config
	consulConfig := &consulConfig{}
	consulConfig.URL = viper.GetString("consul.url")
	consulConfig.Token = viper.GetString("consul.token")
	consulConfig.LockKey = viper.GetString("consul.lock-key")
	consulConfig.LockTimeout = viper.GetDuration("consul.lock-timeout")

	// Azure Blob output config
	azureOutputConfig := &azureOutputConfig{}
	azureOutputConfig.ContainerName = viper.GetString("azure-blob.container-name")
	azureOutputConfig.ContainerPath = viper.GetString("azure-blob.container-path")
	azureOutputConfig.StorageAccount = viper.GetString("azure-blob.storage-account")
	azureOutputConfig.StoraceAccessKey = viper.GetString("azure-blob.storage-access-key")

	// Local output config
	localOutputConfig := &localOutputConfig{}
	localOutputConfig.DestinationPath = viper.GetString("local.destination-path")

	c.Cron = viper.GetString("cron")
	c.FilenamePrefix = viper.GetString("filename-prefix")
	c.FileExtension = viper.GetString("file-extension")
	c.Outputs = viper.GetStringSlice("outputs")
	c.ConsulConfig = *consulConfig
	c.AzureOutputConfig = *azureOutputConfig
	c.LocalOutputConfig = *localOutputConfig

	return nil
}

func (c *config) String() string {
	conf, err := json.MarshalIndent(c, "", "  ")
	// conf, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}
	return string(conf)
}
