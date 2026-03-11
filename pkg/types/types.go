// Package types defines core domain types for the gps project mapping system.
package types

import (
	"time"
)

// ProjectType represents the type of project being mapped.
type ProjectType string

const (
	ProjectTypeGo     ProjectType = "go"
	ProjectTypeNode   ProjectType = "node"
	ProjectTypePython ProjectType = "python"
	ProjectTypeRust   ProjectType = "rust"
	ProjectTypeJava   ProjectType = "java"
	ProjectTypeMixed  ProjectType = "mixed"
	ProjectTypeOther  ProjectType = "other"
)

// OutputFormat represents the output format for project mapping.
type OutputFormat string

const (
	FormatTOON OutputFormat = "toon"
	FormatJSON OutputFormat = "json"
	FormatTree OutputFormat = "tree"
	FormatFlat OutputFormat = "flat"
)

// File represents a single file in the project.
type File struct {
	// Path is the relative path from the project root.
	Path string `json:"path"`

	// Size is the file size in bytes.
	Size int64 `json:"size"`

	// Lines is the line count for text files (0 for binary files).
	Lines int `json:"lines"`

	// Type is the language or file type (e.g., "go", "js", "md", "binary").
	Type string `json:"type"`

	// ModTime is the file modification timestamp.
	ModTime time.Time `json:"mod_time"`

	// IsBinary indicates whether the file is binary (true) or text (false).
	IsBinary bool `json:"is_binary"`

	// Importance is a score (0-100) indicating the file's relevance to the project.
	// Higher scores indicate more important files.
	Importance int `json:"importance"`

	// IsGenerated indicates whether the file is auto-generated.
	IsGenerated bool `json:"is_generated"`
}

// Directory represents a directory in the project structure.
type Directory struct {
	// Path is the relative path from the project root.
	Path string `json:"path"`

	// Files contains the files directly in this directory.
	Files []File `json:"files"`

	// Subdirs contains nested subdirectories.
	Subdirs []Directory `json:"subdirs"`

	// IsExpanded indicates whether this directory should be expanded in tree rendering.
	IsExpanded bool `json:"is_expanded"`
}

// Stats contains aggregate statistics about the project.
type Stats struct {
	// FileCount is the total number of files in the project.
	FileCount int `json:"file_count"`

	// TotalSize is the total size of all files in bytes.
	TotalSize int64 `json:"total_size"`

	// TotalLines is the total lines of code across all text files.
	TotalLines int `json:"total_lines"`

	// ByType contains the count of files grouped by their type.
	ByType map[string]int `json:"by_type"`
}

// KeyFiles identifies important files in the project.
type KeyFiles struct {
	// EntryPoints contains paths to main/entry point files.
	EntryPoints []string `json:"entry_points"`

	// Configs contains paths to configuration files.
	Configs []string `json:"configs"`

	// Tests contains paths to test files.
	Tests []string `json:"tests"`

	// Docs contains paths to documentation files.
	Docs []string `json:"docs"`
}

// Project represents the complete project structure and metadata.
type Project struct {
	// Name is the project name (typically the root directory name).
	Name string `json:"name"`

	// Type is the detected project type (go, node, python, etc.).
	Type ProjectType `json:"type"`

	// Root is the absolute path to the project root directory.
	Root string `json:"root"`

	// Stats contains aggregate statistics about the project.
	Stats Stats `json:"stats"`

	// KeyFiles identifies important files in the project.
	KeyFiles KeyFiles `json:"key_files"`

	// Tree is the root directory structure containing all files and subdirectories.
	Tree Directory `json:"tree"`
}

// Config holds configuration options for project scanning.
type Config struct {
	// Format specifies the output format (toon, json, tree, flat).
	Format OutputFormat `json:"format"`

	// Depth is the maximum traversal depth (-1 for unlimited).
	Depth int `json:"depth"`

	// IncludeHidden determines whether to include hidden files/directories.
	IncludeHidden bool `json:"include_hidden"`

	// ExcludePatterns contains glob patterns for files/directories to exclude.
	ExcludePatterns []string `json:"exclude_patterns"`

	// IncludePatterns contains glob patterns for files/directories to include.
	// If specified, only matching files are included.
	IncludePatterns []string `json:"include_patterns"`

	// ShowMetadata determines whether to include file metadata (size, lines, etc.).
	ShowMetadata bool `json:"show_metadata"`

	// ShowProjectInfo determines whether to include project-level information.
	ShowProjectInfo bool `json:"show_project_info"`

	// SmartTraverse enables intelligent traversal that prioritizes important files.
	SmartTraverse bool `json:"smart_traverse"`

	// FocusPath specifies a subpath to focus scanning on (empty for full project).
	FocusPath string `json:"focus_path"`
}

// NewConfig creates a new Config with sensible defaults.
func NewConfig() *Config {
	return &Config{
		Format:          FormatTOON,
		Depth:           -1,
		IncludeHidden:   false,
		ExcludePatterns: []string{},
		IncludePatterns: []string{},
		ShowMetadata:    true,
		ShowProjectInfo: true,
		SmartTraverse:   true,
		FocusPath:       "",
	}
}
