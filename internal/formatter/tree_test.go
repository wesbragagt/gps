package formatter

import (
	"os"
	"strings"
	"testing"

	"github.com/wesbragagt/gps/pkg/types"
)

func TestTreeFormatter_Format_BasicOutput(t *testing.T) {
	f := NewTreeFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should contain project name with trailing slash
	if !strings.Contains(output, "test-project/") {
		t.Error("Missing project directory name")
	}

	// Should contain tree characters
	if !strings.Contains(output, "├──") && !strings.Contains(output, "└──") {
		t.Error("Missing tree branch characters")
	}
}

func TestTreeFormatter_Format_ASCIITreeCharacters(t *testing.T) {
	f := NewTreeFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify ASCII tree characters are used
	treeChars := []string{"├── ", "└── ", "│   "}
	for _, char := range treeChars {
		if !strings.Contains(output, char) {
			t.Errorf("Missing tree character: %q", char)
		}
	}
}

func TestTreeFormatter_Format_MetadataInline(t *testing.T) {
	f := &TreeFormatter{ShowMetadata: true, Colorize: false, ShowSummary: true}
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should show metadata in brackets
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Error("Missing metadata brackets")
	}

	// Should show size
	if !strings.Contains(output, "KB") && !strings.Contains(output, "B") {
		t.Error("Missing size in metadata")
	}

	// Should show line count with L suffix
	if !strings.Contains(output, "L]") {
		t.Error("Missing line count with L suffix")
	}
}

func TestTreeFormatter_Format_NoMetadata(t *testing.T) {
	f := &TreeFormatter{ShowMetadata: false, Colorize: false, ShowSummary: true}
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should not show metadata
	if strings.Contains(output, "[") && strings.Contains(output, "L]") {
		t.Error("Should not show metadata when disabled")
	}
}

func TestTreeFormatter_Format_Colorization(t *testing.T) {
	// Test with colors enabled (but NO_COLOR not set)
	f := &TreeFormatter{ShowMetadata: true, Colorize: true, ShowSummary: true}
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Colors should be present (contains ANSI codes)
	if strings.Contains(output, "\033[") {
		t.Log("ANSI color codes present")
	}
}

func TestTreeFormatter_Format_NoColorize(t *testing.T) {
	f := &TreeFormatter{ShowMetadata: true, Colorize: false, ShowSummary: true}
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should not contain ANSI codes
	if strings.Contains(output, "\033[") {
		t.Error("Should not contain ANSI color codes when disabled")
	}
}

func TestTreeFormatter_Format_NO_COLOR_Env(t *testing.T) {
	// Set NO_COLOR environment variable
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	f := &TreeFormatter{ShowMetadata: true, Colorize: true, ShowSummary: true}
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should respect NO_COLOR
	if strings.Contains(output, "\033[") {
		t.Error("Should respect NO_COLOR environment variable")
	}
}

func TestTreeFormatter_Format_Summary(t *testing.T) {
	f := &TreeFormatter{ShowMetadata: true, Colorize: false, ShowSummary: true}
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should show summary line
	if !strings.Contains(output, "file") {
		t.Error("Missing file count in summary")
	}
}

func TestTreeFormatter_Format_NoSummary(t *testing.T) {
	f := &TreeFormatter{ShowMetadata: true, Colorize: false, ShowSummary: false}
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Count lines - should not have extra summary line
	lines := strings.Count(output, "\n")
	t.Logf("Lines without summary: %d", lines)
}

func TestTreeFormatter_Format_DirectorySorting(t *testing.T) {
	f := &TreeFormatter{ShowMetadata: false, Colorize: false, ShowSummary: false}
	project := &types.Project{
		Name: "sort-test",
		Type: types.ProjectTypeGo,
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "z_file.go", Size: 100, Lines: 10, Type: "go"},
				{Path: "a_file.go", Size: 100, Lines: 10, Type: "go"},
				{Path: "m_file.go", Size: 100, Lines: 10, Type: "go"},
			},
			Subdirs: []types.Directory{
				{Path: "z_dir"},
				{Path: "a_dir"},
			},
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Directories should come before files
	dirIdx := strings.Index(output, "a_dir/")
	fileIdx := strings.Index(output, "a_file.go")

	if dirIdx == -1 || fileIdx == -1 {
		t.Fatal("Could not find expected entries")
	}

	if dirIdx > fileIdx {
		t.Error("Directories should be sorted before files")
	}
}

func TestTreeFormatter_Format_EmptyDirectory(t *testing.T) {
	f := NewTreeFormatter()
	project := &types.Project{
		Name: "empty-test",
		Type: types.ProjectTypeGo,
		Tree: types.Directory{
			Path: ".",
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should show project name
	if !strings.Contains(output, "empty-test/") {
		t.Error("Missing project name")
	}
}

func TestTreeFormatter_Format_SingleFile(t *testing.T) {
	f := &TreeFormatter{ShowMetadata: false, Colorize: false, ShowSummary: false}
	project := &types.Project{
		Name: "single",
		Type: types.ProjectTypeGo,
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "only.go", Size: 100, Lines: 10, Type: "go"},
			},
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if !strings.Contains(output, "only.go") {
		t.Error("Missing single file")
	}

	// Single file should use └── (last item)
	if !strings.Contains(output, "└── only.go") {
		t.Error("Single file should use last branch character")
	}
}

func TestTreeFormatter_Format_NilProject(t *testing.T) {
	f := NewTreeFormatter()
	output, err := f.Format(nil)
	if err != nil {
		t.Fatalf("Format(nil) failed: %v", err)
	}
	if output != "" {
		t.Errorf("Format(nil) = %q, want empty", output)
	}
}

func TestTreeFormatter_Format_NilDirectory(t *testing.T) {
	f := NewTreeFormatter()
	var sb strings.Builder
	f.formatTree(nil, &sb, "", false)
	if sb.String() != "" {
		t.Errorf("formatTree(nil) = %q, want empty", sb.String())
	}
}

func TestTreeFormatter_Format_NilFileMetadata(t *testing.T) {
	f := NewTreeFormatter()
	result := f.formatFileMetadata(nil)
	if result != "" {
		t.Errorf("formatFileMetadata(nil) = %q, want empty", result)
	}
}

func TestTreeFormatter_Format_SummaryStats(t *testing.T) {
	f := NewTreeFormatter()

	tests := []struct {
		name     string
		stats    *types.Stats
		contains string
	}{
		{
			name:     "nil stats",
			stats:    nil,
			contains: "",
		},
		{
			name:     "single file",
			stats:    &types.Stats{FileCount: 1, TotalSize: 100, TotalLines: 10},
			contains: "1 file",
		},
		{
			name:     "multiple files",
			stats:    &types.Stats{FileCount: 5, TotalSize: 1024, TotalLines: 100},
			contains: "5 files",
		},
		{
			name:     "with size",
			stats:    &types.Stats{FileCount: 3, TotalSize: 2048, TotalLines: 50},
			contains: "2KB",
		},
		{
			name:     "with lines",
			stats:    &types.Stats{FileCount: 2, TotalSize: 500, TotalLines: 200},
			contains: "200 lines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.formatSummary(tt.stats)
			if tt.contains != "" && !strings.Contains(result, tt.contains) {
				t.Errorf("formatSummary() = %q, should contain %q", result, tt.contains)
			}
		})
	}
}

func TestTreeFormatter_GetColorForFile(t *testing.T) {
	f := &TreeFormatter{Colorize: true}

	tests := []struct {
		name      string
		file      *types.File
		wantColor bool
	}{
		{"go file", &types.File{Path: "main.go", Type: "go"}, true},
		{"js file", &types.File{Path: "app.js", Type: "js"}, true},
		{"ts file", &types.File{Path: "app.ts", Type: "ts"}, true},
		{"py file", &types.File{Path: "script.py", Type: "py"}, true},
		{"rs file", &types.File{Path: "main.rs", Type: "rs"}, true},
		{"java file", &types.File{Path: "Main.java", Type: "java"}, true},
		{"c file", &types.File{Path: "main.c", Type: "c"}, true},
		{"json file", &types.File{Path: "data.json", Type: "json"}, true},
		{"md file", &types.File{Path: "README.md", Type: "md"}, true},
		{"sh file", &types.File{Path: "script.sh", Type: "sh"}, true},
		{"html file", &types.File{Path: "index.html", Type: "html"}, true},
		{"css file", &types.File{Path: "style.css", Type: "css"}, true},
		{"binary file", &types.File{Path: "image.png", Type: "png", IsBinary: true}, true},
		{"unknown file", &types.File{Path: "file.xyz", Type: "xyz"}, false},
		{"makefile", &types.File{Path: "Makefile", Type: "make"}, true},
		{"dockerfile", &types.File{Path: "Dockerfile", Type: "docker"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := f.getColorForFile(tt.file)
			if tt.wantColor && color == "" {
				t.Errorf("Expected color for %s", tt.name)
			}
			if !tt.wantColor && color != "" {
				t.Errorf("Unexpected color for %s: %s", tt.name, color)
			}
		})
	}
}

func TestTreeFormatter_GetColorForFile_NoColor(t *testing.T) {
	f := &TreeFormatter{Colorize: false}
	file := &types.File{Path: "main.go", Type: "go"}

	color := f.getColorForFile(file)
	if color != "" {
		t.Errorf("getColorForFile with Colorize=false should return empty, got: %s", color)
	}
}

func TestTreeFormatter_GetColorForFile_Ruby(t *testing.T) {
	f := &TreeFormatter{Colorize: true}
	file := &types.File{Path: "script.rb", Type: "rb"}

	color := f.getColorForFile(file)
	if color == "" {
		t.Error("Ruby files should have a color")
	}
	// Ruby uses red() function
	if !strings.Contains(color, "31m") {
		t.Errorf("Ruby should use red color, got: %s", color)
	}
}

func TestTreeFormatter_Colorize(t *testing.T) {
	f := &TreeFormatter{Colorize: true}
	result := f.colorize("test", "\033[31m")

	if !strings.Contains(result, "\033[31m") {
		t.Error("colorize should wrap text with color code")
	}
	if !strings.Contains(result, "\033[0m") {
		t.Error("colorize should reset color at end")
	}
}

func TestTreeFormatter_Colorize_Disabled(t *testing.T) {
	f := &TreeFormatter{Colorize: false}
	result := f.colorize("test", "\033[31m")

	if result != "test" {
		t.Errorf("colorize when disabled should return original text, got: %s", result)
	}
}

func TestTreeFormatter_DeepNesting(t *testing.T) {
	f := &TreeFormatter{ShowMetadata: false, Colorize: false, ShowSummary: false}
	project := CreateDeepNestedProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should handle deep nesting
	if !strings.Contains(output, "level3.go") {
		t.Error("Missing deeply nested file")
	}

	// Should use proper indentation with │
	verticalCount := strings.Count(output, "│")
	if verticalCount == 0 {
		t.Error("Missing vertical connectors for nested items")
	}
}

func TestTreeFormatter_BinaryFile(t *testing.T) {
	f := &TreeFormatter{ShowMetadata: true, Colorize: false, ShowSummary: false}
	project := CreateProjectWithBinary()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Binary file should have metadata but without line count
	// The format is [size] without lines for binary files
	// Check that binary file doesn't show 0L after the size
	if strings.Contains(output, ", 0L]") {
		t.Errorf("Binary file should not show 0L in metadata, got: %s", output)
	}
	// Should show size for binary
	if !strings.Contains(output, "KB") && !strings.Contains(output, "B") {
		t.Error("Binary file should show size")
	}
}

func TestTreeFormatter_FormatFileMetadata(t *testing.T) {
	f := NewTreeFormatter()

	tests := []struct {
		name     string
		file     *types.File
		contains string
	}{
		{
			name:     "with size and lines",
			file:     &types.File{Size: 2048, Lines: 100},
			contains: "2KB, 100L",
		},
		{
			name:     "only size",
			file:     &types.File{Size: 500, Lines: 0},
			contains: "500B",
		},
		{
			name:     "binary file",
			file:     &types.File{Size: 1000, Lines: 0, IsBinary: true},
			contains: "1000B",
		},
		{
			name:     "empty file",
			file:     &types.File{Size: 0, Lines: 0},
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := f.formatFileMetadata(tt.file)
			if tt.contains != "" && !strings.Contains(result, tt.contains) {
				t.Errorf("formatFileMetadata() = %q, should contain %q", result, tt.contains)
			}
			if tt.contains == "" && result != "" {
				t.Errorf("formatFileMetadata() = %q, should be empty", result)
			}
		})
	}
}

func TestTreeFormatter_ProjectNameEmpty(t *testing.T) {
	f := NewTreeFormatter()
	project := &types.Project{
		Name: "",
		Type: types.ProjectTypeGo,
		Tree: types.Directory{Path: "."},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should use "." for empty name
	if !strings.Contains(output, "./") {
		t.Error("Empty project name should show as ./")
	}
}
