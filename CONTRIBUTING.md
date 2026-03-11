# Contributing to gps

Thank you for your interest in contributing to gps! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

This project follows the [Go Code of Conduct](https://go.dev/conduct). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- A code editor (VS Code, GoLand, etc.)

### Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/gps.git
   cd gps
   ```

3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/wesbragagt/gps.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Build the project:
   ```bash
   go build ./cmd/gps
   ```

6. Run tests:
   ```bash
   go test ./...
   ```

## Project Structure

```
gps/
├── cmd/
│   ├── gps/           # Main application entry point
│   │   └── main.go
│   ├── root.go        # Root command definition
│   └── version.go     # Version command
├── internal/          # Private application code
│   ├── detector/      # Project type detection
│   ├── formatter/     # Output formatters (TOON, JSON, tree, flat)
│   ├── metadata/      # File metadata extraction
│   ├── scanner/       # Filesystem scanning
│   └── version/       # Version information
├── pkg/
│   └── types/         # Public types (Project, File, Directory, etc.)
├── docs/              # Documentation
├── go.mod
├── go.sum
├── README.md
└── CONTRIBUTING.md
```

### Key Packages

- **`pkg/types`** - Core domain types. Changes here affect the entire project.
- **`internal/scanner`** - Filesystem traversal logic.
- **`internal/formatter`** - Output formatting. Add new formats here.
- **`internal/detector`** - Project type and key file detection.

## Making Changes

### Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `go fmt` before committing
- Use meaningful variable names
- Add godoc comments for exported functions and types
- Keep functions focused and small

### Adding a New Output Format

1. Create a new file in `internal/formatter/`:
   ```go
   package formatter

   // MyFormatter formats projects in custom format.
   type MyFormatter struct {
       // Options...
   }

   // NewMyFormatter creates a new formatter.
   func NewMyFormatter() *MyFormatter {
       return &MyFormatter{}
   }

   // Format converts a Project to the custom format.
   func (f *MyFormatter) Format(project *types.Project) (string, error) {
       // Implementation
   }
   ```

2. Register in `cmd/root.go`:
   ```go
   func getFormatter(format string) Formatter {
       switch format {
       case "myformat":
           return formatter.NewMyFormatter()
       // ...
       }
   }
   ```

3. Add tests in `internal/formatter/myformat_test.go`

### Adding Project Detection

1. Add patterns to `internal/detector/detector.go`
2. Update `projectTypeRules` for new languages
3. Add entry point patterns if applicable
4. Add tests

## Testing

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Run Specific Package Tests

```bash
go test ./internal/formatter/...
```

### Run Specific Test

```bash
go test -run TestToonFormatter ./internal/formatter/
```

### Writing Tests

- Place tests in the same package with `_test.go` suffix
- Use table-driven tests for multiple cases
- Test edge cases (nil inputs, empty projects, etc.)
- Use testdata directories for fixture files

Example:
```go
func TestMyFormatter(t *testing.T) {
    tests := []struct {
        name    string
        project *types.Project
        want    string
        wantErr bool
    }{
        {
            name: "empty project",
            project: &types.Project{Name: "empty"},
            want: "project[empty]{}",
        },
        // More cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f := NewMyFormatter()
            got, err := f.Format(tt.project)
            if (err != nil) != tt.wantErr {
                t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("Format() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Commit Guidelines

### Commit Message Format

```
<type>: <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Adding/updating tests
- `chore`: Maintenance tasks

### Examples

```
feat: add YAML output format

Add support for YAML output with -f yaml flag.
Includes full project structure and metadata.

Closes #42
```

```
fix: handle permission denied errors gracefully

Scanner now logs permission errors and continues
instead of failing the entire scan.
```

## Pull Request Process

1. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make your changes** following the guidelines above

3. **Run tests**:
   ```bash
   go test ./...
   go fmt ./...
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/my-feature
   ```

6. **Create a Pull Request** on GitHub

### PR Checklist

- [ ] Code compiles without warnings
- [ ] Tests pass locally
- [ ] New code has tests
- [ ] Documentation updated (if needed)
- [ ] Commit messages follow guidelines
- [ ] PR description explains the change

### Review Process

1. Maintainers will review your PR
2. Address any feedback
3. Once approved, a maintainer will merge

## Questions?

- Open an issue for bugs or feature requests
- Start a discussion for questions or ideas

Thank you for contributing!
