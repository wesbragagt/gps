package scanner

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// Filter determines if files/directories should be included or excluded.
type Filter struct {
	config          *Config
	gitignore       *GitignoreMatcher
	defaultExcludes []string
}

// NewFilter creates a new Filter with the given configuration.
func NewFilter(config *Config, gitignore *GitignoreMatcher) *Filter {
	return &Filter{
		config:    config,
		gitignore: gitignore,
		defaultExcludes: []string{
			".git", ".svn", ".hg",
			"node_modules", "vendor",
			"__pycache__", ".pytest_cache",
			"target", "build", "dist",
			".idea", ".vscode",
			".DS_Store", "Thumbs.db",
		},
	}
}

// ShouldSkip determines if a directory entry should be skipped entirely.
// Returns true if the entry should be excluded from traversal.
func (f *Filter) ShouldSkip(relPath string, entry fs.DirEntry) bool {
	name := entry.Name()
	isDir := entry.IsDir()

	// Check hidden files
	if !f.config.IncludeHidden && isHidden(name) {
		return true
	}

	// Check default excluded directories
	if isDir && f.isDefaultExcluded(name) {
		return true
	}

	// Check gitignore
	if f.config.RespectGitignore && f.gitignore != nil {
		if f.gitignore.Matches(relPath, isDir) {
			return true
		}
	}

	// Check exclude patterns
	if f.matchesExcludePattern(relPath, isDir) {
		// But allow include patterns to override
		if len(f.config.IncludePatterns) > 0 {
			if f.matchesIncludePattern(relPath, isDir) {
				return false
			}
		}
		return true
	}

	return false
}

// ShouldInclude determines if a file should be included in the result.
// This is a secondary check after ShouldSkip returns false.
func (f *Filter) ShouldInclude(relPath string, entry fs.DirEntry) bool {
	// Directories are always included if they pass ShouldSkip
	if entry.IsDir() {
		return true
	}

	// If include patterns specified, must match at least one
	if len(f.config.IncludePatterns) > 0 {
		if !f.matchesIncludePattern(relPath, false) {
			return false
		}
	}

	// If include extensions specified, must match at least one
	if len(f.config.IncludeExtensions) > 0 {
		if !f.matchesIncludeExtension(relPath) {
			return false
		}
	}

	// Check exclude extensions (but include patterns override)
	if len(f.config.IncludeExtensions) == 0 && len(f.config.IncludePatterns) == 0 {
		if f.matchesExcludeExtension(relPath) {
			return false
		}
	}

	return true
}

// matchesPattern checks if a path matches a glob pattern.
// Supports ** for recursive matching via doublestar library.
func matchesPattern(path, pattern string) bool {
	// Normalize paths for cross-platform compatibility
	path = filepath.ToSlash(path)
	pattern = filepath.ToSlash(pattern)

	// Handle simple patterns with filepath.Match first (faster)
	if !strings.Contains(pattern, "**") {
		matched, _ := filepath.Match(pattern, path)
		if matched {
			return true
		}
		// Also try matching just the filename
		matched, _ = filepath.Match(pattern, filepath.Base(path))
		return matched
	}

	// Use doublestar for ** patterns
	matched, _ := doublestar.Match(pattern, path)
	return matched
}

// matchesExcludePattern checks if path matches any exclude pattern.
func (f *Filter) matchesExcludePattern(relPath string, isDir bool) bool {
	for _, pattern := range f.config.ExcludePatterns {
		if matchesPattern(relPath, pattern) {
			return true
		}
		// For directories, also check with trailing slash
		if isDir {
			if matchesPattern(relPath+"/", pattern) {
				return true
			}
		}
	}
	return false
}

// matchesIncludePattern checks if path matches any include pattern.
func (f *Filter) matchesIncludePattern(relPath string, isDir bool) bool {
	for _, pattern := range f.config.IncludePatterns {
		if matchesPattern(relPath, pattern) {
			return true
		}
		// For directories, also check with trailing slash
		if isDir {
			if matchesPattern(relPath+"/", pattern) {
				return true
			}
		}
	}
	return false
}

// matchesExcludeExtension checks if file has an excluded extension.
func (f *Filter) matchesExcludeExtension(relPath string) bool {
	ext := strings.ToLower(filepath.Ext(relPath))
	for _, excludeExt := range f.config.ExcludeExtensions {
		if strings.ToLower(excludeExt) == ext {
			return true
		}
	}
	return false
}

// matchesIncludeExtension checks if file has an included extension.
func (f *Filter) matchesIncludeExtension(relPath string) bool {
	ext := strings.ToLower(filepath.Ext(relPath))
	for _, includeExt := range f.config.IncludeExtensions {
		if strings.ToLower(includeExt) == ext {
			return true
		}
	}
	return false
}

// isDefaultExcluded checks if a directory name is in the default exclude list.
func (f *Filter) isDefaultExcluded(name string) bool {
	// Case-insensitive on Windows
	if runtime.GOOS == "windows" {
		name = strings.ToLower(name)
		for _, excl := range f.defaultExcludes {
			if strings.ToLower(excl) == name {
				return true
			}
		}
		return false
	}

	for _, excl := range f.defaultExcludes {
		if excl == name {
			return true
		}
	}
	return false
}

// MatchesPattern is the public version of matchesPattern for testing.
func MatchesPattern(path, pattern string) bool {
	return matchesPattern(path, pattern)
}
