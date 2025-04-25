# Makefile for zlx

PROJECT_NAME := zlx
BINARY_NAME := $(PROJECT_NAME)
BUILD_DIR := build
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")
LDFLAGS := -ldflags "-s -w"
UPX_OPTIONS := --best -q
DEB_DIR := $(BUILD_DIR)/$(PROJECT_NAME)-deb
VERSION := 0.0.1
ARCH := $(shell dpkg --print-architecture)
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
PACKAGE_BASENAME := $(BINARY_NAME)_v$(VERSION)_$(OS)_$(ARCH)

.PHONY: all build compress clean install user-install help deb tarball

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

deb: build compress
	@echo "Packaging .deb file..."
	rm -rf $(DEB_DIR)
	mkdir -p $(DEB_DIR)/DEBIAN
	mkdir -p $(DEB_DIR)/usr/local/bin
	install -m 0755 $(BUILD_DIR)/$(BINARY_NAME) $(DEB_DIR)/usr/local/bin/$(BINARY_NAME)
	@echo "Package: $(PROJECT_NAME)" > $(DEB_DIR)/DEBIAN/control
	@echo "Version: $(VERSION)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Section: utils" >> $(DEB_DIR)/DEBIAN/control
	@echo "Priority: optional" >> $(DEB_DIR)/DEBIAN/control
	@echo "Architecture: $(ARCH)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Maintainer: You <you@example.com>" >> $(DEB_DIR)/DEBIAN/control
	@echo "Description: $(PROJECT_NAME) built with Go" >> $(DEB_DIR)/DEBIAN/control
	chmod -R 0755 $(DEB_DIR)
	dpkg-deb --build $(DEB_DIR) $(BUILD_DIR)/$(PACKAGE_BASENAME).deb
	@echo "Debian package created: $(BUILD_DIR)/$(PACKAGE_BASENAME).deb"

tarball: build compress
	@echo "Creating tar.gz archive..."
	mkdir -p $(BUILD_DIR)/tarball
	cp $(BUILD_DIR)/$(BINARY_NAME) $(BUILD_DIR)/tarball/
	tar -czf $(BUILD_DIR)/$(PACKAGE_BASENAME).tar.gz -C $(BUILD_DIR)/tarball $(BINARY_NAME)
	rm -rf $(BUILD_DIR)/tarball
	@echo "Created: $(BUILD_DIR)/$(PACKAGE_BASENAME).tar.gz"

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  all            Builds and compresses the binary"
	@echo "  build          Builds the binary"
	@echo "  compress       Compresses the binary with UPX"
	@echo "  clean          Cleans the build directory"
	@echo "  install        Installs the binary to /usr/local/bin"
	@echo "  user-install   Installs the binary to ~/.local/bin"
	@echo "  deb            Builds a .deb package"
	@echo "  tarball        Creates a tar.gz archive"
	@echo "  help           Shows this help message"

.DEFAULT_GOAL := help