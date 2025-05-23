package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/ruizink/consul-snapshotter/logger"
	"github.com/ruizink/consul-snapshotter/version"
)

type consulConfig struct {
	URL         string        `json:"url"`
	Token       string        `json:"token"`
	LockKey     string        `json:"lock-key"`
	LockTimeout time.Duration `json:"lock-timeout"`
}

type localOutputConfig struct {
	DestinationPath   string        `json:"destination-path"`
	RetentionPeriod   time.Duration `json:"retention-period"`
	CreateDestination bool          `json:"create-destination"`
}

type azureOutputConfig struct {
	ContainerName    string        `json:"container-name"`
	ContainerPath    string        `json:"container-path"`
	StorageAccount   string        `json:"storage-account"`
	CloudDomain      string        `json:"cloud-domain"`
	StorageAccessKey string        `json:"storage-access-key"`
	StorageSASToken  string        `json:"storage-sas-token"`
	CreateContainer  bool          `json:"create-container"`
	BlockSize        int64         `json:"block-size"`
	Parallelism      uint16        `json:"parallelism"`
	RetentionPeriod  time.Duration `json:"retention-period"`
	Emulated         bool          `json:"emulated"`
	EmulatorUrl      string        `json:"emulator-url"`
}

type config struct {
	Cron              string            `json:"cron"`
	Outputs           []string          `json:"outputs"`
	ConsulConfig      consulConfig      `json:"consul"`
	AzureOutputConfig azureOutputConfig `json:"azure-blob"`
	LocalOutputConfig localOutputConfig `json:"local"`
	FilenamePrefix    string            `json:"filename-prefix"`
	FileExtension     string            `json:"file-extension"`
	LogLevel          string            `json:"log-level"`
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

func regFlagBool(flag string, value bool, usage string) {
	if pflag.Lookup(flag) == nil {
		pflag.Bool(flag, value, usage)
	}
}

func regFlagInt64(flag string, value int64, usage string) {
	if pflag.Lookup(flag) == nil {
		pflag.Int64(flag, value, usage)
	}
}

func regFlagUint(flag string, value uint, usage string) {
	if pflag.Lookup(flag) == nil {
		pflag.Uint(flag, value, usage)
	}
}

func (c *config) loadConfig() error {

	// set defaults
	viper.SetDefault("filename-prefix", "consul-snapshot-")
	viper.SetDefault("file-extension", ".snap")
	viper.SetDefault("configdir", ".")
	viper.SetDefault("log-level", "info")
	viper.SetDefault("consul.url", "http://127.0.0.1:8500")
	viper.SetDefault("consul.lock-key", "consul-snapshotter/.lock")
	viper.SetDefault("consul.lock-timeout", 10*time.Minute)
	viper.SetDefault("outputs", []string{"local"})
	viper.SetDefault("local.destination-path", ".")
	viper.SetDefault("local.create-destination", false)
	viper.SetDefault("local.retention-period", 0)
	viper.SetDefault("azure-blob.cloud-domain", "blob.core.windows.net")
	viper.SetDefault("azure-blob.create-container", false)
	viper.SetDefault("azure-blob.block-size", 4*1024*1024)
	viper.SetDefault("azure-blob.parallelism", 16)
	viper.SetDefault("azure-blob.retention-period", 0)
	viper.SetDefault("azure-blob.emulated", false)
	viper.SetDefault("azure-blob.emulator-url", "http://127.0.0.1:10000")

	// read command flags
	regFlagString("configdir", viper.GetString("configdir"), "The path to look for the configuration file")
	regFlagString("cron", viper.GetString("cron"), "Cron expression to define when to run")
	regFlagString("filename-prefix", viper.GetString("filename-prefix"), "Prefix to use in the snapshot name")
	regFlagString("file-extension", viper.GetString("file-extension"), "File extension to use in the snapshot name")
	regFlagString("log-level", viper.GetString("log-level"), "Verbosity (info, warn, debug) of the log")
	regFlagString("consul.url", viper.GetString("consul.url"), "Consul Agent URL")
	regFlagString("consul.token", viper.GetString("consul.token"), "Consul Agent authentication token")
	regFlagString("consul.lock-key", viper.GetString("consul.lock-key"), "Key to use in the KV lock")
	regFlagDuration("consul.lock-timeout", viper.GetDuration("consul.lock-timeout"), "Timeout for the session lock")
	regFlagStringSliceP("outputs", "o", viper.GetStringSlice("outputs"), "List of outputs to push the snapshot to")
	regFlagString("azure-blob.container-name", "", "Name of the Azure Blob container to use")
	regFlagString("azure-blob.container-path", "", "Path to use inside the Azure Blob container")
	regFlagString("azure-blob.storage-account", "", "Azure Blob storage account to use")
	regFlagString("azure-blob.cloud-domain", viper.GetString("azure-blob.cloud-domain"), "The domain for the Azure Blob service, depending on the cloud you are using")
	regFlagString("azure-blob.storage-access-key", "", "Azure Blob storage access key to use (mutually exclusive with azure-blob.storage-sas-token)")
	regFlagString("azure-blob.storage-sas-token", "", "Azure Blob storage SAS token to use (mutually exclusive with azure-blob.storage-access-key)")
	regFlagBool("azure-blob.create-container", viper.GetBool("azure-blob.create-container"), "Behavior when the container-name does not exist (default: false)")
	regFlagInt64("azure-blob.block-size", viper.GetInt64("azure-blob.block-size"), "Size in bytes of each block")
	regFlagUint("azure-blob.parallelism", viper.GetUint("azure-blob.parallelism"), "Maximum number of blocks to upload in parallel")
	regFlagDuration("azure-blob.retention-period", viper.GetDuration("azure-blob.retention-period"), "Duration that Azure Blob snapshots need to be retained (default: \"0s\" - keep forever)")
	regFlagBool("azure-blob.emulated", viper.GetBool("azure-blob.emulated"), "If enabled, it will try to connect to a local Azure Blob Emulator using <emulator-url>/<storage-account>/<container-name> (default: false)")
	regFlagString("azure-blob.emulator-url", viper.GetString("azure-blob.emulator-url"), "URL of the Azure Blob Emulator")
	regFlagString("local.destination-path", viper.GetString("local.destination-path"), "Local path where to save the snapshots")
	regFlagBool("local.create-destination", viper.GetBool("local.create-destination"), "Behavior when the destination-path does not exist (default: false)")
	regFlagDuration("local.retention-period", viper.GetDuration("local.retention-period"), "Duration that Local snapshots need to be retained (default: \"0s\" - keep forever)")
	regFlagBoolP("help", "h", false, "Prints this help message")
	regFlagBoolP("version", "V", false, "Prints the version")

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

	// print version if --version or -V
	if viper.GetBool("version") {
		fmt.Fprintf(os.Stderr, "Version: %s\n", version.Version)
		fmt.Fprintf(os.Stderr, "(Build date: %s, Git commit: %s)\n", version.BuildDate, version.GitCommit)
		os.Exit(0)
	}

	// bind env vars
	viper.BindEnv("consul.url", "CONSUL_HTTP_ADDR")
	viper.BindEnv("consul.token", "CONSUL_HTTP_TOKEN")
	viper.BindEnv("azure-blob.cloud-domain", "AZURE_CLOUD_DOMAIN")
	viper.BindEnv("azure-blob.storage-account", "AZURE_STORAGE_ACCOUNT")
	viper.BindEnv("azure-blob.storage-access-key", "AZURE_STORAGE_ACCESS_KEY")
	viper.BindEnv("azure-blob.storage-sas-token", "AZURE_STORAGE_SAS_TOKEN")

	// load config from file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(viper.GetString("configdir"))

	if err := viper.ReadInConfig(); err != nil {
		logger.Warn("Could not load config file: ", err)
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
	azureOutputConfig.CloudDomain = viper.GetString("azure-blob.cloud-domain")
	azureOutputConfig.StorageAccessKey = viper.GetString("azure-blob.storage-access-key")
	azureOutputConfig.StorageSASToken = viper.GetString("azure-blob.storage-sas-token")
	azureOutputConfig.CreateContainer = viper.GetBool("azure-blob.create-container")
	azureOutputConfig.BlockSize = viper.GetInt64("azure-blob.block-size")
	azureOutputConfig.Parallelism = uint16(viper.GetUint("azure-blob.parallelism"))
	azureOutputConfig.RetentionPeriod = viper.GetDuration("azure-blob.retention-period")
	azureOutputConfig.Emulated = viper.GetBool("azure-blob.emulated")
	azureOutputConfig.EmulatorUrl = viper.GetString("azure-blob.emulator-url")

	// Local output config
	localOutputConfig := &localOutputConfig{}
	localOutputConfig.DestinationPath = viper.GetString("local.destination-path")
	localOutputConfig.RetentionPeriod = viper.GetDuration("local.retention-period")
	localOutputConfig.CreateDestination = viper.GetBool("local.create-destination")

	c.Cron = viper.GetString("cron")
	c.FilenamePrefix = viper.GetString("filename-prefix")
	c.FileExtension = viper.GetString("file-extension")
	c.LogLevel = viper.GetString("log-level")
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
