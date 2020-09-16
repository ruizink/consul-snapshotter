# consul-snapshot

A binary to help you perform Consul backups periodically

## Build

Build for all platforms:

`make build`

Build for specific platform:

`make build-linux`

Clean the build environment:

`make clean`

## Run

Run with default settings:

`consul-snapshot`

Run every 10 seconds and save the snapshot to the current directory

`consul-snapshot --cron "@every 10s" --local.destination-path "." --outputs "local"`

Run with config from file:

`consul-snapshot --configdir /etc/consul-snapshot`

Usage:

`consul-snapshot --help`

```text
Usage of consul-snapshot:
      --azure-blob.container-name string       The name of the Azure Blob container to use
      --azure-blob.container-path string       The path to use inside the Azure Blob container
      --azure-blob.storage-access-key string   The Azure Blob storage access key to use
      --azure-blob.storage-account string      The Azure Blob storage account to use
      --configdir string                       The path to look for the configuration file (default ".")
      --consul.lock-key string                 The Key to use in the KV lock (default "consul-snapshot/.lock")
      --consul.lock-timeout duration           The timeout for the session lock (default 10m0s)
      --consul.token string                    The Consul Agent auth token
      --consul.url string                      The Consul Agent URL (default "http://127.0.0.1:8500")
      --cron string                            The cron expression to define when to run (default "@every 1h")
      --file-extension string                  The file extension to use in the snapshot name (default ".snap")
      --filename-prefix string                 The prefix to use in the snapshot name (default "consul-snapshot-")
  -h, --help                                   Prints this help message
      --local.destination-path string          The local path where to save the snapshots (default ".")
  -o, --outputs strings                        The list of outputs to push the snapshot to (default [local])
```

## TODO

* Write documentation
* Write tests
