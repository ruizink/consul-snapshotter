BIN_PATH=bin
ARCHIVE_PATH=archive
BINARY_NAME=consul-snapshotter

build: dependencies build-linux build-darwin

dependencies:
	@echo "Fetching dependencies"
	go get -v -t -d ./...

build-linux:
	@echo "Building binary for Linux"
	GOOS=linux GOARCH=amd64 go build -o $(BIN_PATH)/$(BINARY_NAME)_linux_amd64 -v .

build-darwin:
	@echo "Building binary for Darwin"
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_PATH)/$(BINARY_NAME)_darwin_amd64 -v .

archive: linux-zip linux-tgz darwin-zip darwin-tgz

.mk-archive:
	mkdir -p $(ARCHIVE_PATH)

linux-zip: .mk-archive
	@echo "Creating zip for linux_amd64"
	zip --junk-paths $(ARCHIVE_PATH)/$(BINARY_NAME)_linux_amd64.zip $(BIN_PATH)/$(BINARY_NAME)_linux_amd64

darwin-zip: .mk-archive
	@echo "Creating zip for darwin_amd64"
	zip --junk-paths $(ARCHIVE_PATH)/$(BINARY_NAME)_darwin_amd64.zip $(BIN_PATH)/$(BINARY_NAME)_darwin_amd64

linux-tgz: .mk-archive
	@echo "Creating tgz for linux_amd64"
	tar -czvf $(ARCHIVE_PATH)/$(BINARY_NAME)_linux_amd64.tar.gz -C $(BIN_PATH) $(BINARY_NAME)_linux_amd64

darwin-tgz: .mk-archive
	@echo "Creating tgz for darwin_amd64"
	tar -czvf $(ARCHIVE_PATH)/$(BINARY_NAME)_darwin_amd64.tar.gz -C $(BIN_PATH) $(BINARY_NAME)_darwin_amd64

clean:
	go clean
	rm -rf $(BIN_PATH)
	rm -rf $(ARCHIVE_PATH)

start-docker-env:
	@docker compose up -d

stop-docker-env:
	@docker compose down
	@docker compose rm
