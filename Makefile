# build
BUILD_PATH     := $(CURDIR)/build
BIN_PATH       := $(BUILD_PATH)/bin
CHECKSUM_PATH  := $(BUILD_PATH)/checksum
BIN_NAME       := consul-snapshotter

# git
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_SHA   ?= $(shell git rev-parse --short HEAD)
GIT_TAG   ?= $(shell git describe --tags --exact-match "$(GIT_SA)" 2>/dev/null || true)

# app
OS         ?= linux
ARCH       ?= amd64
VERSION    ?= $(GIT_TAG:v%=%)
ifeq ($(VERSION),)
	VERSION := dev
endif
BUILD_DATE ?= $(shell date --iso=seconds)
T          := github.com/ruizink/consul-snapshotter
LDFLAGS    := -X '$(T)/version.Version=$(VERSION)' -X '$(T)/version.BuildDate=$(BUILD_DATE)' -X '$(T)/version.GitCommit=$(GIT_SHA)$(GIT_DIRTY)'

.PHONY: mkdirs build build-docker checksum clean start-docker-env stop-docker-env

build:
	$(info Building binary for $(OS) $(ARCH))
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BIN_PATH)/$(BIN_NAME)_$(OS)_$(ARCH) -trimpath -buildvcs=false

build-docker: export OS=linux
build-docker: build
	$(info Building docker image for $(OS)/$(ARCH))
	docker build \
		--tag $(BIN_NAME):$(VERSION) \
		--platform $(OS)/$(ARCH) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--target alpine \
		.

mkdirs:
	@mkdir -p $(CHECKSUM_PATH)

checksum: mkdirs
ifneq ($(wildcard $(BIN_PATH)/$(BIN_NAME)*),)
	$(info Generating checksum)
	@cd $(BIN_PATH) && sha256sum $(BIN_NAME)* | tee $(CHECKSUM_PATH)/SHA256SUM
else
	$(error Could not find files to checksum with the following pattern: $(BIN_PATH)/$(BIN_NAME)*)
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
