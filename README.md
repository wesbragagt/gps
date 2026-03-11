# gps (go-project-gps)

[![Go Reference](https://pkg.go.dev/badge/github.com/wesbragagt/gps.svg)](https://pkg.go.dev/github.com/wesbragagt/gps)
[![Go Report Card](https://goreportcard.com/badge/github.com/wesbragagt/gps)](https://goreportcard.com/report/github.com/wesbragagt/gps)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A tree-like CLI tool optimized for AI agents to map and traverse projects efficiently.

**If you know how to use `tree`, you already know how to use `gps`.**

## Why gps?

When AI agents analyze codebases, they need efficient project overviews. Traditional tools like `tree` or `find` produce verbose output. JSON is structured but token-heavy. **gps** solves this with:

- **Token-efficient output** - 50-70% reduction compared to JSON
- **Rich metadata** - File sizes, line counts, language detection
- **Project intelligence** - Entry points, configs, tests auto-detected
- **Multiple formats** - TOON (default), JSON, tree, flat
- **Smart filtering** - gitignore support, glob patterns

## Installation

### From Binary

Download the latest release for your platform from the [releases page](https://github.com/wesbragagt/gps/releases).

```bash
# Linux/macOS
curl -sL https://github.com/wesbragagt/gps/releases/latest/download/gps-$(uname -s)-$(uname -m) -o gps
chmod +x gps
sudo mv gps /usr/local/bin/
```

### From Source

```bash
git clone https://github.com/wesbragagt/gps.git
cd gps
go build ./cmd/gps
./gps --help
```

### With Go

```bash
go install github.com/wesbragagt/gps/cmd/gps@latest
```

## Quick Start

```bash
# Map current directory (TOON format by default)
gps

# Map a specific project
gps /path/to/project

# Limit depth (like tree -L)
gps -L 2

# JSON output for structured processing
gps -f json

# Just the project summary
gps --summary
```

## Usage

### Tree-Compatible Flags

If you're familiar with `tree`, these flags work the same way:

```bash
gps -L 2                  # Limit depth to 2 levels
gps -d                    # Directories only (no files)
gps -a                    # Include hidden files/directories
gps -I "node_modules"     # Exclude pattern
gps -I "node_modules,dist,*.log"  # Multiple exclude patterns
gps -P "*.go"             # Include only Go files
gps -P "*.go,*.md"        # Include only Go and Markdown files
gps -e go,md              # Include only .go and .md extensions
gps -e ".go,.ts"          # Leading dot optional
gps --exclude-ext log,tmp # Exclude .log and .tmp extensions
```

### AI-Optimized Flags

Flags designed for AI agent workflows:

```bash
gps -f toon               # Token-optimized output (default)
gps -f json               # Structured JSON output
gps -f tree               # Traditional tree view with colors
gps -f flat               # CSV-like flat format

gps --summary             # High-level project overview only
gps --entry-points        # Show only detected entry points
gps --focus src/          # Focus on a specific subdirectory

gps --no-meta             # Skip metadata (faster, smaller output)
gps --no-project-info     # Skip project detection info

# Token counting (NEW!)
gps --tokens              # Show token count in output
gps -t                    # Short form
gps --compare             # Compare token counts across all formats
gps --tokenizer tiktoken  # Use tiktoken (GPT tokenizer)
```

### Token Counting

Understand token usage for AI agents:

```bash
# Show token count with output
gps --tokens

# Compare all formats
gps --compare

# Example output:
# Format          Tokens  Bytes   vs TOON
# -------         ------  -----   -------
# TOON              538    2149   baseline
# JSON             2608   10411   +385%
# JSON (compact)   1413    5645   +163%
# Tree              642    2571   +19%
# Flat              150     612   -72%
```

**Token counting modes:**
- `--tokenizer approx` (default): Fast, 4 chars ≈ 1 token
- `--tokenizer tiktoken`: Accurate GPT tokenizer (requires tiktoken-go)

## Output Formats

### TOON (Token-Optimized Object Notation) - Default

Compact, human-readable format optimized for LLM token efficiency:

```
project[myapp]{
  type: go
  files: 23
  size: 156KB
  lines: 3420
}
root[31]{
  cmd[3]{
    myapp[1]{
      main.go [3.2KB, 89L, go]
    }
  }
  internal[18]{
    scanner[5]{
      scanner.go [8.2KB, 245L, go]
      filter.go [2.1KB, 67L, go]
    }
  }
  go.mod [512B, 15L, go]
  README.md [4.5KB, 142L, md]
}
keyfiles{
  entry: cmd/myapp/main.go
  config: go.mod
  tests: 5 files
  docs: README.md
}
```

### JSON

Full structured output for programmatic processing:

```bash
gps -f json
```

```json
{
  "project": {
    "name": "myapp",
    "type": "go",
    "root": "/home/user/myapp",
    "stats": {
      "file_count": 23,
      "total_size": 159744,
      "total_lines": 3420,
      "by_type": {
        "go": 18,
        "md": 3,
        "yaml": 2
      }
    },
    "key_files": {
      "entry_points": ["cmd/myapp/main.go"],
      "configs": ["go.mod", "go.sum"],
      "tests": ["scanner_test.go", "formatter_test.go"],
      "docs": ["README.md"]
    },
    "tree": { ... }
  }
}
```

### Tree

Traditional tree view with syntax highlighting:

```bash
gps -f tree
```

```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go [3.2KB, 89L]
├── internal/
│   ├── scanner/
│   │   ├── scanner.go [8.2KB, 245L]
│   │   └── filter.go [2.1KB, 67L]
│   └── formatter/
│       └── toon.go [4.1KB, 128L]
├── go.mod [512B, 15L]
└── README.md [4.5KB, 142L]

23 files, 156KB, 3420 lines
```

### Flat

CSV-like format for easy parsing and analysis:

```bash
gps -f flat
```

```
path,size,lines,type,importance
cmd/myapp/main.go,3276,89,go,100
go.mod,512,15,go,90
internal/scanner/scanner.go,8396,245,go,80
README.md,4608,142,md,40
```

## Special Modes

### Summary Mode

Get a quick project overview:

```bash
gps --summary
```

```
project[myapp]{
  type: go
  files: 23
  size: 156KB
  entry: cmd/myapp/main.go
  tests: 5 files
  docs: README.md
}
```

### Entry Points Mode

Find entry points quickly:

```bash
gps --entry-points
```

```
entry-points{
  main: cmd/myapp/main.go
}
```

### Focus Mode

Analyze a specific subdirectory:

```bash
gps --focus internal/scanner
```

## Configuration

Create a `.gps.yaml` file in your project or home directory:

```yaml
# Output format: toon, json, tree, flat
format: toon

# Max depth (-1 for unlimited)
depth: -1

# Include hidden files
all: false

# Include metadata
meta: true

# Default exclude patterns
exclude:
  - node_modules
  - vendor
  - dist
  - build
  - "*.log"
  - "*.tmp"

# Include only specific extensions
ext:
  - .go
  - .md
  - .yaml

# Exclude specific extensions
exclude-ext:
  - .log
  - .tmp
  - .bak
```

## Project Detection

gps automatically detects project types and key files:

| Language  | Config Files                    | Entry Points                    |
|-----------|--------------------------------|---------------------------------|
| Go        | go.mod, go.sum                 | main.go, cmd/*/main.go          |
| Node.js   | package.json                   | index.js, server.js, app.js     |
| Python    | requirements.txt, pyproject.toml | __main__.py, main.py, app.py  |
| Rust      | Cargo.toml                     | src/main.rs                     |
| Java      | pom.xml, build.gradle          | Main.java, Application.java     |

## Use Cases for AI Agents

### Quick Project Overview

```bash
gps --summary
```
Perfect for understanding a codebase in one command.

### Find Entry Points

```bash
gps --entry-points
```
Quickly identify where to start reading code.

### Analyze Specific Areas

```bash
gps --focus src/api -L 2
```
Focus analysis on specific subdirectories.

### Token-Efficient Context

```bash
gps -f toon -L 3 --no-meta
```
Get structure without metadata for maximum token savings.

### Structured Processing

```bash
gps -f json | jq '.project.key_files.entry_points'
```
Parse output programmatically with JSON.

## Comparison with tree

| Feature                  | tree        | gps           |
|--------------------------|-------------|---------------|
| Basic directory listing  | ✅          | ✅            |
| Depth limit (-L)         | ✅          | ✅            |
| Directories only (-d)    | ✅          | ✅            |
| Include hidden (-a)      | ✅          | ✅            |
| Exclude patterns (-I)    | ✅          | ✅            |
| Include patterns (-P)    | ✅          | ✅            |
| Extension filtering      | ❌          | ✅            |
| File metadata            | ❌          | ✅            |
| Line counts              | ❌          | ✅            |
| Language detection       | ❌          | ✅            |
| Project type detection   | ❌          | ✅            |
| Entry point detection    | ❌          | ✅            |
| gitignore support        | ❌          | ✅            |
| Multiple output formats  | ❌          | ✅            |
| Token-optimized output   | ❌          | ✅            |
| Config file support      | ❌          | ✅            |

## Development

### Prerequisites

- Go 1.23 or later
- Make (optional, for build commands)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/wesbragagt/gps.git
cd gps

# Build using Make
make build        # Build for current platform
make release      # Build for all platforms

# Or build directly with Go
go build ./cmd/gps
```

### Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build binary for current platform |
| `make test` | Run tests with race detection |
| `make coverage` | Generate HTML coverage report |
| `make lint` | Run golangci-lint |
| `make fmt` | Format Go source files |
| `make vet` | Run go vet |
| `make check` | Run fmt, vet, and test |
| `make release` | Build binaries for all platforms |
| `make archives` | Create release archives |
| `make install` | Install to GOPATH/bin |
| `make docker` | Build Docker image |
| `make clean` | Remove build artifacts |
| `make version` | Display version info |

### Testing

```bash
# Run all tests
make test

# Or with Go directly
go test -v -race ./...

# Generate coverage report
make coverage
```

### Running

```bash
./bin/gps --help
./bin/gps .
```

### Releasing

The release process is automated via GitHub Actions:

1. **Create a tag:**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions will:**
   - Run tests
   - Build binaries for Linux, macOS, Windows (amd64, arm64)
   - Create GitHub release with archives and checksums

**Alternative: Using GoReleaser**
```bash
# Install goreleaser: go install github.com/goreleaser/goreleaser/v2@latest
GITHUB_TOKEN=$(your_token) goreleaser release --clean
```

### Docker

```bash
# Build image
make docker

# Run container
docker run --rm gps:latest --help

# With volume mount
docker run --rm -v $(pwd):/data gps:latest /data
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on:

- Code of conduct
- Development setup
- Submitting pull requests
- Coding standards

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the classic `tree` command
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Configuration powered by [Viper](https://github.com/spf13/viper)
