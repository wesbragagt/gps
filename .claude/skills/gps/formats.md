# GPS Output Formats Guide

Choose the right format for your use case. GPS provides 4 output formats optimized for different scenarios.

## Quick Decision

```
┌─────────────────────────────────────────────────────────────┐
│                    What's your goal?                        │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│  AI Context?  │   │  Scripting?   │   │ Human Reading?│
└───────────────┘   └───────────────┘   └───────────────┘
        │                   │                   │
        ▼                   ▼                   │
   ┌─────────┐         ┌─────────┐              │
   │  TOON   │         │  JSON   │              │
   │(default)│         │         │              │
   └─────────┘         └─────────┘              │
        │                   │                   │
        │                   │          ┌────────┴────────┐
        │                   │          ▼                 ▼
        │                   │   ┌─────────────┐   ┌─────────────┐
        │                   │   │    Tree     │   │    Flat     │
        │                   │   │(visual/pretty)│  │(CSV/analysis)│
        │                   │   └─────────────┘   └─────────────┘
        │                   │
        └───────────────────┴───────────────────────────────┐
                                                            ▼
                                                    ┌─────────────┐
                                                    │  Summary    │
                                                    │(--summary)  │
                                                    └─────────────┘
```

---

## Format Overview

| Format | Best For | Token Efficiency | Human Readable | Machine Parseable |
|--------|----------|------------------|----------------|-------------------|
| **TOON** | AI agents, LLM context | ⭐⭐⭐⭐⭐ 60% savings | ✅ Yes | ⚠️ Custom parsing |
| **JSON** | Scripts, CI/CD, tooling | ⭐⭐ Baseline | ✅ Yes | ✅ Native support |
| **Tree** | Terminal, documentation | ⭐⭐⭐ 48% savings | ✅✅ Very readable | ❌ Visual format |
| **Flat** | Spreadsheets, analysis | ⭐⭐⭐⭐ 55% savings | ✅ Yes | ✅ CSV tools |

### Command Reference

```bash
gps                    # TOON (default)
gps -f toon            # Explicit TOON
gps -f json            # JSON output
gps -f tree            # Tree view
gps -f flat            # Flat/CSV format
gps --summary          # Quick overview only
```

---

## TOON Format

**Token-Optimized Object Notation** - Default format for AI efficiency.

### When to Use ✅

- AI/LLM context windows
- Copying to chat interfaces
- Documentation with examples
- Quick terminal overview
- Token budget optimization

### When NOT to Use ❌

- Programmatic processing (use JSON)
- CI/CD pipelines (use JSON)
- Spreadsheet import (use Flat)
- Visual documentation (use Tree)

### Example Output

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
      flat.go [4.2KB, 181L, go]
    }
    scanner[4]{
      scanner.go [7.1KB, 274L, go]
      filter.go [3.2KB, 112L, go]
    }
  }
  go.mod [286B, 10L, go]
}
keyfiles{
  entry: cmd/gps/main.go
  config: go.mod
  tests: 8 files
  docs: README.md
}
```

### Token Efficiency Techniques

| Technique | Example | Savings |
|-----------|---------|---------|
| Compact brackets | `root[N]{}` vs `{"root": {"count": N}}` | ~40% |
| Inline metadata | `file.go [1KB, 50L, go]` | ~60% |
| Abbreviated units | `KB`, `L` (lines) | ~20% |
| No quotes on keys | `type:` vs `"type":` | ~15% |
| Aggregated counts | `[8 files]` | ~30% |

### Reading TOON

```
directory[count]{
  └─ subdirectory[count]{     ← nested directories
       file.ext [size, lines, type]  ← file with metadata
     }
}
```

---

## JSON Format

Structured output for programmatic processing.

### When to Use ✅

- Shell scripts with `jq`
- CI/CD pipelines
- Integration with other tools
- Complex queries
- API responses
- Configuration management

### When NOT to Use ❌

- AI context windows (use TOON)
- Quick terminal viewing (use Tree)
- Token-sensitive environments

### Example Output

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
    "tree": {
      "path": ".",
      "files": [...],
      "subdirs": [...]
    }
  }
}
```

### Processing Examples

```bash
# Extract entry points
gps -f json | jq '.project.key_files.entry_points[]'
# Output: "cmd/gps/main.go"

# Count files by type
gps -f json | jq '.project.stats.by_type'
# Output: {"go": 18, "md": 3, "yaml": 2}

# Get project type
gps -f json | jq -r '.project.type'
# Output: go

# Find all Go files
gps -f json | jq '[.. | .path? // empty | select(endswith(".go"))]'

# Get total size
gps -f json | jq '.project.stats.total_size'
# Output: 159744

# Extract test files
gps -f json | jq '.project.key_files.tests[]'

# Get file count
gps -f json | jq '.project.stats.file_count'
# Output: 23

# CI/CD: Check for specific file type
gps -f json | jq -e '.project.stats.by_type.go > 0' && echo "Go project"
```

---

## Tree Format

Traditional tree view with color coding and inline metadata.

### When to Use ✅

- Human terminal reading
- Visual project exploration
- Documentation/screenshots
- Quick size/line overview
- Onboarding new developers

### When NOT to Use ❌

- Programmatic processing (use JSON)
- AI context (use TOON)
- Data analysis (use Flat)

### Example Output

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
│   │   ├── toon.go [6.8KB, 227L]
│   │   └── tree.go [8.1KB, 308L]
│   └── scanner/
│       ├── filter.go [3.2KB, 112L]
│       └── scanner.go [7.1KB, 274L]
├── pkg/
│   └── types/
│       └── types.go [4.2KB, 169L]
├── go.mod [286B, 10L]
└── go.sum [1.1KB, 26L]

23 files, 156KB, 3420 lines
```

### Color Coding

| File Type | Color | Example |
|-----------|-------|---------|
| Go (.go) | Cyan | `main.go` |
| JavaScript/TypeScript | Yellow | `index.ts` |
| Python (.py) | Green | `app.py` |
| Rust (.rs) | Bold Yellow | `main.rs` |
| Java/Kotlin | Magenta | `Main.java` |
| Config files | Magenta | `package.json` |
| Documentation | White | `README.md` |
| Binary files | Dim | `binary` |
| Directories | Bold Blue | `src/` |

Disable colors: `NO_COLOR=1 gps -f tree`

---

## Flat Format

CSV-like format for data analysis and processing.

### When to Use ✅

- Spreadsheet import (Excel, Google Sheets)
- Shell scripting with awk/sed
- Data pipelines
- Sorting and filtering
- Statistical analysis
- Report generation

### When NOT to Use ❌

- Visual inspection (use Tree)
- AI context (use TOON)
- Hierarchical queries (use JSON)

### Example Output

```
path,size,lines,type,importance
cmd/gps/main.go,456,13,go,100
cmd/gps/root.go,11468,404,go,80
cmd/gps/version.go,512,21,go,70
go.mod,286,10,go,90
go.sum,1126,26,go,85
internal/detector/detector.go,15564,593,go,80
internal/formatter/flat.go,4300,181,go,70
internal/formatter/json.go,1433,59,go,70
internal/formatter/toon.go,6963,227,go,70
internal/formatter/tree.go,8294,308,go,70
```

### Field Reference

| Field | Description | Example |
|-------|-------------|---------|
| path | Relative file path | `cmd/gps/main.go` |
| size | File size in bytes | `456` |
| lines | Line count (0 for binary) | `13` |
| type | Detected file type | `go` |
| importance | Importance score (0-100) | `100` |

### Processing Examples

```bash
# Sort by size (largest first)
gps -f flat | tail -n +2 | sort -t',' -k2 -rn | head -5

# Find all Go files
gps -f flat | grep ",go,"

# Calculate total lines
gps -f flat | tail -n +2 | awk -F',' '{sum+=$3} END {print sum}'

# Find most important files (top 10)
gps -f flat | tail -n +2 | sort -t',' -k5 -rn | head -10

# Filter by minimum size (1KB)
gps -f flat | tail -n +2 | awk -F',' '$2 >= 1024'

# Get file count
gps -f flat | tail -n +2 | wc -l

# Average file size
gps -f flat | tail -n +2 | awk -F',' '{sum+=$2; count++} END {print int(sum/count)}'

# Find config files
gps -f flat | grep -E "(json|yaml|toml|mod),"

# Export to actual CSV
gps -f flat > project_files.csv
```

---

## Format Comparison Example

Same project (23 Go files) in all formats:

### TOON (~850 tokens)
```
project[gps]{
  type: go
  files: 23
  size: 156KB
}
root[31]{
  cmd[3]{
    gps[1]{ main.go [456B, 13L, go] }
  }
  internal[22]{ ... }
}
keyfiles{
  entry: cmd/gps/main.go
}
```

### JSON (~2,100 tokens)
```json
{
  "project": {
    "name": "gps",
    "type": "go",
    "stats": {
      "file_count": 23,
      "total_size": 159744
    },
    "tree": {
      "path": ".",
      "files": [...],
      "subdirs": [...]
    }
  }
}
```

### Tree (~1,100 tokens)
```
gps/
├── cmd/gps/main.go [456B, 13L]
├── internal/...
└── go.mod [286B, 10L]

23 files, 156KB
```

### Flat (~950 tokens)
```
path,size,lines,type,importance
cmd/gps/main.go,456,13,go,100
...
```

---

## Performance Characteristics

| Characteristic | TOON | JSON | Tree | Flat |
|----------------|------|------|------|------|
| **Token Efficiency** | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Parse Speed** | Medium | Fast | Slow | Fast |
| **Memory Usage** | Low | Medium | Low | Low |
| **Streaming** | No | Yes | No | Yes |
| **Schema** | None | Strict | None | CSV |
| **IDE Support** | None | Native | None | CSV tools |

### Token Reduction vs JSON

| Project Size | TOON | Tree | Flat |
|--------------|------|------|------|
| Small (<50 files) | 55-65% | 40-50% | 50-60% |
| Medium (50-200 files) | 60-70% | 45-55% | 55-65% |
| Large (200+ files) | 65-75% | 50-60% | 60-70% |

---

## Common Patterns

### AI Workflow

```bash
# 1. Quick project understanding
gps --summary

# 2. Find entry points
gps --entry-points

# 3. Get efficient context
gps -f toon -L 3

# 4. Focus on specific area
gps --focus src/api -L 2

# 5. Minimal context for structure only
gps -f toon --no-meta
```

### Scripting & Automation

```bash
# Check project type in CI
PROJECT_TYPE=$(gps -f json | jq -r '.project.type')

# Fail if no tests found
gps -f json | jq -e '.project.key_files.tests | length > 0' || exit 1

# Get file count threshold
FILE_COUNT=$(gps -f json | jq '.project.stats.file_count')
[ $FILE_COUNT -gt 1000 ] && echo "Large project warning"

# Generate report
gps -f flat > "report-$(date +%Y%m%d).csv"
```

### Documentation

```bash
# Include in README
gps -f tree -L 2 > STRUCTURE.md

# Visual overview for wiki
gps -f tree --no-meta > docs/project-layout.txt
```

---

## Quick Reference Card

```
┌────────────────────────────────────────────────────────────┐
│                    GPS FORMAT CHEATSHEET                   │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  AI/LLM Context     →  gps              (TOON default)    │
│  Scripting/CI/CD    →  gps -f json                         │
│  Human Terminal     →  gps -f tree                         │
│  Data Analysis      →  gps -f flat                         │
│  Quick Overview     →  gps --summary                       │
│  Find Entry Points  →  gps --entry-points                  │
│                                                            │
├────────────────────────────────────────────────────────────┤
│  TOKEN SAVINGS vs JSON                                     │
│  TOON: 60%  │  Tree: 48%  │  Flat: 55%                     │
└────────────────────────────────────────────────────────────┘
```

---

*Use `gps --compare` to see token counts for all formats on your project.*
