module github.com/ruizink/consul-snapshotter

go 1.19

replace github.com/ruizink/consul-snapshotter/consul => ./consul

replace github.com/ruizink/consul-snapshotter/azure => ./azure

replace github.com/ruizink/consul-snapshotter/logger => ./logger

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.3.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.6.1
	github.com/hashicorp/consul v1.13.7
	github.com/hashicorp/consul/api v1.18.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/rboyer/safeio v0.2.1
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.4.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/raft v1.3.11 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hashicorp/yamux v0.0.0-20211028200310-0bc27b27de87 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
