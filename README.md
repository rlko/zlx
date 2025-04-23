# zlx

[![Go Report Card](https://goreportcard.com/badge/github.com/rlko/zlx)](https://goreportcard.com/report/github.com/rlko/zlx)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`zlx` is a [Zipline](https://zipline.diced.sh/) upload client with smart defaults and runtime flexibility.

## The problem with the official generated shell script

1. Generate rigid configs that do **one thing**
2. Constantly edit files for simple changes

## How zlx Fixes This

- **No More Config Juggling**
  Set sane defaults once, override anytime with flags

- **True Per-Upload Control**
  Need max-views just this once? Add `-m 5` and move on

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

## Clipboard Behavior

By default, zlx **echoes the URL to stdout**. Use `-c` to **additionally** copy to clipboard:

```bash
zlx up file.txt       # Only echoes URL
zlx up file.txt -c    # Echoes AND copies to clipboard
```

## Installation

### From Source (Makefile)

#### Clone the repository
```bash
git clone https://github.com/rlko/zlx.git
cd zlx
```

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
- `build` - Compiles the binary
- `compress` - Compresses binary with UPX
- `install` - Installs to `/usr/local/bin`
- `user-install` - Installs to `~/.local/bin`
- `clean` - Removes build artifacts

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
| `--max-views` | `-m` | Sets maximum views | `upload.max_views` |
| `--original-name` | `-o` | Preserves filename | `upload.original_name` |

## Key Features
- **Explicit Output** - Always echoes URL to stdout
- **Additive Clipboard** - `-c` supplements (not replaces) echo behavior
- **No Silent Mode** - Always visible output with optional clipboard copy

## License

MIT © [rlko](https://github.com/rlko)
