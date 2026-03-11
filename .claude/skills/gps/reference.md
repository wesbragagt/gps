# GPS Reference Documentation

> Tree-like CLI tool optimized for AI agents with token-efficient output.

## Table of Contents

1. [Overview](#1-overview)
2. [Installation](#2-installation)
3. [Command Reference](#3-command-reference)
4. [Output Formats](#4-output-formats)
5. [Special Modes](#5-special-modes)
6. [Configuration](#6-configuration)
7. [Project Detection](#7-project-detection)
8. [Performance](#8-performance)

---

## 1. Overview

### What is GPS?

GPS (Go Project GPS) is a tree-like CLI tool designed for AI agents to map and traverse projects efficiently. If you know how to use `tree`, you already know how to use `gps`.

### Key Features

| Feature | Description |
|---------|-------------|
| Token-efficient output | 50-70% reduction compared to JSON |
| Rich metadata | File sizes, line counts, language detection |
| Project intelligence | Entry points, configs, tests auto-detected |
| Multiple formats | TOON (default), JSON, tree, flat |
| Smart filtering | gitignore support, glob patterns |

### Comparison with tree/find

| Feature | tree | find | gps |
|---------|------|------|-----|
| Basic directory listing | ✅ | ✅ | ✅ |
| Depth limit (-L) | ✅ | ✅ | ✅ |
| Directories only (-d) | ✅ | ✅ | ✅ |
| Include hidden (-a) | ✅ | ✅ | ✅ |
| Exclude patterns (-I) | ✅ | ✅ | ✅ |
| Include patterns (-P) | ✅ | ✅ | ✅ |
| Extension filtering (-e) | ❌ | ❌ | ✅ |
| File metadata | ❌ | ❌ | ✅ |
| Line counts | ❌ | ❌ | ✅ |
| Language detection | ❌ | ❌ | ✅ |
| Project type detection | ❌ | ❌ | ✅ |
| Entry point detection | ❌ | ❌ | ✅ |
| gitignore support | ❌ | ❌ | ✅ |
| Multiple output formats | ❌ | ❌ | ✅ |
| Token-optimized output | ❌ | ❌ | ✅ |
| Config file support | ❌ | ❌ | ✅ |

---

## 2. Installation

### From Binary (Recommended)

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

### Verify Installation

```bash
gps --version
gps --help
```

---

## 3. Command Reference

### Basic Usage

```bash
# Map current directory (TOON format by default)
gps

# Map a specific project
gps /path/to/project

# Map with depth limit
gps -L 2

# Map with JSON output
gps -f json
```

### Tree-Compatible Flags

| Flag | Description | Example |
|------|-------------|---------|
| `-L <n>` | Limit depth to n levels | `gps -L 2` |
| `-d` | Directories only (no files) | `gps -d` |
| `-a` | Include hidden files/directories | `gps -a` |
| `-I <pattern>` | Exclude pattern | `gps -I "node_modules"` |
| `-I <patterns>` | Multiple exclude patterns (comma-separated) | `gps -I "node_modules,dist,*.log"` |
| `-P <pattern>` | Include only matching files | `gps -P "*.go"` |
| `-P <patterns>` | Multiple include patterns | `gps -P "*.go,*.md"` |
| `-e <ext>` | Include only files with these extensions | `gps -e go,md,yaml` |
| `--exclude-ext <ext>` | Exclude files with these extensions | `gps --exclude-ext log,tmp,bak` |

**Examples:**

```bash
# Show only 2 levels deep
gps -L 2

# Show only directories
gps -d

# Include hidden files like .gitignore
gps -a

# Exclude common directories
gps -I "node_modules,vendor,dist,build"

# Show only Go files
gps -P "*.go"

# Show only source and config files
gps -P "*.go,*.yaml,*.json"

# Show only Go and Markdown files by extension
gps -e go,md

# Leading dot is optional
gps -e ".go,.md,.yaml"

# Exclude log and temp files
gps --exclude-ext log,tmp,bak

# Combine extension filter with exclude
gps -e go --exclude-ext test.go
```

### AI-Optimized Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--format <fmt>` | `-f` | Output format: toon, json, tree, flat |
| `--summary` | | High-level project overview only |
| `--entry-points` | | Show only detected entry points |
| `--focus <dir>` | | Focus on a specific subdirectory |
| `--no-meta` | | Skip metadata (faster, smaller output) |
| `--no-project-info` | | Skip project detection info |
| `--tokens` | `-t` | Show token count in output |
| `--compare` | | Compare token counts across all formats |
| `--tokenizer <type>` | | Tokenizer: approx (default), tiktoken |

**Examples:**

```bash
# Token-optimized output (default)
gps -f toon

# Structured JSON output
gps -f json

# Traditional tree view with colors
gps -f tree

# CSV-like flat format
gps -f flat

# Quick project overview
gps --summary

# Find entry points
gps --entry-points

# Focus on subdirectory
gps --focus src/

# Faster output, less metadata
gps --no-meta

# Show token count
gps --tokens

# Compare all formats
gps --compare

# Use GPT tokenizer
gps --tokenizer tiktoken
```

### Token Counting

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

**Tokenizer modes:**

| Mode | Description |
|------|-------------|
| `approx` | Fast, 4 chars ≈ 1 token (default) |
| `tiktoken` | Accurate GPT tokenizer (requires tiktoken-go) |

---

## 4. Output Formats

### TOON (Token-Optimized Object Notation)

**Default format** - Optimized for LLM token efficiency.

```bash
gps -f toon
```

**Structure:**

```
project[name]{
  type: <project-type>
  files: <count>
  size: <total-size>
  lines: <total-lines>
}
<directory>[<item-count>]{
  <file> [<size>, <lines>L, <type>]
  <subdirectory>[<count>]{ ... }
}
keyfiles{
  entry: <entry-point>
  config: <config-file>
  tests: <test-count> files
  docs: <doc-file>
}
```

**Example:**

```
project[gps]{
  type: go
  files: 23
  size: 156KB
  lines: 3420
}
root[31]{
  cmd[3]{
    gps[1]{
      main.go [456B, 13L, go]
    }
  }
  internal[22]{
    formatter[8]{
      toon.go [6.8KB, 227L, go]
      json.go [1.4KB, 59L, go]
      tree.go [8.1KB, 308L, go]
    }
    scanner[4]{
      scanner.go [7.1KB, 274L, go]
      filter.go [3.2KB, 112L, go]
    }
  }
  go.mod [286B, 10L, go]
  go.sum [1.1KB, 26L, go]
}
keyfiles{
  entry: cmd/gps/main.go
  config: go.mod
  tests: 8 files
  docs: README.md
}
```

**Token Efficiency Techniques:**

| Technique | Description |
|-----------|-------------|
| Compact syntax | Brackets instead of braces, no quotes on keys |
| Inline metadata | File info on single line with brackets |
| Aggregated counts | `[N files]` instead of listing each file |
| Abbreviated units | `KB`, `L` (lines) instead of verbose names |
| Minimal punctuation | No trailing commas, no colons in lists |

---

### JSON

Full structured output for programmatic processing.

```bash
gps -f json
```

**Structure:**

```json
{
  "project": {
    "name": "string",
    "type": "go|node|python|rust|java|mixed|other",
    "root": "string",
    "stats": {
      "file_count": 0,
      "total_size": 0,
      "total_lines": 0,
      "by_type": {}
    },
    "key_files": {
      "entry_points": [],
      "configs": [],
      "tests": [],
      "docs": []
    },
    "tree": {
      "path": "string",
      "files": [],
      "subdirs": []
    }
  }
}
```

**Example:**

```json
{
  "project": {
    "name": "gps",
    "type": "go",
    "root": "/home/user/gps",
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
      "entry_points": ["cmd/gps/main.go"],
      "configs": ["go.mod", "go.sum"],
      "tests": ["internal/formatter/json_test.go"],
      "docs": ["README.md"]
    },
    "tree": { ... }
  }
}
```

**Processing with jq:**

```bash
# Get entry points
gps -f json | jq '.project.key_files.entry_points[]'
# Output: "cmd/gps/main.go"

# Count files by type
gps -f json | jq '.project.stats.by_type'
# Output: {"go": 18, "md": 3, "yaml": 2}

# Get project type
gps -f json | jq -r '.project.type'
# Output: go
```

---

### Tree

Traditional tree view with syntax highlighting.

```bash
gps -f tree
```

**Example:**

```
gps/
├── cmd/
│   └── gps/
│       ├── main.go [456B, 13L]
│       ├── root.go [11.2KB, 404L]
│       └── version.go [512B, 21L]
├── internal/
│   ├── detector/
│   │   └── detector.go [15.2KB, 593L]
│   ├── formatter/
│   │   ├── flat.go [4.2KB, 181L]
│   │   ├── json.go [1.4KB, 59L]
│   │   └── toon.go [6.8KB, 227L]
│   └── scanner/
│       ├── filter.go [3.2KB, 112L]
│       └── scanner.go [7.1KB, 274L]
├── go.mod [286B, 10L]
└── go.sum [1.1KB, 26L]

23 files, 156KB, 3420 lines
```

**Color Coding:**

| File Type | Color |
|-----------|-------|
| Go (.go) | Cyan |
| JavaScript/TypeScript | Yellow |
| Python (.py) | Green |
| Rust (.rs) | Bold Yellow |
| Java/Kotlin | Magenta |
| Config files | Magenta |
| Documentation | White |
| Binary files | Dim |
| Directories | Bold Blue |

Disable colors with: `NO_COLOR=1 gps`

---

### Flat

CSV-like format for easy parsing.

```bash
gps -f flat
```

**Structure:**

```
path,size,lines,type,importance
<file-path>,<bytes>,<lines>,<type>,<score>
```

**Example:**

```
path,size,lines,type,importance
cmd/gps/main.go,456,13,go,100
cmd/gps/root.go,11468,404,go,80
cmd/gps/version.go,512,21,go,70
go.mod,286,10,go,90
go.sum,1126,26,go,85
internal/detector/detector.go,15564,593,go,80
internal/formatter/flat.go,4300,181,go,70
internal/scanner/scanner.go,7270,274,go,80
```

**Fields:**

| Field | Description |
|-------|-------------|
| path | Relative file path |
| size | File size in bytes |
| lines | Line count (0 for binary files) |
| type | Detected file type |
| importance | Importance score (0-100) |

**Processing Examples:**

```bash
# Sort by size (largest first)
gps -f flat | tail -n +2 | sort -t',' -k2 -rn | head -5

# Find all Go files
gps -f flat | grep ",go,"

# Calculate total lines
gps -f flat | tail -n +2 | awk -F',' '{sum+=$3} END {print sum}'

# Find most important files
gps -f flat | tail -n +2 | sort -t',' -k5 -rn | head -10
```

---

### Format Comparison

| Format | Tokens | Best For |
|--------|--------|----------|
| TOON | ~850 | AI context, quick overview |
| JSON | ~2,100 | Scripting, CI/CD, tooling |
| Tree | ~1,100 | Human reading, documentation |
| Flat | ~950 | Data analysis, spreadsheets |

**Recommendations:**

| Use Case | Format |
|----------|--------|
| AI agent context | TOON |
| Scripting/automation | JSON |
| Terminal exploration | Tree |
| Data analysis | Flat |
| Quick overview | TOON or --summary |

---

## 5. Special Modes

### Summary Mode

Get a quick project overview without full tree.

```bash
gps --summary
```

**Example:**

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

---

### Entry Points Mode

Find entry points quickly.

```bash
gps --entry-points
```

**Example:**

```
entry-points{
  main: cmd/gps/main.go
}
```

---

### Focus Mode

Analyze a specific subdirectory.

```bash
gps --focus internal/scanner
```

**With depth limit:**

```bash
gps --focus src/api -L 2
```

---

## 6. Configuration

Create a `.gps.yaml` file in your project root or home directory.

### Configuration File Structure

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
```

### Configuration Priority

1. Command-line flags (highest priority)
2. Project `.gps.yaml`
3. Home directory `~/.gps.yaml`
4. Built-in defaults (lowest priority)

### Example Configurations

**For Node.js projects:**

```yaml
format: toon
exclude:
  - node_modules
  - dist
  - build
  - .next
  - coverage
  - "*.min.js"
```

**For Go projects:**

```yaml
format: toon
exclude:
  - vendor
  - bin
  - "*.test"
```

**For Python projects:**

```yaml
format: toon
exclude:
  - __pycache__
  - .venv
  - venv
  - "*.pyc"
  - .pytest_cache
  - .mypy_cache
```

---

## 7. Project Detection

GPS automatically detects project types and key files.

### Supported Project Types

| Language | Config Files | Entry Points |
|----------|-------------|--------------|
| Go | go.mod, go.sum | main.go, cmd/*/main.go |
| Node.js | package.json | index.js, server.js, app.js |
| Python | requirements.txt, pyproject.toml | __main__.py, main.py, app.py |
| Rust | Cargo.toml | src/main.rs |
| Java | pom.xml, build.gradle | Main.java, Application.java |

### Detection Logic

1. **Project Type**: Detected by presence of config files
2. **Entry Points**: Language-specific main files
3. **Tests**: Files matching `*_test.go`, `*.test.js`, `test_*.py`
4. **Docs**: README.md, docs/ directory

### Mixed Projects

For monorepos or mixed-language projects:

```bash
gps -f json | jq '.project.type'
# Output: "mixed"
```

---

## 8. Performance

### Typical Scan Times

| Project Size | Files | Scan Time |
|--------------|-------|-----------|
| Small | <100 | <100ms |
| Medium | 100-1000 | 100-500ms |
| Large | 1000-10000 | 0.5-2s |
| Very Large | >10000 | 2-10s |

### Optimization Tips

**1. Limit Depth:**

```bash
gps -L 2  # Only 2 levels deep
```

**2. Skip Metadata:**

```bash
gps --no-meta  # Faster, smaller output
```

**3. Exclude Large Directories:**

```bash
gps -I "node_modules,vendor,build"
```

**4. Focus on Subdirectory:**

```bash
gps --focus src/  # Only scan src/
```

**5. Use Summary Mode:**

```bash
gps --summary  # Skip tree generation
```

**6. Pattern Filtering:**

```bash
gps -P "*.go"  # Only Go files
```

### Performance Comparison

| Command | Output Size | Speed |
|---------|-------------|-------|
| `gps` | ~2KB | Fast |
| `gps -f json` | ~10KB | Fast |
| `gps --no-meta` | ~1KB | Very Fast |
| `gps --summary` | ~500B | Very Fast |

---

## Quick Reference Card

```bash
# Basic
gps                    # Map current directory
gps /path/to/project   # Map specific project
gps -L 2               # Limit depth

# Formats
gps -f toon            # Token-optimized (default)
gps -f json            # Structured JSON
gps -f tree            # Traditional tree
gps -f flat            # CSV-like

# Filtering
gps -d                 # Directories only
gps -a                 # Include hidden
gps -I "node_modules"  # Exclude pattern
gps -P "*.go"          # Include pattern
gps -e go,md           # Include extensions
gps --exclude-ext log,tmp  # Exclude extensions

# Special Modes
gps --summary          # Quick overview
gps --entry-points     # Find entry points
gps --focus src/       # Focus subdirectory

# AI Optimization
gps --no-meta          # Skip metadata
gps --tokens           # Show token count
gps --compare          # Compare formats
```
