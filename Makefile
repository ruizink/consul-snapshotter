# build
BUILD_PATH     := build
BIN_PATH       := $(BUILD_PATH)/bin
PACKAGE_PATH   := $(BUILD_PATH)/package
CHECKSUM_PATH  := $(BUILD_PATH)/checksum
PACKAGE_FORMAT ?= zip
BIN_NAME       := consul-snapshotter

# git
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_SHA   ?= $(shell git rev-parse --short HEAD)
GIT_TAG   ?= $(shell git describe --tags --exact-match "$(GIT_SA)" 2>/dev/null || true)

# app
OS         ?= linux
ARCH       ?= amd64
BIN_TARGET := $(BIN_PATH)/$(OS)/$(ARCH)
VERSION    ?= $(GIT_TAG:v%=%)
ifeq ($(VERSION),)
	VERSION := dev
endif
BUILD_DATE ?= $(shell date --iso=seconds)
T          := github.com/ruizink/consul-snapshotter
LDFLAGS    := -X '$(T)/version.Version=$(VERSION)' -X '$(T)/version.BuildDate=$(BUILD_DATE)' -X '$(T)/version.GitCommit=$(GIT_SHA)$(GIT_DIRTY)'

.PHONY: mkdirs build build-docker package checksum clean start-docker-env stop-docker-env

build:
	$(info Building binary for $(OS) $(ARCH))
	GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags "$(LDFLAGS)" -o $(BIN_TARGET)/ -trimpath -buildvcs=false

build-docker: export OS=linux
build-docker: build
	$(info Building docker image for $(OS)/$(ARCH))
	docker build \
		--tag $(BIN_NAME):$(VERSION) \
		--platform $(OS)/$(ARCH) \
		.

mkdirs:
	@mkdir -p $(PACKAGE_PATH)
	@mkdir -p $(CHECKSUM_PATH)

package: mkdirs
ifneq ($(wildcard $(BIN_TARGET)/$(BIN_NAME)),)
ifeq ($(PACKAGE_FORMAT), zip)
	$(info Creating zip for $(OS) $(ARCH))
	zip --junk-paths $(PACKAGE_PATH)/$(BIN_NAME)_$(VERSION)_$(OS)_$(ARCH).zip $(BIN_TARGET)/$(BIN_NAME)
else
ifeq ($(PACKAGE_FORMAT), tgz)
	$(info Creating tgz for $(OS) $(ARCH))
	tar -czvf $(PACKAGE_PATH)/$(BIN_NAME)_$(VERSION)_$(OS)_$(ARCH).tar.gz -C $(BIN_TARGET) $(BIN_NAME)
endif
endif
else
	$(error Could not find a build for $(OS) $(ARCH))
endif

checksum: mkdirs
ifneq ($(wildcard $(PACKAGE_PATH)/$(BIN_NAME)*),)
	$(info Generating checksum)
	@cd $(PACKAGE_PATH) && sha256sum $(BIN_NAME)* | tee $(CHECKSUM_PATH)/SHA256SUM
else
	$(error Could not find files to checksum)
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
