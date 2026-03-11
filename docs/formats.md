# Output Formats

gps supports multiple output formats optimized for different use cases.

## TOON (Token-Optimized Object Notation)

**Default format** - Optimized for LLM token efficiency while remaining human-readable.

### Structure

```
project[name]{
  type: <project-type>
  files: <count>
  size: <total-size>
  lines: <total-lines>
}
<directory>[<item-count>]{
  <file> [<size>, <lines>L, <type>]
  <subdirectory>[<count>]{
    ...
  }
}
keyfiles{
  entry: <entry-point>
  config: <config-file>
  tests: <test-count> files
  docs: <doc-file>
}
```

### Example

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
      gitignore.go [2.8KB, 98L, go]
    }
    detector[2]{
      detector.go [15.2KB, 593L, go]
    }
    metadata[2]{
      extractor.go [4.5KB, 156L, go]
    }
    version[1]{
      version.go [1.1KB, 41L, go]
    }
  }
  pkg[6]{
    types[1]{
      types.go [4.2KB, 169L, go]
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

### Token Efficiency

TOON achieves 50-70% token reduction compared to JSON by:

1. **Compact syntax** - Brackets instead of braces, no quotes on keys
2. **Inline metadata** - File info on single line with brackets
3. **Aggregated counts** - `[N files]` instead of listing each file
4. **Abbreviated units** - `KB`, `L` (lines) instead of verbose names
5. **Minimal punctuation** - No trailing commas, no colons in lists

### When to Use

- AI agent context windows
- Quick terminal overviews
- Copying project structure to chat
- Documentation

---

## JSON

Full structured output for programmatic processing.

### Structure

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

### Example

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
      "tests": [
        "internal/formatter/json_test.go",
        "internal/formatter/tree_test.go",
        "internal/scanner/scanner_test.go"
      ],
      "docs": ["README.md"]
    },
    "tree": {
      "path": ".",
      "files": [
        {
          "path": "go.mod",
          "size": 286,
          "lines": 10,
          "type": "go",
          "mod_time": "2024-01-15T10:30:00Z",
          "is_binary": false,
          "importance": 90,
          "is_generated": false
        }
      ],
      "subdirs": [
        {
          "path": "cmd",
          "files": [],
          "subdirs": [...],
          "is_expanded": false
        }
      ],
      "is_expanded": false
    }
  }
}
```

### Processing Examples

```bash
# Get entry points
gps -f json | jq '.project.key_files.entry_points[]'
# Output: "cmd/gps/main.go"

# Count files by type
gps -f json | jq '.project.stats.by_type'
# Output: {"go": 18, "md": 3, "yaml": 2}

# Find largest files
gps -f json | jq '[.project.tree.files[].size] | max'
# Output: 159744

# Get project type
gps -f json | jq -r '.project.type'
# Output: go
```

### When to Use

- Scripting and automation
- CI/CD pipelines
- Integration with other tools
- Complex queries with jq

---

## Tree

Traditional tree view with syntax highlighting and inline metadata.

### Structure

```
<root>/
├── <directory>/
│   ├── <file> [metadata]
│   └── <file> [metadata]
├── <file> [metadata]
└── <file> [metadata]

<summary>
```

### Example

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
│       ├── gitignore.go [2.8KB, 98L]
│       └── scanner.go [7.1KB, 274L]
├── pkg/
│   └── types/
│       └── types.go [4.2KB, 169L]
├── go.mod [286B, 10L]
└── go.sum [1.1KB, 26L]

23 files, 156KB, 3420 lines
```

### Color Coding

The tree format uses ANSI colors for different file types:

| File Type       | Color   |
|-----------------|---------|
| Go (.go)        | Cyan    |
| JavaScript/TypeScript | Yellow |
| Python (.py)    | Green   |
| Rust (.rs)      | Bold Yellow |
| Java/Kotlin     | Magenta |
| Config files    | Magenta |
| Documentation   | White   |
| Binary files    | Dim     |
| Directories     | Bold Blue |

Colors can be disabled with `NO_COLOR=1` environment variable.

### When to Use

- Human reading in terminal
- Visual project exploration
- Screenshots and documentation
- Quick size/line overview

---

## Flat

CSV-like format for easy parsing, sorting, and analysis.

### Structure

```
path,size,lines,type,importance
<file-path>,<bytes>,<lines>,<type>,<score>
...
```

### Example

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
internal/scanner/filter.go,3276,112,go,70
internal/scanner/gitignore.go,2867,98,go,70
internal/scanner/scanner.go,7270,274,go,80
pkg/types/types.go,4300,169,go,70
```

### Processing Examples

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

### Fields

| Field       | Description                        |
|-------------|------------------------------------|
| path        | Relative file path                 |
| size        | File size in bytes                 |
| lines       | Line count (0 for binary files)    |
| type        | Detected file type                 |
| importance  | Importance score (0-100)           |

### When to Use

- Spreadsheet analysis
- Custom sorting/filtering
- Shell scripting
- Data pipelines

---

## Token Efficiency Comparison

Example: 23-file Go project

| Format  | Approximate Tokens | Reduction vs JSON |
|---------|-------------------|-------------------|
| JSON    | ~2,100            | -                 |
| TOON    | ~850              | 60%               |
| Tree    | ~1,100            | 48%               |
| Flat    | ~950              | 55%               |

**Recommendation**: Use TOON for AI context, JSON for tooling, Tree for humans.

## Choosing a Format

| Use Case                | Recommended Format |
|-------------------------|-------------------|
| AI agent context        | TOON              |
| Scripting/automation    | JSON              |
| Terminal exploration    | Tree              |
| Data analysis           | Flat              |
| Quick overview          | TOON or --summary |
| Find entry points       | --entry-points    |
| CI/CD pipelines         | JSON              |
| Documentation           | Tree or TOON      |
