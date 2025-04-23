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
zlx config set pathname "/custom/api/upload"  # Default: "/api/upload"
zlx config set upload.clipboard true      # Default: false
zlx config set upload.max_views 10        # Default: 0 (unlimited)
zlx config set upload.original_name true  # Default: false
```

## Core Principle

```bash
# Persistent config sets safe defaults
zlx config set upload.clipboard true

# Flags provide per-upload control
zlx up file.txt                 # Uses clipboard (from config)
zlx up secret.txt -c=false      # Disables just for this upload
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
make compress
sudo make install
```

#### OR for local user installation
```bash
make build
make compress
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

### Custom One-Time Upload
```bash
zlx up sensitive.doc \
  -m 1 \             # Max 1 view
  -o \               # Keep original name
  -c=false           # Disable clipboard copy
```

## Flag Reference
| Flag | Short | Config Equivalent | Default |
|------|-------|--------------------|---------|
| `--max-views` | `-m` | `upload.max_views` | 0 |
| `--original-name` | `-o` | `upload.original_name` | false |
| `--clipboard` | `-c` | `upload.clipboard` | false |

## Key Features
- **Smart URL Handling**  
  `servername` auto-prepends https:// if missing
- **Config Hierarchy**  
  Flags > Persistent Config > Defaults
- **No Negative Flags (yet)**  
  Use `-c=false` instead of non-existent `--no-clipboard`

## License

MIT © [rlko](https://github.com/rlko)