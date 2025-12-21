.PHONY: build install clean test

BINARY_NAME=openskill
BUILD_DIR=build

build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/openskill

install: build
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed to /usr/local/bin/openskill"

clean:
	@rm -rf $(BUILD_DIR)

test:
	@go test -v ./...

deps:
	@go mod tidy
