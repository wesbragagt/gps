package scanner

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"

	ignore "github.com/sabhiram/go-gitignore"
)

// GitignoreMatcher manages gitignore patterns for a project.
// It supports nested .gitignore files in subdirectories.
//
// The matcher loads .gitignore files lazily and caches the results.
// It supports the standard gitignore syntax including:
//   - Pattern negation (!)
//   - Directory-specific patterns (trailing /)
//   - Wildcards (* and **)
//   - Root-relative patterns (leading /)
type GitignoreMatcher struct {
	root       string
	rootIgnore *ignore.GitIgnore
	// Cache of subdirectory gitignore patterns
	subdirIgnores map[string]*ignore.GitIgnore
	mu            sync.RWMutex
}

// NewGitignoreMatcher creates a new matcher for the given root directory.
func NewGitignoreMatcher(root string) *GitignoreMatcher {
	return &GitignoreMatcher{
		root:          root,
		subdirIgnores: make(map[string]*ignore.GitIgnore),
	}
}

// Load reads .gitignore from the root directory.
func (g *GitignoreMatcher) Load() error {
	gitignorePath := filepath.Join(g.root, ".gitignore")
	ign, err := ignore.CompileIgnoreFile(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No .gitignore file is fine
		}
		return err
	}

	g.mu.Lock()
	g.rootIgnore = ign
	g.mu.Unlock()

	return nil
}

// Matches checks if a path should be ignored based on gitignore patterns.
// The path should be relative to the project root.
func (g *GitignoreMatcher) Matches(relPath string, isDir bool) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Normalize path for matching (gitignore uses forward slashes)
	relPath = filepath.ToSlash(relPath)

	// Check root gitignore first
	if g.rootIgnore != nil && g.rootIgnore.MatchesPath(relPath) {
		return true
	}

	// Check for gitignore in parent directories
	dir := filepath.Dir(relPath)
	for dir != "." && dir != "" {
		if ign, ok := g.subdirIgnores[dir]; ok {
			// Adjust path relative to the subdirectory gitignore
			subRelPath := strings.TrimPrefix(relPath, dir+"/")
			if ign.MatchesPath(subRelPath) {
				return true
			}
		}
		dir = filepath.Dir(dir)
	}

	return false
}

// LoadSubdirGitignore loads a .gitignore file from a subdirectory.
// This should be called when traversing into new directories.
func (g *GitignoreMatcher) LoadSubdirGitignore(relPath string) error {
	gitignorePath := filepath.Join(g.root, relPath, ".gitignore")
	ign, err := ignore.CompileIgnoreFile(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No .gitignore in this directory
		}
		return err
	}

	g.mu.Lock()
	g.subdirIgnores[relPath] = ign
	g.mu.Unlock()

	return nil
}

// HasSubdirGitignore checks if a subdirectory has a .gitignore file.
func (g *GitignoreMatcher) HasSubdirGitignore(relPath string) bool {
	gitignorePath := filepath.Join(g.root, relPath, ".gitignore")
	_, err := os.Stat(gitignorePath)
	return err == nil
}

// ParseGitignorePatterns parses gitignore-style patterns from a byte slice.
// This is useful for testing or custom pattern lists.
func ParseGitignorePatterns(data []byte) ([]string, error) {
	var patterns []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns, scanner.Err()
}

// DefaultIgnorePatterns returns common patterns that should typically be ignored.
var DefaultIgnorePatterns = []string{
	// Version control
	".git", ".svn", ".hg",
	// Dependencies
	"node_modules", "vendor",
	// Python
	"__pycache__", ".pytest_cache", "*.pyc", "*.pyo",
	// Build outputs
	"target", "build", "dist", "out", "bin",
	// IDE
	".idea", ".vscode", "*.swp", "*.swo",
	// OS files
	".DS_Store", "Thumbs.db",
	// Logs
	"*.log",
	// Environment
	".env", ".env.local",
}
