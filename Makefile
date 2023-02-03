# build
BUILD_PATH=$(CURDIR)/build
BIN_PATH=$(BUILD_PATH)/bin
ARCHIVE_PATH=$(BUILD_PATH)/archive
ARCHIVE_FORMAT?=zip
BIN_NAME=consul-snapshotter

# git
GIT_SHA?=$(shell git rev-parse --short HEAD)
GIT_TAG?=$(shell git describe --tags --exact-match "$(GIT_SA)" 2>/dev/null || true)

# golang
OS?=linux
ARCH?=amd64
BIN_TARGET=$(BIN_PATH)/$(OS)/$(ARCH)
VERSION?=$(GIT_TAG:v%=%)
ifeq ($(VERSION),)
	VERSION=$(GIT_SHA)
endif
BUILD_DATE?=$(shell date --iso=seconds)
T=github.com/ruizink/consul-snapshotter
LDFLAGS=-X '$(T)/version.Version=$(VERSION)' -X '$(T)/version.BuildDate=$(BUILD_DATE)' -X '$(T)/version.GitCommit=$(GIT_SHA)'

.PHONY: archivedir build archive checksum clean start-docker-env stop-docker-env

build:
	$(info Building binary for $(OS) $(ARCH))
	GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags "$(LDFLAGS)" -o $(BIN_TARGET)/ -trimpath -buildvcs=false

archivedir:
	@mkdir -p $(ARCHIVE_PATH)

archive: archivedir
ifneq ($(wildcard $(BIN_TARGET)/$(BIN_NAME)),)
ifeq ($(ARCHIVE_FORMAT), zip)
	$(info Creating zip for $(OS) $(ARCH))
	zip --junk-paths $(ARCHIVE_PATH)/$(BIN_NAME)_$(VERSION)_$(OS)_$(ARCH).zip $(BIN_TARGET)/$(BIN_NAME)
else
ifeq ($(ARCHIVE_FORMAT), tgz)
	$(info Creating tgz for $(OS) $(ARCH))
	tar -czvf $(ARCHIVE_PATH)/$(BIN_NAME)_$(VERSION)_$(OS)_$(ARCH).tar.gz -C $(BIN_TARGET) $(BIN_NAME)
endif
endif
else
	$(info Could not find a build for $(OS) $(ARCH). Skipping...)
endif

checksum: archivedir
ifneq ($(wildcard $(ARCHIVE_PATH)/$(BIN_NAME)*),)
	$(info Generating checksum)
	@cd $(ARCHIVE_PATH) && sha256sum $(BIN_NAME)* > SHA256SUM
else
	$(info Could not find files to checksum. Skipping...)
endif

clean:
	$(info Cleaning go environment)
	@go clean
	$(info Removing $(BUILD_PATH) directory)
	@rm -rf $(BUILD_PATH)

start-docker-env:
	$(info Starting containers with docker compose)
	@docker compose -f ./docker/docker-compose.yml up -d

stop-docker-env:
	$(info Stopping containers with docker compose)
	@docker compose -f ./docker/docker-compose.yml down
	@docker compose -f ./docker/docker-compose.yml rm
