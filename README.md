# zlx

[![Go Report Card](https://goreportcard.com/badge/github.com/rlko/zlx)](https://goreportcard.com/report/github.com/rlko/zlx)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A flexible [Zipline](https://zipline.diced.sh/) upload client with configurable settings and command-line overrides.

## Features

- **Persistent Configuration** - Set your preferences once, use them everywhere
- **Runtime Flexibility** - Override any setting with command-line flags when needed
- **Cross-Platform** - Native support for Linux, macOS, and Windows (CMD/Powershell/WSL)
- **Distribution Packages** - Debian packages and tarballs for easy distribution
- **Flexible Output** - Echo to stdout with optional clipboard support

## Configuration Essentials

### Mandatory Settings
```bash
zlx config set servername "zipline.example.com"  # Defaults to https://
zlx config set servername "http://local.example" # Explicit scheme
zlx config set token "your_token_here"
```

### Optional Settings
```bash
zlx config set pathname "/custom/upload"  # Default: "/api/upload"
zlx config set upload.clipboard true      # Default: false (echo only)
zlx config set upload.max_views 10        # Default: 0 (unlimited)
zlx config set upload.original_name true  # Default: false
```

### View Current Settings
```bash
zlx config list  # Shows all current configuration values
```

## Clipboard Behavior

By default, zlx **echoes the URL to stdout** without copying to clipboard. The clipboard behavior can be controlled in two ways:

1. **Config File Setting**
   ```bash
   zlx config set upload.clipboard true  # Always copy to clipboard
   zlx up file.txt  # Now copies automatically
   ```

2. **Command Line Flags**
   ```bash
   zlx up file.txt -c    # Echoes AND copies to clipboard (overrides config)
   zlx up file.txt -n    # Only echoes, never copies (overrides config)
   ```

The `-n` flag will always disable clipboard copying, even if `upload.clipboard` is set to `true` in the config file.

### Clipboard Prerequisites

To use clipboard functionality, you need to have one of these clipboard utilities installed:

- **Linux**: `xsel`, `xclip`, or `wl-clipboard`
- **macOS**: `pbcopy` (included by default)
- **Windows**: Built-in clipboard support in CMD/PowerShell

If no clipboard utility is available, zlx will still echo the URL to stdout but will show an error message when attempting to copy to clipboard.

## Installation

### Prerequisites
- Go (1.24.2 or later)
- Make

### From Source (Makefile)

#### Build and install (system-wide)
```bash
make build
sudo make install
```

#### OR for local user installation
```bash
make build
make user-install
```

Available Makefile targets:
- `build` - Build the binary
- `install` - Install to `/usr/local/bin`
- `user-install` - Install to `~/.local/bin`
- `compress` - Compress binary with UPX
- `deb` - Create Debian package
- `tarball` - Create tar.gz archive
- `clean` - Remove build artifacts

## Usage Examples

### Basic Upload
```bash
zlx up file.txt
```

### Persistent Clipboard Setting
```bash
zlx config set upload.clipboard true  # Always copy to clipboard
zlx up file.txt  # Now copies automatically
```

## Flag Reference
| Flag | Short | Effect | Config Equivalent |
|------|-------|--------|--------------------|
| `--clipboard` | `-c` | Echoes URL AND copies to clipboard | `upload.clipboard` |
| `--no-clipboard` | `-n` | Disables copying to clipboard | |
| `--max-views` | `-m` | Sets maximum views | `upload.max_views` |
| `--original-name` | `-o` | Preserves filename | `upload.original_name` |

## License

MIT © [rlko](https://github.com/rlko)
