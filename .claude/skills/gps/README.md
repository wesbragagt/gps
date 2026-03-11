# GPS Skill for AI Coding Assistants

> 🗺️ **Token-efficient project structure analysis for Claude Code, Cursor, GitHub Copilot, and Windsurf**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## What is GPS?

**GPS (Go Project Structure)** is a skill that enables AI assistants to map and traverse projects with **50-70% token reduction** compared to JSON output.

### Key Benefits

| Benefit | Impact |
|---------|--------|
| 🎯 **Token Efficiency** | 50-70% reduction vs JSON |
| 📊 **Rich Metadata** | File sizes, line counts, language detection |
| 🔍 **Project Intelligence** | Auto-detected entry points, configs, tests |
| 📁 **Multiple Formats** | TOON, JSON, tree, flat output |
| ⚡ **Easy Integration** | Works with 4 major AI assistants |

> **If you know `tree`, you already know GPS.**

---

## Why Use This Skill?

1. **🚀 Quick Codebase Exploration** - Understand any project in seconds
2. **💰 Token Budget Optimization** - Maximize AI context window efficiency
3. **🔍 Entry Point Discovery** - Find where to start reading code
4. **📝 Documentation Generation** - Create structure docs automatically
5. **🤖 Multi-Agent Support** - Same skill, multiple AI assistants

---

## Features

- ✅ **Multi-Agent Support** - Claude Code, Cursor, GitHub Copilot, Windsurf
- ✅ **Easy Integration** - Copy templates or use one-line setup
- ✅ **Comprehensive Docs** - Full command reference, examples, formats
- ✅ **Token Optimization** - Built-in token counting and comparison
- ✅ **Smart Filtering** - gitignore support, glob patterns
- ✅ **Project Detection** - Auto-detect Go, Node.js, Python, Rust, Java

---

## Quick Start

### 1️⃣ Install GPS

```bash
# Linux/macOS
curl -sL https://github.com/wesbragagt/gps/releases/latest/download/gps-$(uname -s)-$(uname -m) -o gps
chmod +x gps && sudo mv gps /usr/local/bin/

# Or with Go
go install github.com/wesbragagt/gps/cmd/gps@latest
```

### 2️⃣ Install Skill

```bash
# Claude Code (personal)
mkdir -p ~/.claude/skills/gps
cp -r /path/to/gps-skill/* ~/.claude/skills/gps/

# Claude Code (project)
mkdir -p .claude/skills/gps
cp -r /path/to/gps-skill/* .claude/skills/gps/
```

### 3️⃣ Use It

```
"Show me the project structure"
"Map this codebase"
"What are the entry points?"
```

---

## Documentation

| File | Description |
|------|-------------|
| 📖 **[INSTALL.md](INSTALL.md)** | Complete installation guide for all AI assistants |
| 📚 **[reference.md](reference.md)** | Full GPS command reference |
| 💡 **[examples.md](examples.md)** | Real-world use cases and patterns |
| 📋 **[formats.md](formats.md)** | Detailed output format guide |

---

## Common Commands

```bash
# Quick project overview
gps --summary

# Find entry points
gps --entry-points

# Map with depth limit
gps -L 2

# JSON for processing
gps -f json

# Focus on subdirectory
gps --focus src/api -L 3

# Token-efficient (default)
gps -f toon --no-meta

# Compare token savings
gps --compare

# Filter by extension
gps -e go,md,yaml

# Exclude extensions
gps --exclude-ext log,tmp,bak

# Filter by file type (pattern)
gps -P "*.go"

# Exclude patterns
gps -I "node_modules,dist"
```

---

## Supported AI Assistants

| Assistant | Integration | Setup |
|-----------|-------------|-------|
| **Claude Code** | Native skill system | `~/.claude/skills/gps/` |
| **Cursor** | Rules or .cursorrules | `.cursor/rules/gps.md` |
| **GitHub Copilot** | Workspace instructions | `.github/copilot-instructions.md` |
| **Windsurf** | IDE rules | `.windsurf/rules` |

---

## Examples

### Example 1: Quick Overview

```bash
gps --summary
```

**Output:**
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

### Example 2: Find Entry Points

```bash
gps --entry-points
```

**Output:**
```
entry-points{
  main: cmd/myapp/main.go
}
```

### Example 3: Focused Analysis

```bash
gps --focus internal/scanner -L 2
```

**Output:**
```
internal[18]{
  scanner[5]{
    scanner.go [8.2KB, 245L, go]
    filter.go [2.1KB, 67L, go]
    filter_test.go [1.8KB, 54L, go]
  }
}
```

---

## Use Cases

### 1. New Project Exploration

```bash
gps --summary && gps --entry-points
```

Perfect for understanding unfamiliar codebases quickly.

### 2. Code Review Preparation

```bash
gps -L 2 --focus src/
```

Get structure before diving into details.

### 3. Documentation Generation

```bash
gps -f tree > docs/structure.txt
```

Auto-generate project structure documentation.

### 4. AI Context Optimization

```bash
gps -f toon -L 3 --no-meta
```

Maximum token savings for AI context windows.

---

## Output Formats

| Format | Token Efficiency | Best For |
|--------|-----------------|----------|
| **TOON** | Baseline (best) | AI agent context (default) |
| **JSON** | +385% | Scripting, automation |
| **Tree** | +19% | Human reading |
| **Flat** | -72% | Data analysis, CSV import |

---

## Project Structure

```
gps-skill/
├── README.md              # This file
├── INSTALL.md             # Installation guide
├── SKILL.md               # Skill definition
├── reference.md           # Command reference
├── examples.md            # Use cases
├── formats.md             # Format guide
└── templates/
    ├── .cursorrules.template
    ├── copilot-instructions.md.template
    └── .windsurfrules.template
```

---

## Integration Templates

Use templates for quick setup:

### Cursor

```bash
cp templates/.cursorrules.template .cursorrules
```

### GitHub Copilot

```bash
mkdir -p .github
cp templates/copilot-instructions.md.template .github/copilot-instructions.md
```

### Windsurf

```bash
mkdir -p .windsurf
cp templates/.windsurfrules.template .windsurf/rules
```

---

## Token Efficiency

GPS provides significant token savings:

```bash
gps --compare
```

**Sample Output:**
```
Format          Tokens  Bytes   vs TOON
-------         ------  -----   -------
TOON              538    2149   baseline
JSON             2608   10411   +385%
JSON (compact)   1413    5645   +163%
Tree              642    2571   +19%
Flat              150     612   -72%
```

**For a typical project:**
- JSON: ~2,600 tokens
- TOON: ~540 tokens
- **Savings: ~2,060 tokens (79%)**

---

## Customization

Create `.gps.yaml` in your project:

```yaml
# Output settings
format: toon
depth: 3

# Exclude patterns
exclude:
  - node_modules
  - vendor
  - dist
  - "*.log"

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

# Metadata options
metadata:
  show_size: true
  show_lines: true
  show_type: true
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| **Command not found** | Install GPS: `go install github.com/wesbragagt/gps/cmd/gps@latest` |
| **Skill not activating** | Check file location, restart AI assistant, be explicit in prompts |
| **Too much output** | Add `-L 2` to limit depth, or `--no-meta` to skip metadata |

**More help:** See [INSTALL.md](INSTALL.md#troubleshooting)

---

## Contributing

Contributions welcome! 

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

---

## Requirements

| Requirement | Version |
|-------------|---------|
| **GPS CLI** | v1.0.0+ |
| **Claude Code** | Latest |
| **Cursor** | v0.40+ |
| **GitHub Copilot** | Latest |
| **Windsurf** | Latest |

---

## License

MIT License - see [LICENSE](https://opensource.org/licenses/MIT)

---

## Links

- 📦 **GPS Repository:** [github.com/wesbragagt/gps](https://github.com/wesbragagt/gps)
- 🐛 **Issues:** [github.com/wesbragagt/gps/issues](https://github.com/wesbragagt/gps/issues)
- 💬 **Discussions:** [github.com/wesbragagt/gps/discussions](https://github.com/wesbragagt/gps/discussions)

---

<p align="center">
  <strong>GPS Skill</strong> - Token-efficient project structure analysis for AI coding assistants
</p>
