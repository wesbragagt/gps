# GPS

[![Go Reference](https://pkg.go.dev/badge/github.com/wesbragagt/gps.svg)](https://pkg.go.dev/github.com/wesbragagt/gps)
[![Go Report Card](https://goreportcard.com/badge/github.com/wesbragagt/gps)](https://goreportcard.com/report/github.com/wesbragagt/gps)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## TL;DR

**gps** maps your project structure with metadata optimized for AI agents.

```bash
gps              # That's it. See your project structure.
```

## Why gps?

When AI agents analyze codebases, they need efficient project overviews. **gps** provides:

- **Token-efficient output** - 50-70% reduction compared to JSON
- **Rich metadata** - File sizes, line counts, language detection
- **Project intelligence** - Entry points, configs, tests auto-detected
- **Multiple formats** - TOON (default), JSON, tree, flat
- **Smart filtering** - gitignore support, glob patterns

## Installation

```bash
go install github.com/wesbragagt/gps/cmd/gps@latest
# Or download from https://github.com/wesbragagt/gps/releases
```

## Quick Start

```bash
gps              # See your project structure
gps --summary    # Quick overview
gps -L 2         # Limit depth (like tree)
gps -f json      # JSON output
```

Output example:

```
project[gps]{
  type: go
  files: 47
  size: 284KB
  lines: 6231
}
root[47]{
  cmd[1]{
    gps[1]{
      main.go [1.8KB, 48L, go]
    }
  }
  internal[28]{
    scanner[8]{
      scanner.go [12KB, 356L, go]
      filter.go [4.2KB, 128L, go]
    }
  }
  go.mod [512B, 15L, go]
}
keyfiles{
  entry: cmd/gps/main.go
  config: go.mod
  tests: 12 files
}
```

## For AI Agents

gps is optimized for LLM token efficiency.

### Token Budget Guide

| Context Window | Recommended Command |
|----------------|---------------------|
| 4k tokens | `gps -L 2 --no-meta` |
| 8k tokens | `gps -L 3` |
| 16k+ tokens | `gps` (full output) |
| Unlimited | `gps -f json` (structured data) |

### Recommended Invocations

```bash
# Quick context (for any code task)
gps --summary && gps --entry-points

# Structure overview (for refactoring/navigation)
gps -L 3 -f toon

# Full analysis (for comprehensive tasks)
gps -f json | jq '.project.key_files'

# Minimal context (token-constrained)
gps -L 2 --no-meta --no-project-info
```

### Format Selection

| Format | Tokens | Best For |
|--------|--------|----------|
| `toon` | Baseline | Default, human + AI readable |
| `flat` | -72% | Max token savings |
| `json` | +385% | Structured processing |
| `tree` | +19% | Human presentations |

### Token Counting

```bash
gps --tokens    # Show token count
gps --compare   # Compare all formats
```

Modes: `--tokenizer approx` (fast) or `--tokenizer tiktoken` (accurate).

### System Prompt Snippet

```markdown
## Project Navigation
Use `gps` to understand project structure:
- `gps --summary` for quick overview
- `gps --entry-points` to find main files
- `gps -L 2 --no-meta` for token-efficient structure
```

## Usage

### Tree-Compatible Flags

```bash
gps -L 2                  # Limit depth
gps -d                    # Directories only
gps -a                    # Include hidden files
gps -I "node_modules"     # Exclude pattern
gps -P "*.go"             # Include only .go files
gps -e go,md              # Filter by extension
```

### AI-Optimized Flags

```bash
gps -f toon               # Token-optimized (default)
gps -f json               # Structured JSON
gps -f tree               # Traditional tree view
gps -f flat               # CSV-like format
gps --summary             # Project overview only
gps --entry-points        # Show entry points
gps --focus src/          # Focus on subdirectory
gps --no-meta             # Skip metadata
```

## Output Formats

### Quick Decision

```
Need max token savings? → flat (-72%)
Need structured data?   → json (+385%)
Human presentation?     → tree (+19%)
Default/unsure?         → toon (baseline)
```

### Comparison

| Format | Tokens | Use Case | Pros |
|--------|--------|----------|------|
| **toon** | 100% | Default, AI agents | Human + AI readable, compact |
| **flat** | 28% | Token-critical | Smallest output, CSV-like |
| **json** | 485% | Programmatic | Fully structured, queryable |
| **tree** | 119% | Presentations | Familiar format, colors |

### Examples

**TOON (default) - 538 tokens:**
```
project[myapp]{ type: go, files: 23 }
root[23]{ cmd[1]{ main.go [3KB, 89L, go] } }
```

**Flat - 150 tokens (-72%):**
```
path,size,lines,type
cmd/main.go,3072,89,go
go.mod,512,15,go
```

**JSON - 2608 tokens (+385%):**
```json
{"project":{"name":"myapp","tree":{"cmd":{"main.go":{"size":3072}}}}}
```

**Tree - 642 tokens (+19%):**
```
myapp/
├── cmd/main.go [3KB, 89L]
└── go.mod [512B, 15L]
```

## Special Modes

```bash
gps --summary           # Project overview
gps --entry-points      # Find entry points
gps --focus src/api     # Analyze subdirectory
```

## Common Patterns

### Understanding a New Codebase

```bash
gps --summary           # Quick overview
gps --entry-points      # Find where to start
gps -L 3                # See structure
```

### Preparing Context for AI

```bash
gps -L 2 --no-meta      # Minimal (4k budget)
gps -L 3 -f toon        # Standard (8k budget)
gps -f json             # Full (16k+ budget)
```

### Analyzing Specific Areas

```bash
gps --focus src/api -L 3    # Focus on subdirectory
gps -e go                   # Only Go files
gps -I "generated,*.pb.go"  # Exclude patterns
```

### CI/CD Integration

```bash
gps --summary > project-info.txt
gps -f json | jq '.project.stats.file_count'
```

### Documentation

```bash
gps -f tree > STRUCTURE.md
gps -f flat > files.csv
```

## Configuration

```yaml
# .gps.yaml
format: toon
depth: -1
exclude:
  - node_modules
  - dist
  - "*.log"
```

## Project Detection

| Language  | Config Files         | Entry Points              |
|-----------|---------------------|---------------------------|
| Go        | go.mod, go.sum      | main.go, cmd/*/main.go    |
| Node.js   | package.json        | index.js, server.js       |
| Python    | requirements.txt    | __main__.py, main.py      |
| Rust      | Cargo.toml          | src/main.rs               |
| Java      | pom.xml             | Main.java                 |

## Comparison with tree

| Feature                | tree | gps |
|------------------------|------|-----|
| Basic directory listing| ✅   | ✅  |
| Depth limit (-L)       | ✅   | ✅  |
| Exclude patterns (-I)  | ✅   | ✅  |
| Extension filtering    | ❌   | ✅  |
| File metadata          | ❌   | ✅  |
| Line counts            | ❌   | ✅  |
| Language detection     | ❌   | ✅  |
| Entry point detection  | ❌   | ✅  |
| gitignore support      | ❌   | ✅  |
| Multiple output formats| ❌   | ✅  |
| Token-optimized output | ❌   | ✅  |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup.

## License

MIT License - see [LICENSE](LICENSE).

## Acknowledgments

- Inspired by `tree`
- Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
