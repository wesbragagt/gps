package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewGitignoreMatcher(t *testing.T) {
	matcher := NewGitignoreMatcher("/some/path")

	if matcher == nil {
		t.Fatal("NewGitignoreMatcher returned nil")
	}
	if matcher.root != "/some/path" {
		t.Errorf("root = %q, want %q", matcher.root, "/some/path")
	}
}

func TestGitignoreMatcher_Load(t *testing.T) {
	t.Run("no gitignore file", func(t *testing.T) {
		tmpDir := t.TempDir()
		matcher := NewGitignoreMatcher(tmpDir)

		err := matcher.Load()
		if err != nil {
			t.Errorf("Load() with no .gitignore should not error, got: %v", err)
		}
	})

	t.Run("valid gitignore file", func(t *testing.T) {
		tmpDir := t.TempDir()

		gitignoreContent := `*.log
*.tmp
/build/
`
		if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
			t.Fatalf("Failed to create .gitignore: %v", err)
		}

		matcher := NewGitignoreMatcher(tmpDir)
		err := matcher.Load()
		if err != nil {
			t.Errorf("Load() error = %v", err)
		}

		if matcher.rootIgnore == nil {
			t.Error("rootIgnore should not be nil after Load()")
		}
	})

	t.Run("gitignore with comments and empty lines", func(t *testing.T) {
		tmpDir := t.TempDir()

		gitignoreContent := `# This is a comment

*.log

# Another comment
*.tmp
`
		if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
			t.Fatalf("Failed to create .gitignore: %v", err)
		}

		matcher := NewGitignoreMatcher(tmpDir)
		err := matcher.Load()
		if err != nil {
			t.Errorf("Load() error = %v", err)
		}
	})
}

func TestGitignoreMatcher_Matches(t *testing.T) {
	t.Run("basic patterns", func(t *testing.T) {
		tmpDir := t.TempDir()

		gitignoreContent := `*.log
*.tmp
/build/
secret.txt
`
		if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
			t.Fatalf("Failed to create .gitignore: %v", err)
		}

		matcher := NewGitignoreMatcher(tmpDir)
		if err := matcher.Load(); err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		tests := []struct {
			path      string
			isDir     bool
			wantMatch bool
		}{
			{"debug.log", false, true},
			{"logs/debug.log", false, true},
			{"file.txt", false, false},
			{"temp.tmp", false, true},
			// Note: /build/ pattern matches files inside build, not the dir name itself in this library
			{"build/output.txt", false, true},
			{"secret.txt", false, true},
			{"src/main.go", false, false},
		}

		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				got := matcher.Matches(tt.path, tt.isDir)
				if got != tt.wantMatch {
					t.Errorf("Matches(%q, %v) = %v, want %v", tt.path, tt.isDir, got, tt.wantMatch)
				}
			})
		}
	})

	t.Run("negation patterns", func(t *testing.T) {
		tmpDir := t.TempDir()

		gitignoreContent := `*.log
!important.log
`
		if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
			t.Fatalf("Failed to create .gitignore: %v", err)
		}

		matcher := NewGitignoreMatcher(tmpDir)
		if err := matcher.Load(); err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		// debug.log should be ignored
		if !matcher.Matches("debug.log", false) {
			t.Error("debug.log should match")
		}

		// important.log should NOT be ignored (negation)
		if matcher.Matches("important.log", false) {
			t.Error("important.log should NOT match due to negation")
		}
	})

	t.Run("directory patterns", func(t *testing.T) {
		tmpDir := t.TempDir()

		gitignoreContent := `node_modules/
dist/
*.egg-info/
`
		if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
			t.Fatalf("Failed to create .gitignore: %v", err)
		}

		matcher := NewGitignoreMatcher(tmpDir)
		if err := matcher.Load(); err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		tests := []struct {
			path      string
			isDir     bool
			wantMatch bool
		}{
			// Files inside these directories are matched
			{"node_modules/package.json", false, true},
			{"dist/bundle.js", false, true},
		}

		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				got := matcher.Matches(tt.path, tt.isDir)
				if got != tt.wantMatch {
					t.Errorf("Matches(%q, %v) = %v, want %v", tt.path, tt.isDir, got, tt.wantMatch)
				}
			})
		}
	})

	t.Run("wildcard patterns", func(t *testing.T) {
		tmpDir := t.TempDir()

		gitignoreContent := `*.log
test_*.go
.*.swp
`
		if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
			t.Fatalf("Failed to create .gitignore: %v", err)
		}

		matcher := NewGitignoreMatcher(tmpDir)
		if err := matcher.Load(); err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		tests := []struct {
			path      string
			wantMatch bool
		}{
			{"debug.log", true},
			{"test_main.go", true},
			{"test_utils.go", true},
			{"main_test.go", false}, // pattern is test_*.go, not *_test.go
			{".file.swp", true},
			{".other.swp", true},
		}

		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				got := matcher.Matches(tt.path, false)
				if got != tt.wantMatch {
					t.Errorf("Matches(%q) = %v, want %v", tt.path, got, tt.wantMatch)
				}
			})
		}
	})
}

func TestGitignoreMatcher_SubdirGitignore(t *testing.T) {
	tmpDir := t.TempDir()

	// Create root .gitignore
	rootGitignore := `*.log
`
	if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(rootGitignore), 0644); err != nil {
		t.Fatalf("Failed to create root .gitignore: %v", err)
	}

	// Create subdirectory
	subdir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Create subdirectory .gitignore
	subGitignore := `*.tmp
`
	if err := os.WriteFile(filepath.Join(subdir, ".gitignore"), []byte(subGitignore), 0644); err != nil {
		t.Fatalf("Failed to create subdir .gitignore: %v", err)
	}

	matcher := NewGitignoreMatcher(tmpDir)
	if err := matcher.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Check that root gitignore works
	if !matcher.Matches("debug.log", false) {
		t.Error("debug.log should match root gitignore")
	}

	// Check for subdir gitignore existence
	if !matcher.HasSubdirGitignore("subdir") {
		t.Error("HasSubdirGitignore should return true for subdir")
	}

	// Load subdir gitignore
	if err := matcher.LoadSubdirGitignore("subdir"); err != nil {
		t.Errorf("LoadSubdirGitignore error = %v", err)
	}
}

func TestGitignoreMatcher_HasSubdirGitignore(t *testing.T) {
	tmpDir := t.TempDir()

	matcher := NewGitignoreMatcher(tmpDir)

	// No gitignore exists
	if matcher.HasSubdirGitignore("nonexistent") {
		t.Error("HasSubdirGitignore should return false for nonexistent dir")
	}

	// Create subdirectory with gitignore
	subdir := filepath.Join(tmpDir, "withgitignore")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subdir, ".gitignore"), []byte("*.log\n"), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	if !matcher.HasSubdirGitignore("withgitignore") {
		t.Error("HasSubdirGitignore should return true for dir with .gitignore")
	}
}

func TestGitignoreMatcher_LoadSubdirGitignore(t *testing.T) {
	t.Run("nonexistent directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		matcher := NewGitignoreMatcher(tmpDir)

		err := matcher.LoadSubdirGitignore("nonexistent")
		if err != nil {
			t.Errorf("LoadSubdirGitignore for nonexistent dir should not error, got: %v", err)
		}
	})

	t.Run("directory without gitignore", func(t *testing.T) {
		tmpDir := t.TempDir()
		subdir := filepath.Join(tmpDir, "subdir")
		if err := os.Mkdir(subdir, 0755); err != nil {
			t.Fatalf("Failed to create subdir: %v", err)
		}

		matcher := NewGitignoreMatcher(tmpDir)
		err := matcher.LoadSubdirGitignore("subdir")
		if err != nil {
			t.Errorf("LoadSubdirGitignore for dir without .gitignore should not error, got: %v", err)
		}
	})
}

func TestParseGitignorePatterns(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "basic patterns",
			input: "*.log\n*.tmp\n",
			want:  []string{"*.log", "*.tmp"},
		},
		{
			name:  "with comments",
			input: "# comment\n*.log\n# another\n*.tmp\n",
			want:  []string{"*.log", "*.tmp"},
		},
		{
			name:  "with empty lines",
			input: "\n*.log\n\n*.tmp\n\n",
			want:  []string{"*.log", "*.tmp"},
		},
		{
			name:  "empty input",
			input: "",
			want:  nil,
		},
		{
			name:  "only comments",
			input: "# comment\n# another\n",
			want:  nil,
		},
		{
			name:  "negation patterns",
			input: "*.log\n!important.log\n",
			want:  []string{"*.log", "!important.log"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGitignorePatterns([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGitignorePatterns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ParseGitignorePatterns() = %v, want %v", got, tt.want)
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("ParseGitignorePatterns()[%d] = %v, want %v", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestDefaultIgnorePatterns(t *testing.T) {
	if len(DefaultIgnorePatterns) == 0 {
		t.Error("DefaultIgnorePatterns should not be empty")
	}

	// Check some essential patterns
	essentialPatterns := []string{".git", "node_modules", "*.log"}
	for _, p := range essentialPatterns {
		found := false
		for _, dp := range DefaultIgnorePatterns {
			if dp == p {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("DefaultIgnorePatterns missing essential pattern: %s", p)
		}
	}
}

func TestGitignoreMatcher_UsingTestdata(t *testing.T) {
	// Use the testdata/gitignore_test fixture
	testdataPath := "testdata/gitignore_test"

	matcher := NewGitignoreMatcher(testdataPath)
	if err := matcher.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// file.txt should not be ignored
	if matcher.Matches("file.txt", false) {
		t.Error("file.txt should not be ignored")
	}

	// ignored.log should be ignored
	if !matcher.Matches("ignored.log", false) {
		t.Error("ignored.log should be ignored")
	}
}
