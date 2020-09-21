module github.com/ruizink/consul-snapshotter

go 1.14

replace github.com/ruizink/consul-snapshotter/consul => ./consul

replace github.com/ruizink/consul-snapshotter/azure => ./azure

require (
	github.com/Azure/azure-pipeline-go v0.2.2
	github.com/Azure/azure-storage-blob-go v0.10.0
	github.com/hashicorp/consul v1.8.1
	github.com/hashicorp/consul/api v1.5.0
	github.com/hashicorp/consul/sdk v0.5.0
	github.com/kr/text v0.1.0
	github.com/mitchellh/cli v1.1.0
	github.com/mitchellh/mapstructure v1.3.3
	github.com/namsral/flag v1.7.4-pre
	github.com/rboyer/safeio v0.2.1
	github.com/robfig/cron v1.2.0
	github.com/robfig/cron/v3 v3.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.5.1
	golang.org/x/tools v0.0.0-20200828161849-5deb26317202 // indirect
)
