// Package scanner provides filesystem traversal and project scanning capabilities.
package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/wesbragagt/gps/pkg/types"
)

// Config holds scanner configuration options.
type Config struct {
	// MaxDepth is the maximum directory depth to traverse (-1 for unlimited).
	MaxDepth int

	// IncludeHidden determines whether to include hidden files/directories.
	IncludeHidden bool

	// CurrentlyOpen contains paths that are currently open (for future filtering).
	CurrentlyOpen []string

	// ExcludePatterns contains glob patterns for files/directories to exclude.
	ExcludePatterns []string

	// IncludePatterns contains glob patterns for files/directories to include.
	// If specified, only matching files are included.
	IncludePatterns []string

	// ExcludeExtensions contains file extensions to exclude (e.g., ".log", ".tmp").
	ExcludeExtensions []string

	// IncludeExtensions contains file extensions to include.
	// If specified, only files with these extensions are included.
	IncludeExtensions []string

	// RespectGitignore determines whether to respect .gitignore patterns.
	RespectGitignore bool
}

// NewConfig creates a new Config with sensible defaults.
func NewConfig() *Config {
	return &Config{
		MaxDepth:          -1,
		IncludeHidden:     false,
		CurrentlyOpen:     nil,
		ExcludePatterns:   nil,
		IncludePatterns:   nil,
		ExcludeExtensions: nil,
		IncludeExtensions: nil,
		RespectGitignore:  true,
	}
}

// Scanner performs filesystem traversal and builds directory trees.
type Scanner struct {
	config    *Config
	errors    []error
	filter    *Filter
	gitignore *GitignoreMatcher
	root      string
}

// New creates a new Scanner with the given configuration.
func New(config *Config) *Scanner {
	if config == nil {
		config = NewConfig()
	}
	return &Scanner{
		config: config,
		errors: make([]error, 0),
	}
}

// Errors returns any non-fatal errors encountered during scanning.
func (s *Scanner) Errors() []error {
	return s.errors
}

// Walk traverses the filesystem starting at root and builds a Directory tree.
// The returned Directory represents the root of the tree structure.
func (s *Scanner) Walk(root string) (*types.Directory, error) {
	// Clean and validate root path
	root = filepath.Clean(root)

	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("cannot access root path %q: %w", root, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("root path %q is not a directory", root)
	}

	// Reset errors from previous walks
	s.errors = make([]error, 0)
	s.root = root

	// Initialize gitignore matcher if enabled
	if s.config.RespectGitignore {
		s.gitignore = NewGitignoreMatcher(root)
		if err := s.gitignore.Load(); err != nil {
			// Non-fatal: log and continue without gitignore
			s.errors = append(s.errors, &ScanError{
				Path:  root,
				Op:    "gitignore",
				Cause: err,
			})
			s.gitignore = nil
		}
	}

	// Initialize filter
	s.filter = NewFilter(s.config, s.gitignore)

	// Build the tree recursively
	dir, err := s.walkDir(root, ".", 0)
	if err != nil {
		return nil, err
	}

	return dir, nil
}

// walkDir recursively walks a directory and builds the tree structure.
// basePath is the absolute path, relPath is the relative path from root.
// depth is the current traversal depth.
func (s *Scanner) walkDir(basePath, relPath string, depth int) (*types.Directory, error) {
	dir := &types.Directory{
		Path:       relPath,
		Files:      make([]types.File, 0),
		Subdirs:    make([]types.Directory, 0),
		IsExpanded: false,
	}

	// Read directory entries
	entries, err := os.ReadDir(basePath)
	if err != nil {
		// Handle permission errors - log and return empty directory
		if os.IsPermission(err) {
			s.errors = append(s.errors, &ScanError{
				Path:  basePath,
				Op:    "readdir",
				Cause: err,
			})
			return dir, nil
		}
		return nil, fmt.Errorf("failed to read directory %q: %w", basePath, err)
	}

	// Load gitignore for this directory if it exists
	if s.gitignore != nil && relPath != "." {
		if s.gitignore.HasSubdirGitignore(relPath) {
			if err := s.gitignore.LoadSubdirGitignore(relPath); err != nil {
				// Non-fatal: log and continue
				s.errors = append(s.errors, &ScanError{
					Path:  basePath,
					Op:    "gitignore",
					Cause: err,
				})
			}
		}
	}

	// Separate files and directories for consistent ordering
	var fileEntries []fs.DirEntry
	var dirEntries []fs.DirEntry

	for _, entry := range entries {
		entryRelPath := filepath.Join(relPath, entry.Name())

		// Apply filter to skip excluded entries early
		if s.filter != nil && s.filter.ShouldSkip(entryRelPath, entry) {
			continue
		}

		// For files, check if they should be included
		if !entry.IsDir() && s.filter != nil && !s.filter.ShouldInclude(entryRelPath, entry) {
			continue
		}

		if entry.IsDir() {
			dirEntries = append(dirEntries, entry)
		} else {
			fileEntries = append(fileEntries, entry)
		}
	}

	// Process files first
	for _, entry := range fileEntries {
		file, err := s.createFile(basePath, relPath, entry)
		if err != nil {
			// Non-fatal: log and continue
			s.errors = append(s.errors, err)
			continue
		}
		dir.Files = append(dir.Files, *file)
	}

	// Process subdirectories if within depth limit
	if s.config.MaxDepth < 0 || depth < s.config.MaxDepth {
		for _, entry := range dirEntries {
			entryRelPath := filepath.Join(relPath, entry.Name())
			entryAbsPath := filepath.Join(basePath, entry.Name())

			subdir, err := s.walkDir(entryAbsPath, entryRelPath, depth+1)
			if err != nil {
				// Non-fatal: log and continue
				s.errors = append(s.errors, err)
				continue
			}
			dir.Subdirs = append(dir.Subdirs, *subdir)
		}
	}

	return dir, nil
}

// createFile creates a File struct from a directory entry.
func (s *Scanner) createFile(basePath, relPath string, entry fs.DirEntry) (*types.File, error) {
	// Get file info (type is available from DirEntry, but we need size/modtime)
	info, err := entry.Info()
	if err != nil {
		if os.IsPermission(err) {
			return nil, &ScanError{
				Path:  filepath.Join(basePath, entry.Name()),
				Op:    "stat",
				Cause: err,
			}
		}
		return nil, fmt.Errorf("failed to get file info for %q: %w", entry.Name(), err)
	}

	file := &types.File{
		Path:       filepath.Join(relPath, entry.Name()),
		Size:       info.Size(),
		Lines:      0,  // Set by metadata extractor
		Type:       "", // Set by metadata extractor
		ModTime:    info.ModTime(),
		IsBinary:   false, // Set by metadata extractor
		Importance: 0,     // Set by importance scorer
	}

	return file, nil
}

// isHidden checks if a file/directory name is hidden (starts with dot).
func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

// isDirectory checks if a DirEntry is a directory.
func isDirectory(entry fs.DirEntry) bool {
	return entry.IsDir()
}

// ScanError represents a non-fatal error during scanning.
type ScanError struct {
	Path  string
	Op    string
	Cause error
}

// Error implements the error interface.
func (e *ScanError) Error() string {
	return fmt.Sprintf("scanner: %s failed for %q: %v", e.Op, e.Path, e.Cause)
}

// Unwrap returns the underlying cause.
func (e *ScanError) Unwrap() error {
	return e.Cause
}
