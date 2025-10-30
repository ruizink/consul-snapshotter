module github.com/ruizink/consul-snapshotter

go 1.25.3

replace github.com/ruizink/consul-snapshotter/consul => ./consul

replace github.com/ruizink/consul-snapshotter/azure => ./azure

replace github.com/ruizink/consul-snapshotter/logger => ./logger

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.18.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.6.1
	github.com/hashicorp/consul v1.22.0
	github.com/hashicorp/consul/api v1.33.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/rboyer/safeio v0.2.3
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/pflag v1.0.6
	github.com/spf13/viper v1.20.1
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.1 // indirect
	github.com/armon/go-metrics v0.5.4 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.5.4 // indirect
	github.com/hashicorp/go-msgpack/v2 v2.1.3 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/raft v1.7.3 // indirect
	github.com/hashicorp/serf v0.10.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.8.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20250911091902-df9299821621 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.5.4

replace github.com/hashicorp/go-metrics => github.com/armon/go-metrics v0.5.4
