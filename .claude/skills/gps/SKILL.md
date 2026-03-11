---
name: gps
description: |
  Use GPS (Go Project Structure) for token-efficient project structure analysis.
  Activates when users ask about: project structure, codebase organization, 
  file hierarchies, directory layouts, entry points, project layout, 
  file navigation, codebase overview, map project structure, project summary,
  key files, project architecture, folder structure, or when they need to 
  understand a codebase quickly. GPS provides 50-70% token reduction vs JSON.
version: 1.0.0
---

# GPS (Go Project Structure) Skill

Map and traverse projects efficiently with token-optimized output for AI agents.

## Overview

GPS is a tree-like CLI tool optimized for AI coding assistants. It provides:
- **50-70% token reduction** compared to JSON output
- **Rich metadata** - file sizes, line counts, language detection
- **Project intelligence** - auto-detected entry points, configs, tests
- **Multiple formats** - TOON (default), JSON, tree, flat

> **If you know `tree`, you already know GPS.**

## When to Use This Skill

Use GPS when the user:
- Asks to understand project structure or codebase organization
- Needs to find entry points or key files
- Wants a quick project overview or summary
- Is exploring an unfamiliar codebase
- Needs efficient context for AI token budgets
- Wants to map directory layouts or file hierarchies

## Quick Start

```bash
# Map current directory (TOON format by default)
gps

# Map specific project
gps /path/to/project

# Quick project overview
gps --summary

# Find entry points
gps --entry-points
```

## Common Workflows

### 1. Quick Project Overview
```bash
gps --summary
```
Output:
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

### 2. Find Entry Points
```bash
gps --entry-points
```
Output:
```
entry-points{
  main: cmd/myapp/main.go
}
```

### 3. Focus on Specific Area
```bash
# Analyze a subdirectory
gps --focus src/api -L 2
```

### 4. Token-Efficient Context
```bash
# Maximum token savings for AI context
gps -f toon -L 3 --no-meta
```

### 5. Structured Processing
```bash
# JSON for programmatic use
gps -f json | jq '.project.key_files.entry_points[]'
```

## Output Format Selection

| Use Case | Format | Command |
|----------|--------|---------|
| AI agent context (default) | TOON | `gps` or `gps -f toon` |
| Scripting/automation | JSON | `gps -f json` |
| Human reading | Tree | `gps -f tree` |
| Data analysis | Flat | `gps -f flat` |
| Quick overview | Summary | `gps --summary` |
| Find entry points | Entry | `gps --entry-points` |

## Command Reference

### Tree-Compatible Flags
```bash
gps -L 2                  # Limit depth to 2 levels
gps -d                    # Directories only
gps -a                    # Include hidden files
gps -I "node_modules"     # Exclude pattern
gps -I "node_modules,dist,*.log"  # Multiple excludes
gps -P "*.go"             # Include only Go files
gps -P "*.go,*.md"        # Include only Go and Markdown
gps -e go,md              # Include only .go and .md extensions
gps -e ".go,.ts"          # Leading dot optional
gps --exclude-ext log,tmp # Exclude .log and .tmp extensions
```

### AI-Optimized Flags
```bash
gps -f toon               # Token-optimized (default)
gps -f json               # Structured JSON
gps -f tree               # Traditional tree view
gps -f flat               # CSV-like format

gps --summary             # High-level overview only
gps --entry-points        # Show only entry points
gps --focus src/          # Focus on subdirectory
gps --no-meta             # Skip metadata (smaller output)
gps --tokens              # Show token count
gps --compare             # Compare token counts across formats
```

### Extension Filtering
```bash
gps -e go                 # Only Go files
gps -e go,md,yaml         # Multiple extensions
gps -e ".go,.ts"          # Leading dot optional
gps --exclude-ext log,tmp # Exclude .log and .tmp files
gps -e go --exclude-ext test.go  # Go files except tests
```

## TOON Format Reference

Default format optimized for LLM token efficiency:

```
project[name]{
  type: <project-type>
  files: <count>
  size: <total-size>
  lines: <total-lines>
}
root[N]{
  directory[M]{
    file.ext [size, NL, type]
  }
}
keyfiles{
  entry: <entry-point>
  config: <config-file>
  tests: <test-count> files
  docs: <doc-file>
}
```

**Notation:**
- `[N]` = item count in directory
- `[size, NL, type]` = file metadata (size, lines, type)
- `KB`, `L` = abbreviated units

## Best Practices

### For AI Context Windows
1. Use TOON format (default) - 50-70% token savings
2. Limit depth with `-L` for large projects
3. Use `--no-meta` when structure alone suffices
4. Use `--focus` for targeted analysis

### For Project Exploration
1. Start with `--summary` for quick overview
2. Use `--entry-points` to find starting points
3. Add `-L 2` or `-L 3` for manageable output
4. Filter by type with `-P "*.go"` etc.

### For CI/CD Integration
1. Use JSON format for parsing
2. Combine with `jq` for queries
3. Use exit codes for validation

## Project Detection

GPS auto-detects project types:

| Language | Config Files | Entry Points |
|----------|--------------|--------------|
| Go | go.mod, go.sum | main.go, cmd/*/main.go |
| Node.js | package.json | index.js, server.js, app.js |
| Python | requirements.txt, pyproject.toml | __main__.py, main.py, app.py |
| Rust | Cargo.toml | src/main.rs |
| Java | pom.xml, build.gradle | Main.java, Application.java |

## Reference Files

- **[reference.md](reference.md)** - Complete GPS command reference
- **[examples.md](examples.md)** - Real-world use cases and patterns
- **[formats.md](formats.md)** - Detailed format guide and selection

## Installation

GPS must be installed on the system:

```bash
# From binary (Linux/macOS)
curl -sL https://github.com/wesbragagt/gps/releases/latest/download/gps-$(uname -s)-$(uname -m) -o gps
chmod +x gps
sudo mv gps /usr/local/bin/

# With Go
go install github.com/wesbragagt/gps/cmd/gps@latest
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Command not found | Install GPS or check PATH |
| Too much output | Add `-L 2` to limit depth |
| Need specific files | Use `-P "*.go"` pattern |
| Want smaller output | Add `--no-meta` flag |

---

*GPS provides token-efficient project structure analysis for AI coding assistants.*
