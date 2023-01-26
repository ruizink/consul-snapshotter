# consul-snapshotter

A tool to help you perform Consul backups periodically

## Build

Build for all platforms:

`make build`

Build for specific platform:

`make build-linux`

Clean the build environment:

`make clean`

## Run

Run with default settings:

`consul-snapshotter`

Run every 10 seconds and save the snapshot to the current directory

`consul-snapshotter --cron "@every 10s" --local.destination-path "." --outputs "local"`

Run with config from file:

`consul-snapshotter --configdir /etc/consul-snapshotter`

Usage:

`consul-snapshotter --help`

```text
Usage of consul-snapshotter:
      --azure-blob.block-size int              Size in bytes of each block (default 4194304)
      --azure-blob.container-name string       Name of the Azure Blob container to use
      --azure-blob.container-path string       Path to use inside the Azure Blob container
      --azure-blob.create-container            Behavior when the container-name does not exist (default: false)
      --azure-blob.emulated                    If enabled, it will try to connect to a local Azure Blob Emulator using <emulator-url>/<storage-account>/<container-name> (default: false)
      --azure-blob.emulator-url string         URL of the Azure Blob Emulator (default "http://127.0.0.1:10000")
      --azure-blob.parallelism uint            Maximum number of blocks to upload in parallel (default 16)
      --azure-blob.retention-period duration   Duration that Azure Blob snapshots need to be retained (default: "0s" - keep forever)
      --azure-blob.storage-access-key string   Azure Blob storage access key to use (mutually exclusive with azure-blob.storage-sas-token)
      --azure-blob.storage-account string      Azure Blob storage account to use
      --azure-blob.storage-sas-token string    Azure Blob storage SAS token to use (mutually exclusive with azure-blob.storage-access-key)
      --configdir string                       The path to look for the configuration file (default ".")
      --consul.lock-key string                 Key to use in the KV lock (default "consul-snapshotter/.lock")
      --consul.lock-timeout duration           Timeout for the session lock (default 10m0s)
      --consul.token string                    Consul Agent authentication token
      --consul.url string                      Consul Agent URL (default "http://127.0.0.1:8500")
      --cron string                            Cron expression to define when to run
      --file-extension string                  File extension to use in the snapshot name (default ".snap")
      --filename-prefix string                 Prefix to use in the snapshot name (default "consul-snapshot-")
  -h, --help                                   Prints this help message
      --local.create-destination               Behavior when the destination-path does not exist (default: false)
      --local.destination-path string          Local path where to save the snapshots (default ".")
      --local.retention-period duration        Duration that Local snapshots need to be retained (default: "0s" - keep forever)
      --log-level string                       Verbosity (info, warn, debug) of the log (default "info")
  -o, --outputs strings                        List of outputs to push the snapshot to (default [local])
```
