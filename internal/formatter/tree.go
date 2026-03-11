// Package formatter provides output formatters for project structure.
package formatter

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wesbragagt/gps/pkg/types"
)

// TreeFormatter formats projects as a traditional tree view.
type TreeFormatter struct {
	// ShowMetadata displays inline [size, lines] for files.
	ShowMetadata bool

	// Colorize enables ANSI color output for different file types.
	Colorize bool

	// ShowSummary displays the summary line at the end.
	ShowSummary bool
}

// NewTreeFormatter creates a new tree formatter with default settings.
func NewTreeFormatter() *TreeFormatter {
	return &TreeFormatter{
		ShowMetadata: true,
		Colorize:     true,
		ShowSummary:  true,
	}
}

// Tree box drawing characters
const (
	branch    = "├── "
	last      = "└── "
	vertical  = "│   "
	empty     = "    "
	dirSuffix = "/"
)

// ANSI color codes
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	blue    = "\033[34m"
	cyan    = "\033[36m"
	yellow  = "\033[33m"
	green   = "\033[32m"
	magenta = "\033[35m"
	white   = "\033[37m"
	dim     = "\033[2m"
)

// colorEnabled checks if colors should be used.
func (f *TreeFormatter) colorEnabled() bool {
	if !f.Colorize {
		return false
	}
	// Respect NO_COLOR environment variable
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	return true
}

// colorize wraps text with ANSI color codes if enabled.
func (f *TreeFormatter) colorize(text, color string) string {
	if !f.colorEnabled() {
		return text
	}
	return color + text + reset
}

// getColorForFile returns the appropriate color for a file type.
func (f *TreeFormatter) getColorForFile(file *types.File) string {
	if !f.colorEnabled() {
		return ""
	}

	ext := strings.ToLower(filepath.Ext(file.Path))
	name := strings.ToLower(filepath.Base(file.Path))

	// Check by extension or filename
	switch {
	case file.IsBinary:
		return dim
	case ext == ".go":
		return cyan
	case ext == ".js", ext == ".jsx", ext == ".ts", ext == ".tsx", ext == ".mjs", ext == ".cjs":
		return yellow
	case ext == ".py", ext == ".pyw":
		return green
	case ext == ".rs":
		return bold + yellow
	case ext == ".java", ext == ".kt", ext == ".scala":
		return magenta
	case ext == ".rb":
		return red()
	case ext == ".c", ext == ".cpp", ext == ".cc", ext == ".cxx", ext == ".h", ext == ".hpp":
		return blue
	case name == "makefile", name == "dockerfile", name == "rakefile", name == "gemfile":
		return green
	case ext == ".json", ext == ".yaml", ext == ".yml", ext == ".toml", ext == ".ini", ext == ".cfg", ext == ".conf":
		return magenta
	case ext == ".md", ext == ".rst", ext == ".txt", ext == ".adoc":
		return white
	case ext == ".sh", ext == ".bash", ext == ".zsh", ext == ".fish":
		return green
	case ext == ".css", ext == ".scss", ext == ".sass", ext == ".less":
		return magenta
	case ext == ".html", ext == ".htm", ext == ".xml", ext == ".svg":
		return yellow
	case ext == ".sql":
		return cyan
	default:
		return ""
	}
}

func red() string {
	return "\033[31m"
}

// Format converts a Project to tree format string.
func (f *TreeFormatter) Format(project *types.Project) (string, error) {
	if project == nil {
		return "", nil
	}

	var sb strings.Builder

	// Root directory name
	rootName := project.Name
	if rootName == "" {
		rootName = "."
	}

	sb.WriteString(f.colorize(rootName+dirSuffix, bold+blue))
	sb.WriteString("\n")

	// Format tree
	f.formatTree(&project.Tree, &sb, "", true)

	// Summary line
	if f.ShowSummary && project.Stats.FileCount > 0 {
		sb.WriteString("\n")
		sb.WriteString(f.formatSummary(&project.Stats))
	}

	return sb.String(), nil
}

// formatTree recursively formats a directory tree.
func (f *TreeFormatter) formatTree(dir *types.Directory, sb *strings.Builder, prefix string, isRoot bool) {
	if dir == nil {
		return
	}

	// Collect all items (directories and files) for sorting
	type item struct {
		name   string
		isDir  bool
		file   *types.File
		subdir *types.Directory
	}

	var items []item

	// Add subdirectories
	for i := range dir.Subdirs {
		name := filepath.Base(dir.Subdirs[i].Path)
		if name == "" || name == "." {
			name = "unnamed"
		}
		items = append(items, item{
			name:   name,
			isDir:  true,
			subdir: &dir.Subdirs[i],
		})
	}

	// Add files
	for i := range dir.Files {
		name := filepath.Base(dir.Files[i].Path)
		items = append(items, item{
			name:  name,
			isDir: false,
			file:  &dir.Files[i],
		})
	}

	// Sort items: directories first, then files, alphabetically within each group
	sort.Slice(items, func(i, j int) bool {
		if items[i].isDir != items[j].isDir {
			return items[i].isDir
		}
		return strings.ToLower(items[i].name) < strings.ToLower(items[j].name)
	})

	// Render each item
	for i, it := range items {
		isLast := i == len(items)-1

		// Determine the connector
		connector := branch
		if isLast {
			connector = last
		}

		if it.isDir {
			// Directory
			dirName := f.colorize(it.name+dirSuffix, bold+blue)
			sb.WriteString(prefix)
			sb.WriteString(connector)
			sb.WriteString(dirName)
			sb.WriteString("\n")

			// Calculate next prefix for children
			nextPrefix := prefix
			if isLast {
				nextPrefix += empty
			} else {
				nextPrefix += vertical
			}

			// Recurse into subdirectory
			f.formatTree(it.subdir, sb, nextPrefix, false)
		} else {
			// File
			fileName := it.name
			color := f.getColorForFile(it.file)
			if color != "" {
				fileName = f.colorize(fileName, color)
			}

			sb.WriteString(prefix)
			sb.WriteString(connector)
			sb.WriteString(fileName)

			// Metadata
			if f.ShowMetadata && it.file != nil {
				meta := f.formatFileMetadata(it.file)
				if meta != "" {
					sb.WriteString(" ")
					sb.WriteString(f.colorize(meta, dim))
				}
			}

			sb.WriteString("\n")
		}
	}
}

// formatFileMetadata formats inline metadata for a file.
func (f *TreeFormatter) formatFileMetadata(file *types.File) string {
	if file == nil {
		return ""
	}

	var parts []string

	// Size
	if file.Size > 0 {
		parts = append(parts, formatSize(file.Size))
	}

	// Lines (skip if 0 or binary)
	if file.Lines > 0 && !file.IsBinary {
		parts = append(parts, fmt.Sprintf("%dL", file.Lines))
	}

	if len(parts) == 0 {
		return ""
	}

	return "[" + strings.Join(parts, ", ") + "]"
}

// formatSummary formats the summary line.
func (f *TreeFormatter) formatSummary(stats *types.Stats) string {
	if stats == nil {
		return ""
	}

	var parts []string

	// File count
	if stats.FileCount == 1 {
		parts = append(parts, "1 file")
	} else {
		parts = append(parts, fmt.Sprintf("%d files", stats.FileCount))
	}

	// Total size
	if stats.TotalSize > 0 {
		parts = append(parts, formatSize(stats.TotalSize))
	}

	// Total lines
	if stats.TotalLines > 0 {
		parts = append(parts, fmt.Sprintf("%d lines", stats.TotalLines))
	}

	return strings.Join(parts, ", ")
}
