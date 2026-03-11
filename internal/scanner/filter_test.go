package scanner

import (
	"io/fs"
	"path/filepath"
	"testing"
	"time"
)

func TestNewFilter(t *testing.T) {
	config := NewConfig()
	filter := NewFilter(config, nil)

	if filter == nil {
		t.Fatal("NewFilter returned nil")
	}
	if filter.config != config {
		t.Error("Filter config not set correctly")
	}
	if len(filter.defaultExcludes) == 0 {
		t.Error("Filter should have default excludes")
	}
}

func TestFilter_ShouldSkip_HiddenFiles(t *testing.T) {
	tests := []struct {
		name          string
		includeHidden bool
		entryName     string
		isDir         bool
		wantSkip      bool
	}{
		{"hidden file excluded", false, ".gitignore", false, true},
		{"hidden file included", true, ".gitignore", false, false},
		{"hidden dir excluded", false, ".hidden", true, true},
		{"hidden dir included", true, ".hidden", true, false},
		{"normal file not skipped", false, "file.txt", false, false},
		{"normal dir not skipped", false, "subdir", true, false},
		// Note: .git is excluded by default excludes, not just hidden check
		{"git dir excluded by default", true, ".git", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{IncludeHidden: tt.includeHidden}
			filter := NewFilter(config, nil)

			entry := &mockDirEntry{name: tt.entryName, isDir: tt.isDir}
			got := filter.ShouldSkip("path", entry)

			if got != tt.wantSkip {
				t.Errorf("ShouldSkip() = %v, want %v", got, tt.wantSkip)
			}
		})
	}
}

func TestFilter_ShouldSkip_DefaultExcludes(t *testing.T) {
	defaultExcludes := []string{
		".git", ".svn", ".hg",
		"node_modules", "vendor",
		"__pycache__", ".pytest_cache",
		"target", "build", "dist",
		".idea", ".vscode",
		".DS_Store", "Thumbs.db",
	}

	for _, excl := range defaultExcludes {
		t.Run(excl, func(t *testing.T) {
			config := &Config{IncludeHidden: true} // Include hidden to test default excludes
			filter := NewFilter(config, nil)

			entry := &mockDirEntry{name: excl, isDir: true}
			got := filter.ShouldSkip(excl, entry)

			if !got {
				t.Errorf("ShouldSkip() for default exclude %q should be true", excl)
			}
		})
	}
}

func TestFilter_ShouldSkip_ExcludePatterns(t *testing.T) {
	tests := []struct {
		name            string
		excludePatterns []string
		relPath         string
		isDir           bool
		wantSkip        bool
	}{
		{
			name:            "exact match",
			excludePatterns: []string{"temp"},
			relPath:         "temp",
			isDir:           true,
			wantSkip:        true,
		},
		{
			name:            "glob pattern *.log",
			excludePatterns: []string{"*.log"},
			relPath:         "debug.log",
			isDir:           false,
			wantSkip:        true,
		},
		{
			name:            "glob pattern *.log with path",
			excludePatterns: []string{"*.log"},
			relPath:         "logs/debug.log",
			isDir:           false,
			wantSkip:        true,
		},
		{
			name:            "glob pattern no match",
			excludePatterns: []string{"*.log"},
			relPath:         "debug.txt",
			isDir:           false,
			wantSkip:        false,
		},
		{
			name:            "** pattern",
			excludePatterns: []string{"**/test/**"},
			relPath:         "a/b/test/file.txt",
			isDir:           false,
			wantSkip:        true,
		},
		{
			name:            "directory pattern",
			excludePatterns: []string{"node_modules/**"},
			relPath:         "node_modules/package.json",
			isDir:           false,
			wantSkip:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				IncludeHidden:   true,
				ExcludePatterns: tt.excludePatterns,
			}
			filter := NewFilter(config, nil)

			name := filepath.Base(tt.relPath)
			entry := &mockDirEntry{name: name, isDir: tt.isDir}
			got := filter.ShouldSkip(tt.relPath, entry)

			if got != tt.wantSkip {
				t.Errorf("ShouldSkip() = %v, want %v", got, tt.wantSkip)
			}
		})
	}
}

func TestFilter_ShouldSkip_IncludeOverridesExclude(t *testing.T) {
	config := &Config{
		IncludeHidden:   true,
		ExcludePatterns: []string{"*.log"},
		IncludePatterns: []string{"important.log"},
	}
	filter := NewFilter(config, nil)

	// This file matches exclude but also matches include
	entry := &mockDirEntry{name: "important.log", isDir: false}
	got := filter.ShouldSkip("important.log", entry)

	if got {
		t.Error("Include pattern should override exclude pattern")
	}

	// This file matches exclude and not include
	entry2 := &mockDirEntry{name: "debug.log", isDir: false}
	got2 := filter.ShouldSkip("debug.log", entry2)

	if !got2 {
		t.Error("File matching exclude but not include should be skipped")
	}
}

func TestFilter_ShouldInclude_Extensions(t *testing.T) {
	tests := []struct {
		name              string
		includeExtensions []string
		excludeExtensions []string
		includePatterns   []string
		relPath           string
		wantInclude       bool
	}{
		{
			name:              "include .go files only",
			includeExtensions: []string{".go"},
			relPath:           "main.go",
			wantInclude:       true,
		},
		{
			name:              "include .go files - other excluded",
			includeExtensions: []string{".go"},
			relPath:           "readme.md",
			wantInclude:       false,
		},
		{
			name:              "exclude .log files",
			excludeExtensions: []string{".log"},
			relPath:           "debug.log",
			wantInclude:       false,
		},
		{
			name:              "exclude .log - other included",
			excludeExtensions: []string{".log"},
			relPath:           "main.go",
			wantInclude:       true,
		},
		{
			name:              "include pattern overrides exclude extension",
			includePatterns:   []string{"*.log"},
			excludeExtensions: []string{".log"},
			relPath:           "important.log",
			wantInclude:       true,
		},
		{
			name:              "include extension overrides exclude extension",
			includeExtensions: []string{".log"},
			excludeExtensions: []string{".log"},
			relPath:           "important.log",
			wantInclude:       true,
		},
		{
			name:              "case insensitive extension",
			includeExtensions: []string{".GO"},
			relPath:           "main.go",
			wantInclude:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				IncludeHidden:     true,
				IncludeExtensions: tt.includeExtensions,
				ExcludeExtensions: tt.excludeExtensions,
				IncludePatterns:   tt.includePatterns,
			}
			filter := NewFilter(config, nil)

			entry := &mockDirEntry{name: filepath.Base(tt.relPath), isDir: false}
			got := filter.ShouldInclude(tt.relPath, entry)

			if got != tt.wantInclude {
				t.Errorf("ShouldInclude() = %v, want %v", got, tt.wantInclude)
			}
		})
	}
}

func TestFilter_ShouldInclude_Directories(t *testing.T) {
	config := &Config{}
	filter := NewFilter(config, nil)

	entry := &mockDirEntry{name: "subdir", isDir: true}
	got := filter.ShouldInclude("subdir", entry)

	if !got {
		t.Error("Directories should always be included if they pass ShouldSkip")
	}
}

func TestFilter_MatchesPattern_EdgeCases(t *testing.T) {
	tests := []struct {
		path    string
		pattern string
		want    bool
	}{
		// Empty patterns
		{"file.txt", "", false},
		// Special characters - [ and ] are glob metacharacters, so [1] matches "1"
		// {"file[1].txt", "file[1].txt", true}, // This would need escaping
		// Multiple wildcards
		{"a/b/c/d.txt", "**/c/*.txt", true},
		// Path with dots
		{"file.test.go", "*.go", true},
		// Unicode
		{"файл.txt", "*.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.path+"_"+tt.pattern, func(t *testing.T) {
			got := MatchesPattern(tt.path, tt.pattern)
			if got != tt.want {
				t.Errorf("MatchesPattern(%q, %q) = %v, want %v", tt.path, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestFilter_IsDefaultExcluded(t *testing.T) {
	filter := NewFilter(NewConfig(), nil)

	tests := []struct {
		name     string
		dirName  string
		wantExcl bool
	}{
		{"node_modules", "node_modules", true},
		{"build", "build", true},
		{"dist", "dist", true},
		{"normal dir", "src", false},
		{"partial match", "node_module", false},
		{"partial match 2", "build123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filter.isDefaultExcluded(tt.dirName)
			if got != tt.wantExcl {
				t.Errorf("isDefaultExcluded(%q) = %v, want %v", tt.dirName, got, tt.wantExcl)
			}
		})
	}
}

// mockDirEntry implements fs.DirEntry for testing
type mockDirEntry struct {
	name  string
	isDir bool
}

func (m *mockDirEntry) Name() string {
	return m.name
}

func (m *mockDirEntry) IsDir() bool {
	return m.isDir
}

func (m *mockDirEntry) Type() fs.FileMode {
	if m.isDir {
		return fs.ModeDir
	}
	return 0
}

func (m *mockDirEntry) Info() (fs.FileInfo, error) {
	return &mockFileInfo{name: m.name, isDir: m.isDir}, nil
}

// mockFileInfo implements fs.FileInfo for testing
type mockFileInfo struct {
	name  string
	isDir bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() fs.FileMode  { return 0 }
func (m *mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }
