# Makefile for zlx installation

BINARY_NAME := zlx
PREFIX ?= /usr/local
BINDIR := $(PREFIX)/bin
MANDIR := $(PREFIX)/share/man/man1

.PHONY: install uninstall

install:
	@echo "Installing $(BINARY_NAME)..."
	install -Dm 0755 $(BINARY_NAME) $(BINDIR)/$(BINARY_NAME)
	install -Dm 0644 $(BINARY_NAME).1 $(MANDIR)/$(BINARY_NAME).1
	@echo "Installation complete!"

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f $(BINDIR)/$(BINARY_NAME)
	rm -f $(MANDIR)/$(BINARY_NAME).1
	@echo "Uninstallation complete!" 