package formatter

import (
	"strings"
	"testing"

	"github.com/wesbragagt/gps/pkg/types"
)

func TestToonFormatter_Format_BasicOutput(t *testing.T) {
	f := NewToonFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify project header
	if !strings.Contains(output, "project[test-project]") {
		t.Error("Missing project header")
	}
	if !strings.Contains(output, "type: go") {
		t.Error("Missing project type")
	}
	if !strings.Contains(output, "files: 5") {
		t.Error("Missing file count")
	}
	if !strings.Contains(output, "size: 9.8KB") {
		t.Errorf("Missing or wrong size, got: %s", output)
	}
	if !strings.Contains(output, "lines: 500") {
		t.Error("Missing line count")
	}
}

func TestToonFormatter_Format_ProjectHeader(t *testing.T) {
	f := NewToonFormatter()

	tests := []struct {
		name     string
		project  *types.Project
		contains []string
	}{
		{
			name:     "go project",
			project:  &types.Project{Name: "myapp", Type: types.ProjectTypeGo, Tree: types.Directory{Path: "."}},
			contains: []string{"project[myapp]", "type: go"},
		},
		{
			name:     "node project",
			project:  &types.Project{Name: "webapp", Type: types.ProjectTypeNode, Tree: types.Directory{Path: "."}},
			contains: []string{"project[webapp]", "type: node"},
		},
		{
			name:     "python project",
			project:  &types.Project{Name: "pyscript", Type: types.ProjectTypePython, Tree: types.Directory{Path: "."}},
			contains: []string{"project[pyscript]", "type: python"},
		},
		{
			name:     "rust project",
			project:  &types.Project{Name: "rustapp", Type: types.ProjectTypeRust, Tree: types.Directory{Path: "."}},
			contains: []string{"project[rustapp]", "type: rust"},
		},
		{
			name:     "java project",
			project:  &types.Project{Name: "javaapp", Type: types.ProjectTypeJava, Tree: types.Directory{Path: "."}},
			contains: []string{"project[javaapp]", "type: java"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := f.Format(tt.project)
			if err != nil {
				t.Fatalf("Format failed: %v", err)
			}
			for _, c := range tt.contains {
				if !strings.Contains(output, c) {
					t.Errorf("Missing expected content: %s", c)
				}
			}
		})
	}
}

func TestToonFormatter_Format_TreeStructure(t *testing.T) {
	f := NewToonFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify tree structure (root has item count [3])
	if !strings.Contains(output, "root[3]{") {
		t.Error("Missing root directory with item count")
	}
	if !strings.Contains(output, "main.go") {
		t.Error("Missing main.go")
	}
	if !strings.Contains(output, "internal[2]{") {
		t.Error("Missing internal directory with item count")
	}
	if !strings.Contains(output, "handler[1]{") {
		t.Error("Missing handler subdirectory with item count")
	}
}

func TestToonFormatter_Format_FileMetadata(t *testing.T) {
	f := NewToonFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify file metadata format: [size, lines, type]
	// Note: formatSize adds .0 for even numbers
	if !strings.Contains(output, "2.0KB") && !strings.Contains(output, "2KB") {
		t.Errorf("Missing size for main.go, got: %s", output)
	}
	if !strings.Contains(output, "100L, go") {
		t.Errorf("Missing lines and type for main.go, got: %s", output)
	}
	if !strings.Contains(output, "200B") {
		t.Errorf("Missing size for go.mod, got: %s", output)
	}
}

func TestToonFormatter_Format_KeyFiles(t *testing.T) {
	f := NewToonFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify keyfiles section
	if !strings.Contains(output, "keyfiles{") {
		t.Error("Missing keyfiles section")
	}
	if !strings.Contains(output, "entry: main.go") {
		t.Error("Missing entry point")
	}
	if !strings.Contains(output, "config: go.mod") {
		t.Error("Missing config file")
	}
	if !strings.Contains(output, "tests: main_test.go") {
		t.Error("Missing test file")
	}
	if !strings.Contains(output, "docs: README.md") {
		t.Error("Missing docs file")
	}
}

func TestToonFormatter_Format_KeyFilesMultiple(t *testing.T) {
	f := NewToonFormatter()
	project := CreateProjectWithMultipleKeyFiles()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Multiple files should show count
	if !strings.Contains(output, "entry: 3 files") {
		t.Error("Expected '3 files' for entry points")
	}
	if !strings.Contains(output, "config: 3 files") {
		t.Error("Expected '3 files' for configs")
	}
	if !strings.Contains(output, "tests: 2 files") {
		t.Error("Expected '2 files' for tests")
	}
	if !strings.Contains(output, "docs: 2 files") {
		t.Error("Expected '2 files' for docs")
	}
}

func TestToonFormatter_Format_EmptyFields(t *testing.T) {
	f := NewToonFormatter()
	project := CreateMinimalProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Minimal project should still have basic structure
	if !strings.Contains(output, "project[minimal]") {
		t.Error("Missing project header for minimal project")
	}
	// Should not have keyfiles section if empty
	if strings.Contains(output, "keyfiles{") {
		t.Error("Should not have keyfiles section for empty project")
	}
}

func TestToonFormatter_Format_BinaryFiles(t *testing.T) {
	f := NewToonFormatter()
	project := CreateProjectWithBinary()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Binary file should have 'bin' indicator
	if !strings.Contains(output, "bin") {
		t.Error("Missing binary indicator")
	}
	// Binary file should not have line count
	if strings.Contains(output, "0L, png") {
		t.Error("Binary file should not show 0L")
	}
}

func TestToonFormatter_SizeFormatting(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0B"},
		{500, "500B"},
		{1023, "1023B"},
		{1024, "1KB"},         // Exact value shows without decimal
		{1536, "1.5KB"},       // Fractional shows with decimal
		{2048, "2KB"},         // Exact value
		{1048575, "1024.0KB"}, // Just under 1MB shows with decimal
		{1048576, "1MB"},      // Exact 1MB
		{1572864, "1.5MB"},
		{2097152, "2MB"},         // Exact 2MB
		{1073741823, "1024.0MB"}, // Just under 1GB
		{1073741824, "1GB"},      // Exact 1GB
		{1610612736, "1.5GB"},
		{2147483648, "2GB"}, // Exact 2GB
		{5368709120, "5GB"}, // Exact 5GB
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatSize(%d) = %s, want %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestToonFormatter_Indentation(t *testing.T) {
	f := NewToonFormatter()
	project := CreateDeepNestedProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	lines := strings.Split(output, "\n")

	// Check that nested files are properly indented
	// Root level files have 2 spaces
	// Level 1 (a/) files have 4 spaces
	// Level 2 (a/b/) files have 6 spaces
	// Level 3 (a/b/c/) files have 8 spaces

	for _, line := range lines {
		// Count leading spaces
		indent := 0
		for _, c := range line {
			if c == ' ' {
				indent++
			} else {
				break
			}
		}

		// Check indentation for specific files
		if strings.Contains(line, "root.go") && indent != 2 {
			t.Errorf("root.go should have 2 spaces indent, got %d", indent)
		}
		if strings.Contains(line, "level1.go") && indent != 4 {
			t.Errorf("level1.go should have 4 spaces indent, got %d", indent)
		}
		if strings.Contains(line, "level2.go") && indent != 6 {
			t.Errorf("level2.go should have 6 spaces indent, got %d", indent)
		}
		if strings.Contains(line, "level3.go") && indent != 8 {
			t.Errorf("level3.go should have 8 spaces indent, got %d", indent)
		}
	}
}

func TestToonFormatter_NilProject(t *testing.T) {
	f := NewToonFormatter()
	output, err := f.Format(nil)
	if err != nil {
		t.Fatalf("Format(nil) failed: %v", err)
	}
	if output != "" {
		t.Errorf("Format(nil) = %q, want empty", output)
	}
}

func TestToonFormatter_NilDirectory(t *testing.T) {
	f := NewToonFormatter()
	result := f.formatTree(nil, 0)
	if result != "" {
		t.Errorf("formatTree(nil, 0) = %q, want empty", result)
	}
}

func TestToonFormatter_NilFile(t *testing.T) {
	f := NewToonFormatter()
	result := f.formatFile(nil)
	if result != "" {
		t.Errorf("formatFile(nil) = %q, want empty", result)
	}
}

func TestToonFormatter_NilKeyFiles(t *testing.T) {
	f := NewToonFormatter()
	result := f.formatKeyFiles(nil)
	if result != "" {
		t.Errorf("formatKeyFiles(nil) = %q, want empty", result)
	}
}

func TestToonFormatter_FormatFileWithUnknownType(t *testing.T) {
	f := NewToonFormatter()
	file := &types.File{
		Path:  "unknown.xyz",
		Size:  1000,
		Lines: 50,
		Type:  "unknown",
	}
	result := f.formatFile(file)
	// 'unknown' type should be omitted from metadata
	// The result should have size and lines but NOT the type
	if strings.Contains(result, "go]") || strings.Contains(result, ", unknown]") {
		t.Errorf("Unknown type should be omitted from metadata, got: %s", result)
	}
	// Should still have size and lines
	if !strings.Contains(result, "1000B") || !strings.Contains(result, "50L") {
		t.Errorf("Should have size and lines, got: %s", result)
	}
}

func TestToonFormatter_FormatFileWithEmptyType(t *testing.T) {
	f := NewToonFormatter()
	file := &types.File{
		Path:  "notype.txt",
		Size:  500,
		Lines: 25,
		Type:  "",
	}
	result := f.formatFile(file)
	// Empty type should be omitted
	if strings.Contains(result, ", ,") {
		t.Errorf("Empty type should be omitted, got: %s", result)
	}
}

func TestToonFormatter_HasKeyFiles(t *testing.T) {
	tests := []struct {
		name     string
		kf       *types.KeyFiles
		expected bool
	}{
		{"nil", nil, false},
		{"empty", &types.KeyFiles{}, false},
		{"with entry", &types.KeyFiles{EntryPoints: []string{"main.go"}}, true},
		{"with config", &types.KeyFiles{Configs: []string{"go.mod"}}, true},
		{"with test", &types.KeyFiles{Tests: []string{"test.go"}}, true},
		{"with doc", &types.KeyFiles{Docs: []string{"README.md"}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasKeyFiles(tt.kf)
			if result != tt.expected {
				t.Errorf("hasKeyFiles(%v) = %v, want %v", tt.kf, result, tt.expected)
			}
		})
	}
}

func TestToonFormatter_DirectoryWithItemCount(t *testing.T) {
	f := NewToonFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Directories should show item count
	// Root has 2 files + 1 subdir = 3 items
	if !strings.Contains(output, "root[3]") {
		t.Error("Root should show item count [3]")
	}
}

func TestToonFormatter_EmptyDirectory(t *testing.T) {
	f := NewToonFormatter()
	project := &types.Project{
		Name: "empty-dirs",
		Type: types.ProjectTypeGo,
		Tree: types.Directory{
			Path: ".",
			Subdirs: []types.Directory{
				{Path: "empty"},
			},
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Empty directory should not have count
	if strings.Contains(output, "empty[0]") {
		t.Error("Empty directory should not show [0]")
	}
}
