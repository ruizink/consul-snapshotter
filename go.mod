module github.com/ruizink/consul-snapshotter

go 1.22
toolchain go1.23.7

replace github.com/ruizink/consul-snapshotter/consul => ./consul

replace github.com/ruizink/consul-snapshotter/azure => ./azure

replace github.com/ruizink/consul-snapshotter/logger => ./logger

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.10.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.3.1
	github.com/hashicorp/consul v1.20.1
	github.com/hashicorp/consul/api v1.29.4
	github.com/hashicorp/go-multierror v1.1.1
	github.com/rboyer/safeio v0.2.3
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.18.2
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.5.2 // indirect
	github.com/armon/go-metrics v0.5.3 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/go-msgpack/v2 v2.1.2 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/raft v1.6.1 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.4.1

replace github.com/hashicorp/go-metrics => github.com/armon/go-metrics v0.4.1
