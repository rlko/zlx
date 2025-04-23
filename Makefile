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

install: build
	@echo "Installing to /usr/local/bin..."
	install -D $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Installation to /usr/bin complete!"

user-install: build
	@echo "Installing to user's local bin directory..."
	mkdir -p ~/.local/bin
	install -D $(BUILD_DIR)/$(BINARY_NAME) ~/.local/bin/$(BINARY_NAME)
	@echo "Installation to user's local bin directory complete!"
	@echo "Don't forget to add \$$HOME/.local/bin to \$$PATH:"
	@echo "  export PATH=\"\$$HOME/.local/bin:\$$PATH\""

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  all            Builds and compresses the binary"
	@echo "  build          Builds the binary"
	@echo "  compress       Compresses the binary with UPX"
	@echo "  clean          Cleans the build directory"
	@echo "  install        Installs the binary to /usr/local/bin"
	@echo "  user-install   Installs the binary to ~/.local/bin"
	@echo "  help           Shows this help message"

.DEFAULT_GOAL := help