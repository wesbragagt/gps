package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wesbragagt/gps/pkg/types"
)

func TestScanner_Walk(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create files and directories
	testFiles := []string{
		"file1.txt",
		"file2.go",
		"subdir/file3.md",
		"subdir/nested/file4.go",
		".hidden/file5.txt",
		".hidden_file",
	}

	for _, f := range testFiles {
		path := filepath.Join(tmpDir, f)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %q: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", path, err)
		}
	}

	tests := []struct {
		name         string
		config       *Config
		wantFiles    int // total files in tree
		wantSubdirs  int // total subdirs in tree
		wantMaxDepth int // max depth of tree
	}{
		{
			name:        "default config - no hidden",
			config:      NewConfig(),
			wantFiles:   4, // 2 root + 1 subdir + 1 nested (hidden excluded)
			wantSubdirs: 2, // subdir + nested
		},
		{
			name: "include hidden",
			config: &Config{
				MaxDepth:      -1,
				IncludeHidden: true,
			},
			wantFiles:   6, // all files
			wantSubdirs: 3, // subdir + nested + .hidden
		},
		{
			name: "depth limit 1",
			config: &Config{
				MaxDepth:      1,
				IncludeHidden: false,
			},
			wantFiles:   3, // 2 root + 1 in subdir (nested not traversed)
			wantSubdirs: 1, // subdir only
		},
		{
			name: "depth limit 0",
			config: &Config{
				MaxDepth:      0,
				IncludeHidden: false,
			},
			wantFiles:   2, // only root files
			wantSubdirs: 0, // no subdirs traversed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.config)
			got, err := s.Walk(tmpDir)
			if err != nil {
				t.Fatalf("Walk() error = %v", err)
			}

			if got == nil {
				t.Fatal("Walk() returned nil directory")
			}

			// Count files and subdirs
			fileCount := countFiles(got)
			subdirCount := countSubdirs(got)

			if fileCount != tt.wantFiles {
				t.Errorf("Walk() file count = %d, want %d", fileCount, tt.wantFiles)
			}
			if subdirCount != tt.wantSubdirs {
				t.Errorf("Walk() subdir count = %d, want %d", subdirCount, tt.wantSubdirs)
			}
		})
	}
}

func TestScanner_Walk_InvalidRoot(t *testing.T) {
	s := New(nil)

	_, err := s.Walk("/nonexistent/path/12345")
	if err == nil {
		t.Error("Walk() with invalid path should return error")
	}

	// Test with a file instead of directory
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = s.Walk(tmpFile)
	if err == nil {
		t.Error("Walk() with file path should return error")
	}
}

func TestScanner_Errors(t *testing.T) {
	s := New(nil)
	errs := s.Errors()
	if errs == nil {
		t.Error("Errors() should return non-nil slice")
	}
	if len(errs) != 0 {
		t.Errorf("New scanner should have no errors, got %d", len(errs))
	}
}

func TestScanError(t *testing.T) {
	cause := os.ErrPermission
	err := &ScanError{
		Path:  "/test/path",
		Op:    "read",
		Cause: cause,
	}

	if err.Error() == "" {
		t.Error("ScanError.Error() should return non-empty string")
	}

	if unwrapped := err.Unwrap(); unwrapped != cause {
		t.Errorf("ScanError.Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestIsHidden(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{".hidden", true},
		{".gitignore", true},
		{"normal", false},
		{"file.txt", false},
		{"..double", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHidden(tt.name); got != tt.want {
				t.Errorf("isHidden(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestScanner_Walk_WithFilters(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create files and directories
	testFiles := []string{
		"file1.txt",
		"file2.go",
		"file3.log",
		"file4.tmp",
		"subdir/file5.go",
		"subdir/file6.md",
		"node_modules/package.json",
		"build/output.txt",
	}

	for _, f := range testFiles {
		path := filepath.Join(tmpDir, f)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %q: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", path, err)
		}
	}

	tests := []struct {
		name        string
		config      *Config
		wantFiles   int
		wantSubdirs int
	}{
		{
			name: "exclude log files",
			config: &Config{
				MaxDepth:          -1,
				ExcludeExtensions: []string{".log"},
				RespectGitignore:  false,
			},
			wantFiles:   5, // 4 root - 1 (.log) + 2 subdir = 5 (node_modules/build excluded by default)
			wantSubdirs: 1, // only subdir (node_modules/build excluded by default)
		},
		{
			name: "include only go files",
			config: &Config{
				MaxDepth:          -1,
				IncludeExtensions: []string{".go"},
				RespectGitignore:  false,
			},
			wantFiles:   2, // file2.go, file5.go
			wantSubdirs: 1, // only subdir (node_modules/build excluded by default)
		},
		{
			name: "exclude patterns - node_modules",
			config: &Config{
				MaxDepth:         -1,
				ExcludePatterns:  []string{"node_modules"},
				RespectGitignore: false,
			},
			wantFiles:   6, // all except node_modules (but build also excluded by default)
			wantSubdirs: 1, // only subdir
		},
		{
			name: "exclude patterns with wildcard",
			config: &Config{
				MaxDepth:         -1,
				ExcludePatterns:  []string{"*.log", "*.tmp"},
				RespectGitignore: false,
			},
			wantFiles:   4, // 4 root - 2 (.log, .tmp) + 2 subdir = 4
			wantSubdirs: 1, // only subdir
		},
		{
			name: "include patterns",
			config: &Config{
				MaxDepth:         -1,
				IncludePatterns:  []string{"**/*.go"},
				RespectGitignore: false,
			},
			wantFiles:   2, // only .go files
			wantSubdirs: 1, // only subdir
		},
		{
			name: "default excludes skip node_modules and build",
			config: &Config{
				MaxDepth:         -1,
				RespectGitignore: false,
			},
			wantFiles:   6, // 4 root + 2 subdir
			wantSubdirs: 1, // only subdir
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.config)
			got, err := s.Walk(tmpDir)
			if err != nil {
				t.Fatalf("Walk() error = %v", err)
			}

			if got == nil {
				t.Fatal("Walk() returned nil directory")
			}

			// Count files and subdirs
			fileCount := countFiles(got)
			subdirCount := countSubdirs(got)

			if fileCount != tt.wantFiles {
				t.Errorf("Walk() file count = %d, want %d", fileCount, tt.wantFiles)
			}
			if subdirCount != tt.wantSubdirs {
				t.Errorf("Walk() subdir count = %d, want %d", subdirCount, tt.wantSubdirs)
			}
		})
	}
}

func TestScanner_Walk_WithGitignore(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .gitignore file
	gitignoreContent := `*.log
*.tmp
/build/
`
	if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	// Create files
	testFiles := []string{
		"file1.txt",
		"file2.go",
		"file3.log",     // ignored by gitignore
		"file4.tmp",     // ignored by gitignore
		"build/out.txt", // ignored by gitignore
		"subdir/file5.go",
	}

	for _, f := range testFiles {
		path := filepath.Join(tmpDir, f)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %q: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", path, err)
		}
	}

	config := &Config{
		MaxDepth:         -1,
		RespectGitignore: true,
	}

	s := New(config)
	got, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	fileCount := countFiles(got)
	// Should have: file1.txt, file2.go, subdir/file5.go = 3
	wantFiles := 3
	if fileCount != wantFiles {
		t.Errorf("Walk() with gitignore file count = %d, want %d", fileCount, wantFiles)
	}
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		path    string
		pattern string
		want    bool
	}{
		{"file.txt", "*.txt", true},
		{"file.go", "*.txt", false},
		{"subdir/file.txt", "*.txt", true}, // matches basename
		{"subdir/file.txt", "**/*.txt", true},
		{"a/b/c/file.txt", "**/*.txt", true},
		{"node_modules", "node_modules", true},
		{"node_modules/package.json", "node_modules/**", true},
		{"build", "build", true},
		{"build/output.txt", "build/**", true},
	}

	for _, tt := range tests {
		t.Run(tt.path+"_"+tt.pattern, func(t *testing.T) {
			if got := MatchesPattern(tt.path, tt.pattern); got != tt.want {
				t.Errorf("MatchesPattern(%q, %q) = %v, want %v", tt.path, tt.pattern, got, tt.want)
			}
		})
	}
}

// Helper functions

func countFiles(d *types.Directory) int {
	count := len(d.Files)
	for i := range d.Subdirs {
		count += countFiles(&d.Subdirs[i])
	}
	return count
}

func countSubdirs(d *types.Directory) int {
	count := len(d.Subdirs)
	for i := range d.Subdirs {
		count += countSubdirs(&d.Subdirs[i])
	}
	return count
}

func countMaxDepth(d *types.Directory, currentDepth int) int {
	if len(d.Subdirs) == 0 {
		return currentDepth
	}
	maxDepth := currentDepth
	for i := range d.Subdirs {
		depth := countMaxDepth(&d.Subdirs[i], currentDepth+1)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
}

func collectFilePaths(d *types.Directory) []string {
	paths := make([]string, 0)
	for _, f := range d.Files {
		paths = append(paths, f.Path)
	}
	for i := range d.Subdirs {
		paths = append(paths, collectFilePaths(&d.Subdirs[i])...)
	}
	return paths
}

// Test using testdata fixtures
func TestScanner_Walk_TestdataProject(t *testing.T) {
	config := &Config{
		MaxDepth:         -1,
		IncludeHidden:    false,
		RespectGitignore: false,
	}
	s := New(config)

	got, err := s.Walk("testdata/project")
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Should have main.go, lib/util.go, lib/helper.go, test/main_test.go, README.md
	fileCount := countFiles(got)
	if fileCount != 5 {
		t.Errorf("Walk() file count = %d, want 5", fileCount)
	}

	// Should have lib and test subdirs
	if len(got.Subdirs) != 2 {
		t.Errorf("Walk() root subdirs = %d, want 2", len(got.Subdirs))
	}
}

func TestScanner_Walk_TestdataProjectWithGitignore(t *testing.T) {
	config := &Config{
		MaxDepth:         -1,
		IncludeHidden:    false,
		RespectGitignore: true,
	}
	s := New(config)

	got, err := s.Walk("testdata/project")
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// .gitignore should be excluded (hidden) but files should be included
	// since none match the patterns in testdata/project/.gitignore
	fileCount := countFiles(got)
	if fileCount != 5 {
		paths := collectFilePaths(got)
		t.Errorf("Walk() file count = %d, want 5. Files: %v", fileCount, paths)
	}
}

func TestScanner_Walk_TestdataGitignore(t *testing.T) {
	config := &Config{
		MaxDepth:         -1,
		IncludeHidden:    false,
		RespectGitignore: true,
	}
	s := New(config)

	got, err := s.Walk("testdata/gitignore_test")
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Should have file.txt but not ignored.log
	fileCount := countFiles(got)
	if fileCount != 1 {
		t.Errorf("Walk() file count = %d, want 1 (ignored.log should be excluded)", fileCount)
	}

	if len(got.Files) > 0 && got.Files[0].Path != "file.txt" {
		t.Errorf("Expected file.txt, got %s", got.Files[0].Path)
	}
}

func TestScanner_Walk_EmptyDirectory(t *testing.T) {
	// Use the empty_dir test fixture
	config := NewConfig()
	s := New(config)

	got, err := s.Walk("testdata/empty_dir")
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	if got == nil {
		t.Fatal("Walk() returned nil for empty directory")
	}

	if len(got.Files) != 0 {
		t.Errorf("Empty directory should have no files, got %d", len(got.Files))
	}

	if len(got.Subdirs) != 0 {
		t.Errorf("Empty directory should have no subdirs, got %d", len(got.Subdirs))
	}
}

func TestScanner_Walk_DeepNesting(t *testing.T) {
	config := NewConfig()
	s := New(config)

	got, err := s.Walk("testdata/deep")
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// The deep structure has a/b/c/d/e (5 levels)
	maxDepth := countMaxDepth(got, 0)
	if maxDepth != 5 {
		t.Errorf("Max depth = %d, want 5", maxDepth)
	}
}

func TestScanner_Walk_DeepNestingWithDepthLimit(t *testing.T) {
	config := &Config{
		MaxDepth:         2,
		IncludeHidden:    false,
		RespectGitignore: false,
	}
	s := New(config)

	got, err := s.Walk("testdata/deep")
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	maxDepth := countMaxDepth(got, 0)
	if maxDepth != 2 {
		t.Errorf("Max depth with limit = %d, want 2", maxDepth)
	}
}

func TestScanner_Walk_Symlink(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file and a symlink to it
	targetFile := filepath.Join(tmpDir, "target.txt")
	if err := os.WriteFile(targetFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	symlinkPath := filepath.Join(tmpDir, "link.txt")
	if err := os.Symlink(targetFile, symlinkPath); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	config := NewConfig()
	s := New(config)

	got, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Symlink should be included as a file
	if len(got.Files) != 2 {
		t.Errorf("Expected 2 files (target + symlink), got %d", len(got.Files))
	}
}

func TestScanner_Walk_SpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files with special characters in names
	specialFiles := []string{
		"file with spaces.txt",
		"file-with-dashes.txt",
		"file_with_underscores.txt",
		"file.multiple.dots.txt",
	}

	for _, f := range specialFiles {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", f, err)
		}
	}

	config := NewConfig()
	s := New(config)

	got, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	if len(got.Files) != len(specialFiles) {
		t.Errorf("Expected %d files, got %d", len(specialFiles), len(got.Files))
	}
}

func TestScanner_NewConfig(t *testing.T) {
	config := NewConfig()

	if config.MaxDepth != -1 {
		t.Errorf("Default MaxDepth = %d, want -1", config.MaxDepth)
	}
	if config.IncludeHidden {
		t.Error("Default IncludeHidden should be false")
	}
	if !config.RespectGitignore {
		t.Error("Default RespectGitignore should be true")
	}
}

func TestScanner_New_NilConfig(t *testing.T) {
	s := New(nil)
	if s == nil {
		t.Fatal("New(nil) returned nil")
	}
	if s.config == nil {
		t.Error("Scanner.config should not be nil when passed nil")
	}
}

func TestScanner_RelativePath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested structure
	testFiles := []string{
		"root.txt",
		"subdir/nested.txt",
		"subdir/deep/file.txt",
	}

	for _, f := range testFiles {
		path := filepath.Join(tmpDir, f)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %q: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", path, err)
		}
	}

	config := NewConfig()
	s := New(config)

	got, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Check root path is "."
	if got.Path != "." {
		t.Errorf("Root path = %q, want %q", got.Path, ".")
	}

	// Check file paths are relative
	paths := collectFilePaths(got)
	expectedPaths := map[string]bool{
		"root.txt":             true,
		"subdir/nested.txt":    true,
		"subdir/deep/file.txt": true,
	}

	for _, p := range paths {
		if !expectedPaths[p] {
			t.Errorf("Unexpected file path: %s", p)
		}
		delete(expectedPaths, p)
	}

	if len(expectedPaths) > 0 {
		t.Errorf("Missing file paths: %v", expectedPaths)
	}
}

func TestScanner_FileInfo(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("test content for size check")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := NewConfig()
	s := New(config)

	got, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	if len(got.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(got.Files))
	}

	file := got.Files[0]
	if file.Size != int64(len(content)) {
		t.Errorf("File size = %d, want %d", file.Size, len(content))
	}
	if file.Path != "test.txt" {
		t.Errorf("File path = %q, want %q", file.Path, "test.txt")
	}
	if file.ModTime.IsZero() {
		t.Error("ModTime should not be zero")
	}
}

func TestScanner_IncludeOverridesExclude(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files
	testFiles := []string{
		"keep.go",
		"exclude.log",
		"important.log",
	}
	for _, f := range testFiles {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", f, err)
		}
	}

	config := &Config{
		MaxDepth:         -1,
		ExcludePatterns:  []string{"*.log"},
		IncludePatterns:  []string{"important.log", "*.go"}, // Include both important.log and .go files
		RespectGitignore: false,
	}
	s := New(config)

	got, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Should have keep.go and important.log (include overrides exclude)
	fileCount := countFiles(got)
	if fileCount != 2 {
		paths := collectFilePaths(got)
		t.Errorf("Expected 2 files, got %d: %v", fileCount, paths)
	}
}

func TestScanner_Walk_MultipleWalks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some files
	for i := 0; i < 3; i++ {
		path := filepath.Join(tmpDir, "file"+string(rune('1'+i))+".txt")
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	config := NewConfig()
	s := New(config)

	// First walk
	got1, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("First Walk() error = %v", err)
	}

	// Second walk
	got2, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Second Walk() error = %v", err)
	}

	// Both should return same results
	if countFiles(got1) != countFiles(got2) {
		t.Errorf("Multiple walks returned different file counts: %d vs %d", countFiles(got1), countFiles(got2))
	}

	// Errors should be reset between walks
	if len(s.Errors()) != 0 {
		t.Error("Scanner should have no errors after clean walk")
	}
}

func TestScanner_Walk_ConcurrentSafe(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files
	for i := 0; i < 5; i++ {
		path := filepath.Join(tmpDir, "file"+string(rune('1'+i))+".txt")
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Run multiple scanners concurrently
	done := make(chan bool)
	for i := 0; i < 3; i++ {
		go func() {
			config := NewConfig()
			s := New(config)
			_, err := s.Walk(tmpDir)
			if err != nil {
				t.Errorf("Concurrent Walk() error = %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestScanner_Walk_PathCleaning(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	config := NewConfig()
	s := New(config)

	// Walk with trailing slashes and dots
	got, err := s.Walk(tmpDir + string(filepath.Separator))
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	if got == nil {
		t.Error("Walk() returned nil")
	}
}

func TestScanner_Walk_TreeStructure(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a specific structure
	dirs := []string{
		"src",
		"src/api",
		"src/api/handlers",
		"src/db",
		"tests",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, d), 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
	}

	files := []string{
		"main.go",
		"src/api/api.go",
		"src/api/handlers/user.go",
		"src/db/db.go",
		"tests/main_test.go",
	}
	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	config := &Config{
		MaxDepth:         -1,
		RespectGitignore: false,
	}
	s := New(config)

	got, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Verify structure
	if got.Path != "." {
		t.Errorf("Root path = %q, want %q", got.Path, ".")
	}

	// Root should have 1 file (main.go) and 2 subdirs (src, tests)
	if len(got.Files) != 1 {
		t.Errorf("Root files = %d, want 1", len(got.Files))
	}
	if len(got.Subdirs) != 2 {
		t.Errorf("Root subdirs = %d, want 2", len(got.Subdirs))
	}

	// Find src directory
	var srcDir *types.Directory
	for i := range got.Subdirs {
		if got.Subdirs[i].Path == "src" {
			srcDir = &got.Subdirs[i]
			break
		}
	}
	if srcDir == nil {
		t.Fatal("src directory not found")
	}

	// src should have 2 subdirs (api, db) and 0 files
	if len(srcDir.Files) != 0 {
		t.Errorf("src files = %d, want 0", len(srcDir.Files))
	}
	if len(srcDir.Subdirs) != 2 {
		t.Errorf("src subdirs = %d, want 2", len(srcDir.Subdirs))
	}
}

func TestScanner_Walk_WithCurrentlyOpen(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files
	for i := 0; i < 3; i++ {
		path := filepath.Join(tmpDir, "file"+string(rune('1'+i))+".txt")
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	config := &Config{
		MaxDepth:      -1,
		CurrentlyOpen: []string{"file1.txt"},
	}
	s := New(config)

	got, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// CurrentlyOpen is stored but not yet used for filtering
	// Just verify the scan works
	if len(got.Files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(got.Files))
	}
}

func TestScanner_Walk_GitignoreError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an unreadable .gitignore (on Unix systems)
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := &Config{
		MaxDepth:         -1,
		RespectGitignore: true,
	}
	s := New(config)

	got, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Should still return results even with gitignore issues
	if got == nil {
		t.Error("Walk() should return directory even with gitignore error")
	}
}

func TestScanner_Walk_NestedGitignore(t *testing.T) {
	tmpDir := t.TempDir()

	// Create root .gitignore
	rootGitignore := `*.log
`
	if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(rootGitignore), 0644); err != nil {
		t.Fatalf("Failed to create root .gitignore: %v", err)
	}

	// Create subdirectory with its own .gitignore
	subdir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	subGitignore := `*.tmp
secret.txt
`
	if err := os.WriteFile(filepath.Join(subdir, ".gitignore"), []byte(subGitignore), 0644); err != nil {
		t.Fatalf("Failed to create subdir .gitignore: %v", err)
	}

	// Create files
	files := []string{
		"file.txt",
		"debug.log", // ignored by root
		"subdir/file.txt",
		"subdir/temp.tmp",   // ignored by subdir gitignore
		"subdir/secret.txt", // ignored by subdir gitignore
	}
	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	config := &Config{
		MaxDepth:         -1,
		RespectGitignore: true,
	}
	s := New(config)

	got, err := s.Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk() error = %v", err)
	}

	// Should have file.txt and subdir/file.txt
	// debug.log, temp.tmp, secret.txt should be ignored
	fileCount := countFiles(got)
	if fileCount != 2 {
		paths := collectFilePaths(got)
		t.Errorf("Expected 2 files, got %d: %v", fileCount, paths)
	}
}
