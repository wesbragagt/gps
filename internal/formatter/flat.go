// Package formatter provides output formatters for project structure.
package formatter

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/wesbragagt/gps/pkg/types"
)

// FlatFormatter formats projects as flat CSV-like output.
type FlatFormatter struct {
	// ShowHeader determines whether to show column headers.
	ShowHeader bool

	// SortBy specifies which field to sort by: "path", "size", "lines", "type", "importance", "modtime".
	SortBy string

	// SortDescending reverses the sort order when true.
	SortDescending bool

	// ShowFields specifies which fields to include (path is always included as first column).
	// Available: path, size, lines, type, importance, modtime
	ShowFields []string
}

// NewFlatFormatter creates a new flat formatter with default settings.
func NewFlatFormatter() *FlatFormatter {
	return &FlatFormatter{
		ShowHeader:     true,
		SortBy:         "path",
		SortDescending: false,
		ShowFields:     []string{"path", "size", "lines", "type", "importance"},
	}
}

// flatFile represents a flattened file for sorting and formatting.
type flatFile struct {
	Path       string
	Size       int64
	Lines      int
	Type       string
	Importance int
	ModTime    time.Time
}

// Format converts a Project to flat CSV-like format string.
func (f *FlatFormatter) Format(project *types.Project) (string, error) {
	if project == nil {
		return "", nil
	}

	// Flatten tree to list
	files := f.flattenTree(&project.Tree)

	// Sort files
	f.sortFiles(files)

	// Build output
	var sb strings.Builder

	// Header row
	if f.ShowHeader {
		sb.WriteString(f.formatHeader())
		sb.WriteString("\n")
	}

	// Data rows
	for _, file := range files {
		sb.WriteString(f.formatRow(file))
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// flattenTree recursively flattens directory tree to slice of files.
func (f *FlatFormatter) flattenTree(dir *types.Directory) []flatFile {
	if dir == nil {
		return nil
	}

	var files []flatFile

	// Add files in this directory
	for i := range dir.Files {
		files = append(files, flatFile{
			Path:       dir.Files[i].Path,
			Size:       dir.Files[i].Size,
			Lines:      dir.Files[i].Lines,
			Type:       dir.Files[i].Type,
			Importance: dir.Files[i].Importance,
			ModTime:    dir.Files[i].ModTime,
		})
	}

	// Recurse into subdirectories
	for i := range dir.Subdirs {
		files = append(files, f.flattenTree(&dir.Subdirs[i])...)
	}

	return files
}

// sortFiles sorts the flat file list based on SortBy and SortDescending.
func (f *FlatFormatter) sortFiles(files []flatFile) {
	sort.Slice(files, func(i, j int) bool {
		var less bool

		switch f.SortBy {
		case "size":
			less = files[i].Size < files[j].Size
		case "lines":
			less = files[i].Lines < files[j].Lines
		case "type":
			less = strings.ToLower(files[i].Type) < strings.ToLower(files[j].Type)
		case "importance":
			less = files[i].Importance < files[j].Importance
		case "modtime":
			less = files[i].ModTime.Before(files[j].ModTime)
		default: // "path"
			less = strings.ToLower(files[i].Path) < strings.ToLower(files[j].Path)
		}

		if f.SortDescending {
			return !less
		}
		return less
	})
}

// formatHeader returns the CSV header row.
func (f *FlatFormatter) formatHeader() string {
	var cols []string
	for _, field := range f.ShowFields {
		cols = append(cols, field)
	}
	return strings.Join(cols, ",")
}

// formatRow returns a single CSV row for a file.
func (f *FlatFormatter) formatRow(file flatFile) string {
	var cols []string

	for _, field := range f.ShowFields {
		switch field {
		case "path":
			cols = append(cols, csvEscape(file.Path))
		case "size":
			cols = append(cols, fmt.Sprintf("%d", file.Size))
		case "lines":
			cols = append(cols, fmt.Sprintf("%d", file.Lines))
		case "type":
			cols = append(cols, csvEscape(file.Type))
		case "importance":
			cols = append(cols, fmt.Sprintf("%d", file.Importance))
		case "modtime":
			if file.ModTime.IsZero() {
				cols = append(cols, "")
			} else {
				cols = append(cols, file.ModTime.Format(time.RFC3339))
			}
		default:
			cols = append(cols, "")
		}
	}

	return strings.Join(cols, ",")
}

// csvEscape quotes a field if it contains commas, quotes, or newlines.
func csvEscape(s string) string {
	if strings.ContainsAny(s, ",\"\n\r") {
		// Escape quotes by doubling them
		escaped := strings.ReplaceAll(s, "\"", "\"\"")
		return "\"" + escaped + "\""
	}
	return s
}
