# Makefile for zlx

PROJECT_NAME := zlx
BINARY_NAME := $(PROJECT_NAME)
BUILD_DIR := build
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")
LDFLAGS := -ldflags "-s -w"
UPX_OPTIONS := --best -q

all: build compress

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(GO_FILES)
	@echo "Build complete!"

compress: build
	@echo "Compressing with UPX..."
	upx $(UPX_OPTIONS) $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Compression complete!"

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	@echo "Clean complete!"

run: build
	@echo "Running..."
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  all        Builds and compresses the binary"
	@echo "  build      Builds the binary"
	@echo "  compress   Compresses the binary with UPX"
	@echo "  clean      Cleans the build directory"
	@echo "  run        Runs the application"
	@echo "  help       Shows this help message"

.DEFAULT_GOAL := help