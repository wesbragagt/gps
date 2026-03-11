# GPS Skill Installation Guide

Complete setup guide for GPS (Go Project Structure) skill across AI coding assistants.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Step 1: Install GPS](#step-1-install-gps)
- [Step 2: Install GPS Skill](#step-2-install-gps-skill)
  - [Claude Code](#claude-code)
  - [Cursor](#cursor)
  - [GitHub Copilot](#github-copilot)
  - [Windsurf](#windsurf)
- [Step 3: Configure GPS (Optional)](#step-3-configure-gps-optional)
- [Step 4: Verify Installation](#step-4-verify-installation)
- [Troubleshooting](#troubleshooting)
- [Quick Start Examples](#quick-start-examples)
- [Uninstallation](#uninstallation)
- [Getting Help](#getting-help)

---

## Prerequisites

Before installing GPS Skill, ensure you have:

| Requirement | Description | Check Command |
|-------------|-------------|---------------|
| **GPS Tool** | Go Project Structure CLI | `gps --version` |
| **AI Assistant** | Claude Code, Cursor, Copilot, or Windsurf | - |
| **Git Repository** | A project you want to analyze | `git status` |

---

## Step 1: Install GPS

Choose your installation method:

### Option A: From Binary (Recommended)

**Linux (amd64/arm64):**
```bash
# Download latest release
curl -sL https://github.com/wesbragagt/gps/releases/latest/download/gps-$(uname -s)-$(uname -m) -o gps

# Make executable
chmod +x gps

# Move to PATH
sudo mv gps /usr/local/bin/

# Verify
gps --version
```

**macOS (Intel/Apple Silicon):**
```bash
# Download latest release
curl -sL https://github.com/wesbragagt/gps/releases/latest/download/gps-$(uname -s)-$(uname -m) -o gps

# Make executable
chmod +x gps

# Move to PATH
sudo mv gps /usr/local/bin/

# Verify
gps --version
```

**Windows (PowerShell):**
```powershell
# Download latest release
Invoke-WebRequest -Uri "https://github.com/wesbragagt/gps/releases/latest/download/gps-windows-amd64.exe" -OutFile "gps.exe"

# Move to PATH (adjust as needed)
Move-Item gps.exe "C:\Program Files\gps\"

# Verify
gps --version
```

### Option B: From Source

```bash
# Clone repository
git clone https://github.com/wesbragagt/gps.git
cd gps

# Build
go build -o gps ./cmd/gps

# Install to PATH
sudo mv gps /usr/local/bin/

# Verify
gps --version
```

### Option C: With Go

```bash
# Install directly
go install github.com/wesbragagt/gps/cmd/gps@latest

# Verify (ensure Go bin is in PATH)
gps --version

# If command not found, add Go bin to PATH:
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Verify GPS Installation

```bash
# Check version
gps --version

# Test on current directory
gps --summary
```

Expected output:
```
project[your-project]{
  type: <detected-type>
  files: <count>
  size: <total-size>
}
```

---

## Step 2: Install GPS Skill

Install the GPS skill for your AI assistant:

### Claude Code

Claude Code supports both **personal** (global) and **project-level** skill installation.

#### Personal Use (Global)

Install for all your projects:

```bash
# Create skills directory if needed
mkdir -p ~/.claude/skills/gps

# Copy skill files
cp -r /path/to/gps-skill/* ~/.claude/skills/gps/

# Verify
ls ~/.claude/skills/gps/
```

**File structure:**
```
~/.claude/skills/
└── gps/
    ├── SKILL.md
    ├── reference.md
    ├── examples.md
    ├── formats.md
    └── templates/
```

#### Project Use

Install for a specific project:

```bash
# Create project skills directory
mkdir -p .claude/skills/gps

# Copy skill files
cp -r /path/to/gps-skill/* .claude/skills/gps/

# Verify
ls .claude/skills/gps/
```

**File structure:**
```
your-project/
└── .claude/
    └── skills/
        └── gps/
            ├── SKILL.md
            ├── reference.md
            ├── examples.md
            └── formats.md
```

#### Activate in Claude Code

The skill activates automatically when you mention project structure-related tasks. Try:

```
"Show me the project structure"
"Map out this codebase"
"What are the entry points?"
```

---

### Cursor

Cursor supports both **legacy** and **new** skill formats.

#### New Format (Recommended)

```bash
# Create rules directory
mkdir -p .cursor/rules

# Copy skill as a rule
cp /path/to/gps-skill/SKILL.md .cursor/rules/gps.md

# Verify
ls .cursor/rules/
```

#### Legacy Format

```bash
# Create .cursorrules file
cat > .cursorrules << 'EOF'
# GPS Skill Integration
When analyzing project structure, use the GPS CLI tool.

## Commands
- `gps` - Map current directory
- `gps --summary` - Quick overview
- `gps --entry-points` - Find entry points
- `gps -L 2` - Limit depth

## Reference
See .claude/skills/gps/SKILL.md for full documentation.
EOF

# Copy full skill for reference
mkdir -p .claude/skills/gps
cp -r /path/to/gps-skill/* .claude/skills/gps/
```

#### Activate in Cursor

Mention project structure in your prompts:

```
@codebase What's the structure of this project?
Map the directory layout
Find all entry points in this codebase
```

---

### GitHub Copilot

Configure Copilot to use GPS via workspace instructions.

#### Setup

```bash
# Create GitHub Copilot instructions directory
mkdir -p .github

# Create instructions file
cat > .github/copilot-instructions.md << 'EOF'
# GPS Tool Integration

Use the GPS CLI tool for project structure analysis:

## Commands
- `gps` - Token-optimized structure (default)
- `gps --summary` - Quick project overview
- `gps --entry-points` - Find entry points
- `gps -f json` - JSON output for parsing
- `gps -L 2` - Limit depth to 2 levels

## When to Use
- Understanding project organization
- Finding entry points and key files
- Getting codebase overviews
- Mapping directory structures

## Reference
Full documentation: .claude/skills/gps/SKILL.md
EOF

# Copy skill files for reference
mkdir -p .claude/skills/gps
cp -r /path/to/gps-skill/* .claude/skills/gps/
```

**File structure:**
```
your-project/
├── .github/
│   └── copilot-instructions.md
└── .claude/
    └── skills/
        └── gps/
            └── ...
```

#### Activate in Copilot Chat

```
@workspace Show me the project structure
@workspace What are the main entry points?
@workspace Analyze the codebase organization
```

---

### Windsurf

Configure Windsurf via the IDE settings.

#### Setup

```bash
# Create Windsurf rules directory
mkdir -p .windsurf

# Create rules file
cat > .windsurf/rules << 'EOF'
# GPS Tool Integration

Use GPS CLI for project structure analysis:

## Primary Commands
- `gps` - Map project (TOON format, 50-70% token savings)
- `gps --summary` - Quick overview
- `gps --entry-points` - Find entry points
- `gps -L 2` - Limit depth

## Use Cases
- Project exploration
- Entry point discovery
- Codebase mapping
- Structure documentation

## Reference
Full docs: .claude/skills/gps/SKILL.md
EOF

# Copy skill files
mkdir -p .claude/skills/gps
cp -r /path/to/gps-skill/* .claude/skills/gps/
```

**File structure:**
```
your-project/
├── .windsurf/
│   └── rules
└── .claude/
    └── skills/
        └── gps/
            └── ...
```

#### Activate in Windsurf

Use natural language prompts:

```
What's the structure of this project?
Map out the codebase
Show me the entry points
```

---

## Step 3: Configure GPS (Optional)

Create a `.gps.yaml` configuration file for project-specific settings:

```bash
# Create config in project root
cat > .gps.yaml << 'EOF'
# GPS Configuration
output:
  format: toon      # toon, json, tree, flat
  depth: 3          # Default depth limit

exclude:
  - node_modules
  - ".git"
  - dist
  - build
  - "*.log"

include:
  - "*.go"
  - "*.md"
  - "*.yaml"
  - "*.json"

metadata:
  show_size: true
  show_lines: true
  show_type: true
EOF
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `output.format` | Default output format | `toon` |
| `output.depth` | Default depth limit | unlimited |
| `exclude` | Patterns to exclude | `[".git"]` |
| `include` | Patterns to include only | all files |
| `metadata.show_size` | Show file sizes | `true` |
| `metadata.show_lines` | Show line counts | `true` |
| `metadata.show_type` | Show file types | `true` |

---

## Step 4: Verify Installation

### Verify GPS Tool

```bash
# Check version
gps --version

# Test on current directory
gps --summary

# Test entry points detection
gps --entry-points
```

### Verify Claude Code

```bash
# Check skill exists
ls ~/.claude/skills/gps/SKILL.md   # Personal
ls .claude/skills/gps/SKILL.md     # Project
```

In Claude Code, try:
```
"Use GPS to show me the project structure"
```

### Verify Cursor

```bash
# Check rules exist
ls .cursor/rules/gps.md      # New format
cat .cursorrules             # Legacy format
```

In Cursor, try:
```
"Map this project using GPS"
```

### Verify GitHub Copilot

```bash
# Check instructions exist
cat .github/copilot-instructions.md
```

In Copilot Chat:
```
@workspace Use GPS to analyze the project structure
```

### Verify Windsurf

```bash
# Check rules exist
cat .windsurf/rules
```

In Windsurf:
```
"Show me the project structure using GPS"
```

---

## Troubleshooting

### GPS Not Found

**Error:** `gps: command not found`

**Solutions:**

```bash
# Check if GPS is installed
which gps

# Check PATH
echo $PATH

# If installed via Go, ensure Go bin is in PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Add to shell config for persistence
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

### Skill Not Activating

**Issue:** AI assistant doesn't use GPS

**Solutions:**

1. **Check file locations:**
   ```bash
   # Claude Code
   ls ~/.claude/skills/gps/SKILL.md
   
   # Cursor
   ls .cursor/rules/gps.md
   
   # Copilot
   ls .github/copilot-instructions.md
   
   # Windsurf
   ls .windsurf/rules
   ```

2. **Restart the AI assistant** after installation

3. **Be explicit in prompts:**
   ```
   "Use the gps command to show the project structure"
   ```

### Too Much Output

**Issue:** GPS output is too large

**Solutions:**

```bash
# Limit depth
gps -L 2

# No metadata
gps --no-meta

# Focus on subdirectory
gps --focus src/

# Combine options
gps -L 2 --no-meta --focus src/
```

### Permission Denied

**Error:** `permission denied: gps`

**Solutions:**

```bash
# Make executable
chmod +x /usr/local/bin/gps

# Or reinstall with sudo
sudo curl -sL https://github.com/wesbragagt/gps/releases/latest/download/gps-$(uname -s)-$(uname -m) -o /usr/local/bin/gps
sudo chmod +x /usr/local/bin/gps
```

### Wrong Architecture

**Error:** `cannot execute binary file`

**Solutions:**

```bash
# Check your architecture
uname -m

# Download correct binary
# amd64: x86_64, x64
# arm64: aarch64, arm64

# Or build from source
go install github.com/wesbragagt/gps/cmd/gps@latest
```

---

## Quick Start Examples

### Example 1: New Project Exploration

```bash
# Clone a new project
git clone https://github.com/example/project.git
cd project

# Get quick overview
gps --summary

# Find where to start
gps --entry-points

# Map structure with depth limit
gps -L 3
```

### Example 2: AI-Assisted Code Review

In your AI assistant:

```
I need to review this project. First, use GPS to map the structure 
with depth 2, then identify the main entry points.
```

```bash
# AI should run:
gps -L 2
gps --entry-points
```

### Example 3: Documentation Generation

```bash
# Generate structure for docs
gps -f tree > docs/structure.txt

# Get JSON for processing
gps -f json | jq '.project' > docs/project-info.json

# Compare token savings
gps --compare
```

---

## Uninstallation

### Remove GPS Tool

```bash
# Remove binary
sudo rm /usr/local/bin/gps

# Or if installed via Go
rm $(go env GOPATH)/bin/gps
```

### Remove Claude Code Skill

```bash
# Personal
rm -rf ~/.claude/skills/gps

# Project
rm -rf .claude/skills/gps
```

### Remove Cursor Integration

```bash
# New format
rm .cursor/rules/gps.md

# Legacy format
rm .cursorrules
```

### Remove GitHub Copilot Integration

```bash
rm .github/copilot-instructions.md
rm -rf .claude/skills/gps
```

### Remove Windsurf Integration

```bash
rm .windsurf/rules
rm -rf .claude/skills/gps
```

### Remove Configuration

```bash
rm .gps.yaml
```

---

## Getting Help

### Documentation

- **[SKILL.md](SKILL.md)** - Skill overview and quick start
- **[reference.md](reference.md)** - Complete command reference
- **[examples.md](examples.md)** - Real-world use cases
- **[formats.md](formats.md)** - Output format guide

### Community

- **GitHub Issues:** [github.com/wesbragagt/gps/issues](https://github.com/wesbragagt/gps/issues)
- **Discussions:** [github.com/wesbragagt/gps/discussions](https://github.com/wesbragagt/gps/discussions)

### Quick Commands Reference

| Task | Command |
|------|---------|
| Map project | `gps` |
| Quick overview | `gps --summary` |
| Find entry points | `gps --entry-points` |
| Limit depth | `gps -L 2` |
| JSON output | `gps -f json` |
| Token comparison | `gps --compare` |
| Help | `gps --help` |

---

*GPS Skill - Token-efficient project structure analysis for AI coding assistants.*
