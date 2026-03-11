// Package formatter provides output formatters for project structure.
//
// This package implements multiple output formats:
//   - TOON: Token-Optimized Object Notation (default, most efficient for AI)
//   - JSON: Standard JSON output for programmatic access
//   - Tree: Traditional tree view with colors and inline metadata
//   - Flat: CSV-like format for data analysis
//
// All formatters implement the Formatter interface:
//
//	type Formatter interface {
//	    Format(project *types.Project) (string, error)
//	}
//
// Example usage:
//
//	f := formatter.NewToonFormatter()
//	output, err := f.Format(project)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(output)
package formatter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wesbragagt/gps/pkg/types"
)

// ToonFormatter formats projects in Token-Optimized Object Notation.
type ToonFormatter struct{}

// NewToonFormatter creates a new TOON formatter.
func NewToonFormatter() *ToonFormatter {
	return &ToonFormatter{}
}

// Format converts a Project to TOON format string.
func (f *ToonFormatter) Format(project *types.Project) (string, error) {
	if project == nil {
		return "", nil
	}

	var sb strings.Builder

	// Project header
	sb.WriteString(f.formatProjectHeader(project))
	sb.WriteString("\n")

	// Tree section
	if project.Tree.Path != "" || len(project.Tree.Files) > 0 || len(project.Tree.Subdirs) > 0 {
		sb.WriteString(f.formatTree(&project.Tree, 0))
		sb.WriteString("\n")
	}

	// Key files section
	if hasKeyFiles(&project.KeyFiles) {
		sb.WriteString(f.formatKeyFiles(&project.KeyFiles))
	}

	return sb.String(), nil
}

// formatProjectHeader formats the project metadata section.
func (f *ToonFormatter) formatProjectHeader(p *types.Project) string {
	var sb strings.Builder

	sb.WriteString("project[")
	sb.WriteString(p.Name)
	sb.WriteString("]{\n")

	// Type
	sb.WriteString("  type: ")
	sb.WriteString(string(p.Type))
	sb.WriteString("\n")

	// Files count
	if p.Stats.FileCount > 0 {
		sb.WriteString("  files: ")
		sb.WriteString(fmt.Sprintf("%d", p.Stats.FileCount))
		sb.WriteString("\n")
	}

	// Total size
	if p.Stats.TotalSize > 0 {
		sb.WriteString("  size: ")
		sb.WriteString(formatSize(p.Stats.TotalSize))
		sb.WriteString("\n")
	}

	// Total lines (if meaningful)
	if p.Stats.TotalLines > 0 {
		sb.WriteString("  lines: ")
		sb.WriteString(fmt.Sprintf("%d", p.Stats.TotalLines))
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}

// formatTree formats a directory tree recursively.
func (f *ToonFormatter) formatTree(dir *types.Directory, indent int) string {
	if dir == nil {
		return ""
	}

	var sb strings.Builder
	prefix := strings.Repeat("  ", indent)

	// Directory header
	dirName := filepath.Base(dir.Path)
	if dirName == "" || dirName == "." {
		dirName = "root"
	}

	totalItems := len(dir.Files) + len(dir.Subdirs)
	sb.WriteString(prefix)
	sb.WriteString(dirName)
	if totalItems > 0 {
		sb.WriteString(fmt.Sprintf("[%d]", totalItems))
	}
	sb.WriteString("{\n")

	// Files first (sorted by name is handled by caller if needed)
	for i := range dir.Files {
		sb.WriteString(prefix)
		sb.WriteString("  ")
		sb.WriteString(f.formatFile(&dir.Files[i]))
		sb.WriteString("\n")
	}

	// Subdirectories
	for i := range dir.Subdirs {
		sb.WriteString(f.formatTree(&dir.Subdirs[i], indent+1))
	}

	sb.WriteString(prefix)
	sb.WriteString("}\n")

	return sb.String()
}

// formatFile formats a single file with inline metadata.
func (f *ToonFormatter) formatFile(file *types.File) string {
	if file == nil {
		return ""
	}

	var sb strings.Builder

	// Filename
	sb.WriteString(filepath.Base(file.Path))

	// Metadata in brackets
	var meta []string

	// Size
	if file.Size > 0 {
		meta = append(meta, formatSize(file.Size))
	}

	// Lines (skip if 0 or binary)
	if file.Lines > 0 && !file.IsBinary {
		meta = append(meta, fmt.Sprintf("%dL", file.Lines))
	}

	// Type
	if file.Type != "" && file.Type != "unknown" {
		meta = append(meta, file.Type)
	}

	// Binary indicator
	if file.IsBinary {
		meta = append(meta, "bin")
	}

	if len(meta) > 0 {
		sb.WriteString(" [")
		sb.WriteString(strings.Join(meta, ", "))
		sb.WriteString("]")
	}

	return sb.String()
}

// formatKeyFiles formats the key files section.
func (f *ToonFormatter) formatKeyFiles(kf *types.KeyFiles) string {
	if kf == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("keyfiles{\n")

	// Entry points
	if len(kf.EntryPoints) > 0 {
		sb.WriteString("  entry: ")
		if len(kf.EntryPoints) == 1 {
			sb.WriteString(kf.EntryPoints[0])
		} else {
			sb.WriteString(fmt.Sprintf("%d files", len(kf.EntryPoints)))
		}
		sb.WriteString("\n")
	}

	// Configs
	if len(kf.Configs) > 0 {
		sb.WriteString("  config: ")
		if len(kf.Configs) == 1 {
			sb.WriteString(kf.Configs[0])
		} else {
			sb.WriteString(fmt.Sprintf("%d files", len(kf.Configs)))
		}
		sb.WriteString("\n")
	}

	// Tests
	if len(kf.Tests) > 0 {
		sb.WriteString("  tests: ")
		if len(kf.Tests) == 1 {
			sb.WriteString(kf.Tests[0])
		} else {
			sb.WriteString(fmt.Sprintf("%d files", len(kf.Tests)))
		}
		sb.WriteString("\n")
	}

	// Docs
	if len(kf.Docs) > 0 {
		sb.WriteString("  docs: ")
		if len(kf.Docs) == 1 {
			sb.WriteString(kf.Docs[0])
		} else {
			sb.WriteString(fmt.Sprintf("%d files", len(kf.Docs)))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}

// formatSize converts bytes to human-readable format.
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		v := float64(bytes) / GB
		if v == float64(int64(v)) {
			return fmt.Sprintf("%dGB", int64(v))
		}
		return fmt.Sprintf("%.1fGB", v)
	case bytes >= MB:
		v := float64(bytes) / MB
		if v == float64(int64(v)) {
			return fmt.Sprintf("%dMB", int64(v))
		}
		return fmt.Sprintf("%.1fMB", v)
	case bytes >= KB:
		v := float64(bytes) / KB
		if v == float64(int64(v)) {
			return fmt.Sprintf("%dKB", int64(v))
		}
		return fmt.Sprintf("%.1fKB", v)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

// hasKeyFiles checks if KeyFiles has any content.
func hasKeyFiles(kf *types.KeyFiles) bool {
	if kf == nil {
		return false
	}
	return len(kf.EntryPoints) > 0 ||
		len(kf.Configs) > 0 ||
		len(kf.Tests) > 0 ||
		len(kf.Docs) > 0
}
