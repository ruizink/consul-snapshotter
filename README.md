# consul-snapshot

A binary to help you perform Consul backups periodically

## Build

`go run main.go config.go`

Example: Run every 10 seconds and save the snapshot to the current directory

`go run main.go config.go --cron "@every 10s" --local.destination-path "." --outputs "local"`

## TODO

* Write documentation
* Write tests