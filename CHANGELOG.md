# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Extension Filtering** (NEW!)
  - `--ext` / `-e` flag to include only specific file extensions
  - `--exclude-ext` flag to exclude specific file extensions
  - Supports comma-separated list (e.g., `-e go,md,yaml`)
  - Leading dot optional (`.go` and `go` both work)
  - Case-insensitive matching
  - Config file support via `ext` and `exclude-ext` keys

## [0.1.0] - 2024-01-15

### Added

- **Core Features**
  - Filesystem scanning with recursive directory traversal
  - gitignore support for filtering files
  - Hidden file filtering (toggle with `-a`)
  - Depth limiting (toggle with `-L`)
  - Pattern-based include/exclude filtering

- **Output Formats**
  - TOON (Token-Optimized Object Notation) - default format
  - JSON with pretty-print support
  - Tree view with ANSI color coding
  - Flat CSV-like format

- **Project Intelligence**
  - Automatic project type detection (Go, Node.js, Python, Rust, Java)
  - Entry point detection (main.go, index.js, etc.)
  - Configuration file identification
  - Test file detection
  - Documentation file detection

- **Metadata**
  - File size information
  - Line count for text files
  - Language/file type detection
  - Binary file detection
  - Importance scoring

- **CLI Features**
  - Tree-compatible flags (`-L`, `-d`, `-a`, `-I`, `-P`)
  - Configuration file support (`.gps.yaml`)
  - Focus mode for subdirectory analysis
  - Summary mode for quick overview
  - Entry points mode for finding main files
  - Version command

- **Documentation**
  - README with usage examples
  - Output format documentation
  - Contributing guidelines

### Changed

- N/A (initial release)

### Deprecated

- N/A (initial release)

### Removed

- N/A (initial release)

### Fixed

- N/A (initial release)

---

## Version History

| Version | Date       | Highlights                          |
|---------|------------|-------------------------------------|
| 0.1.0   | 2024-01-15 | Initial release with core features  |

---

[Unreleased]: https://github.com/wesbragagt/gps/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/wesbragagt/gps/releases/tag/v0.1.0
