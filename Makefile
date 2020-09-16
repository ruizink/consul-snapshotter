BIN_PATH=bin
BINARY_NAME=consul-snapshotter

build: build-linux build-darwin

build-linux:
	@echo "Building binary for Linux"
	GOOS=linux GOARCH=amd64 go build -o $(BIN_PATH)/linux_amd64/$(BINARY_NAME) -v

build-darwin:
	@echo "Building binary for Darwin"
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_PATH)/darwin_amd64/$(BINARY_NAME) -v

clean:
	go clean
	rm -rf $(BIN_PATH)