package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/wesbragagt/gps/pkg/types"
)

func TestFlatFormatter_Format_BasicOutput(t *testing.T) {
	f := NewFlatFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should have header row
	if !strings.Contains(output, "path") {
		t.Error("Missing header row")
	}

	// Should have data rows
	if !strings.Contains(output, "main.go") {
		t.Error("Missing file data")
	}
}

func TestFlatFormatter_Format_HeaderRow(t *testing.T) {
	f := &FlatFormatter{
		ShowHeader: true,
		SortBy:     "path",
		ShowFields: []string{"path", "size", "lines", "type", "importance"},
	}
	project := CreateMinimalProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// First line should be header
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		t.Fatal("No output lines")
	}

	expectedHeader := "path,size,lines,type,importance"
	if lines[0] != expectedHeader {
		t.Errorf("Header = %q, want %q", lines[0], expectedHeader)
	}
}

func TestFlatFormatter_Format_NoHeader(t *testing.T) {
	f := &FlatFormatter{
		ShowHeader: false,
		SortBy:     "path",
		ShowFields: []string{"path", "size"},
	}
	project := &types.Project{
		Name: "test",
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "a.go", Size: 100},
			},
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should not contain header labels
	if strings.Contains(output, "path,size") {
		t.Error("Should not have header row")
	}

	// Should only have data
	if !strings.Contains(output, "a.go") {
		t.Error("Missing file data")
	}
}

func TestFlatFormatter_Format_SortOptions(t *testing.T) {
	project := &types.Project{
		Name: "sort-test",
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "z.go", Size: 100, Lines: 10, Type: "go", Importance: 10},
				{Path: "a.go", Size: 500, Lines: 50, Type: "go", Importance: 90},
				{Path: "m.go", Size: 300, Lines: 30, Type: "go", Importance: 50},
			},
		},
	}

	tests := []struct {
		name       string
		sortBy     string
		descending bool
		expected   string
	}{
		{
			name:       "path ascending",
			sortBy:     "path",
			descending: false,
			expected:   "a.go\nm.go\nz.go",
		},
		{
			name:       "path descending",
			sortBy:     "path",
			descending: true,
			expected:   "z.go\nm.go\na.go",
		},
		{
			name:       "size ascending",
			sortBy:     "size",
			descending: false,
			expected:   "z.go\nm.go\na.go",
		},
		{
			name:       "size descending",
			sortBy:     "size",
			descending: true,
			expected:   "a.go\nm.go\nz.go",
		},
		{
			name:       "lines ascending",
			sortBy:     "lines",
			descending: false,
			expected:   "z.go\nm.go\na.go",
		},
		{
			name:       "importance ascending",
			sortBy:     "importance",
			descending: false,
			expected:   "z.go\nm.go\na.go",
		},
		{
			name:       "importance descending",
			sortBy:     "importance",
			descending: true,
			expected:   "a.go\nm.go\nz.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FlatFormatter{
				ShowHeader:     false,
				SortBy:         tt.sortBy,
				SortDescending: tt.descending,
				ShowFields:     []string{"path"},
			}

			output, err := f.Format(project)
			if err != nil {
				t.Fatalf("Format failed: %v", err)
			}

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Sort by %s (desc=%v):\nGot: %q\nWant order: %s", tt.sortBy, tt.descending, output, tt.expected)
			}
		})
	}
}

func TestFlatFormatter_Format_FieldSelection(t *testing.T) {
	project := &types.Project{
		Name: "field-test",
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "test.go", Size: 1000, Lines: 50, Type: "go", Importance: 80},
			},
		},
	}

	tests := []struct {
		name     string
		fields   []string
		contains string
	}{
		{
			name:     "only path",
			fields:   []string{"path"},
			contains: "test.go",
		},
		{
			name:     "path and size",
			fields:   []string{"path", "size"},
			contains: "test.go,1000",
		},
		{
			name:     "all fields",
			fields:   []string{"path", "size", "lines", "type", "importance"},
			contains: "test.go,1000,50,go,80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FlatFormatter{
				ShowHeader: true,
				SortBy:     "path",
				ShowFields: tt.fields,
			}

			output, err := f.Format(project)
			if err != nil {
				t.Fatalf("Format failed: %v", err)
			}

			if !strings.Contains(output, tt.contains) {
				t.Errorf("Output should contain %q, got: %s", tt.contains, output)
			}
		})
	}
}

func TestFlatFormatter_Format_PathQuoting(t *testing.T) {
	f := &FlatFormatter{
		ShowHeader: false,
		SortBy:     "path",
		ShowFields: []string{"path"},
	}
	project := &types.Project{
		Name: "quote-test",
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "file,with,commas.go", Size: 100},
				{Path: `file"with"quotes.go`, Size: 100},
			},
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Paths with commas should be quoted
	if !strings.Contains(output, `"file,with,commas.go"`) {
		t.Error("Path with commas should be quoted")
	}

	// Quotes should be escaped
	if !strings.Contains(output, `"file""with""quotes.go"`) {
		t.Error("Quotes in path should be escaped")
	}
}

func TestFlatFormatter_Format_ModTime(t *testing.T) {
	modTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	f := &FlatFormatter{
		ShowHeader: true,
		SortBy:     "path",
		ShowFields: []string{"path", "modtime"},
	}
	project := &types.Project{
		Name: "modtime-test",
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "test.go", ModTime: modTime},
			},
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should contain RFC3339 formatted time
	if !strings.Contains(output, "2024-01-15T10:30:00Z") {
		t.Errorf("Output should contain formatted modtime, got: %s", output)
	}
}

func TestFlatFormatter_Format_ZeroModTime(t *testing.T) {
	f := &FlatFormatter{
		ShowHeader: false,
		SortBy:     "path",
		ShowFields: []string{"path", "modtime"},
	}
	project := &types.Project{
		Name: "zero-modtime",
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "test.go", ModTime: time.Time{}},
			},
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Zero time should be empty string
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(lines))
	}

	// Should be "test.go," (empty modtime)
	if !strings.HasSuffix(lines[0], ",") {
		t.Errorf("Zero modtime should be empty, got: %q", lines[0])
	}
}

func TestFlatFormatter_Format_NilProject(t *testing.T) {
	f := NewFlatFormatter()
	output, err := f.Format(nil)
	if err != nil {
		t.Fatalf("Format(nil) failed: %v", err)
	}
	if output != "" {
		t.Errorf("Format(nil) = %q, want empty", output)
	}
}

func TestFlatFormatter_Format_EmptyProject(t *testing.T) {
	f := &FlatFormatter{
		ShowHeader: true,
		SortBy:     "path",
		ShowFields: []string{"path"},
	}
	project := &types.Project{
		Name: "empty",
		Tree: types.Directory{Path: "."},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should only have header
	if !strings.Contains(output, "path") {
		t.Error("Should have header")
	}

	// Count lines - should be 2 (header + empty)
	lineCount := strings.Count(output, "\n")
	if lineCount != 1 {
		t.Errorf("Expected 1 newline (header only), got %d", lineCount)
	}
}

func TestFlatFormatter_Format_NilDirectory(t *testing.T) {
	f := NewFlatFormatter()
	files := f.flattenTree(nil)
	if files != nil {
		t.Errorf("flattenTree(nil) = %v, want nil", files)
	}
}

func TestFlatFormatter_SortByModTime(t *testing.T) {
	project := &types.Project{
		Name: "modtime-sort",
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "old.go", ModTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Path: "new.go", ModTime: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)},
				{Path: "mid.go", ModTime: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
	}

	f := &FlatFormatter{
		ShowHeader: false,
		SortBy:     "modtime",
		ShowFields: []string{"path"},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Oldest first
	expected := "old.go\nmid.go\nnew.go"
	if !strings.Contains(output, expected) {
		t.Errorf("Modtime sort failed:\nGot: %q\nWant order: %s", output, expected)
	}
}

func TestFlatFormatter_SortByType(t *testing.T) {
	project := &types.Project{
		Name: "type-sort",
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "main.go", Type: "go"},
				{Path: "readme.md", Type: "md"},
				{Path: "config.json", Type: "json"},
			},
		},
	}

	f := &FlatFormatter{
		ShowHeader: false,
		SortBy:     "type",
		ShowFields: []string{"path", "type"},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Alphabetical by type: go, json, md
	// Verify order in output
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines, got %d: %v", len(lines), lines)
	}

	// First should be go file (type "go" comes before "json" alphabetically)
	if !strings.HasPrefix(lines[0], "main.go") {
		t.Errorf("First file should be main.go, got: %s", lines[0])
	}
	// Second should be json file
	if !strings.HasPrefix(lines[1], "config.json") {
		t.Errorf("Second file should be config.json, got: %s", lines[1])
	}
	// Third should be md file
	if !strings.HasPrefix(lines[2], "readme.md") {
		t.Errorf("Third file should be readme.md, got: %s", lines[2])
	}
}

func TestFlatFormatter_UnknownField(t *testing.T) {
	f := &FlatFormatter{
		ShowHeader: false,
		SortBy:     "path",
		ShowFields: []string{"path", "unknown"},
	}
	project := &types.Project{
		Name: "unknown-field",
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "test.go"},
			},
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Unknown field should produce empty column
	if !strings.Contains(output, "test.go,") {
		t.Errorf("Unknown field should be empty, got: %s", output)
	}
}

func TestFlatFormatter_NestedDirectories(t *testing.T) {
	f := &FlatFormatter{
		ShowHeader: false,
		SortBy:     "path",
		ShowFields: []string{"path"},
	}
	project := CreateDeepNestedProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should include files from nested directories
	if !strings.Contains(output, "level1.go") {
		t.Error("Missing level1.go")
	}
	if !strings.Contains(output, "level2.go") {
		t.Error("Missing level2.go")
	}
	if !strings.Contains(output, "level3.go") {
		t.Error("Missing level3.go")
	}
}

func TestCsvEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple",
			input:    "simple",
			expected: "simple",
		},
		{
			name:     "with comma",
			input:    "with,comma",
			expected: `"with,comma"`,
		},
		{
			name:     "with quote",
			input:    `with"quote`,
			expected: `"with""quote"`,
		},
		{
			name:     "with newline",
			input:    "with\nnewline",
			expected: "\"with\nnewline\"",
		},
		{
			name:     "with carriage return",
			input:    "with\rcarriage",
			expected: "\"with\rcarriage\"",
		},
		{
			name:     "multiple special chars",
			input:    `test,"data",more`,
			expected: `"test,""data"",more"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := csvEscape(tt.input)
			if result != tt.expected {
				t.Errorf("csvEscape(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFlatFormatter_FormatHeader(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		expected string
	}{
		{
			name:     "single field",
			fields:   []string{"path"},
			expected: "path",
		},
		{
			name:     "multiple fields",
			fields:   []string{"path", "size", "lines"},
			expected: "path,size,lines",
		},
		{
			name:     "all fields",
			fields:   []string{"path", "size", "lines", "type", "importance", "modtime"},
			expected: "path,size,lines,type,importance,modtime",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FlatFormatter{ShowFields: tt.fields}
			result := f.formatHeader()
			if result != tt.expected {
				t.Errorf("formatHeader() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFlatFormatter_FormatRow(t *testing.T) {
	f := &FlatFormatter{
		ShowFields: []string{"path", "size", "lines", "type", "importance"},
	}

	file := flatFile{
		Path:       "test.go",
		Size:       1024,
		Lines:      50,
		Type:       "go",
		Importance: 80,
	}

	result := f.formatRow(file)
	expected := "test.go,1024,50,go,80"
	if result != expected {
		t.Errorf("formatRow() = %q, want %q", result, expected)
	}
}

func TestFlatFormatter_DefaultSettings(t *testing.T) {
	f := NewFlatFormatter()

	if !f.ShowHeader {
		t.Error("Default ShowHeader should be true")
	}
	if f.SortBy != "path" {
		t.Errorf("Default SortBy = %q, want 'path'", f.SortBy)
	}
	if f.SortDescending {
		t.Error("Default SortDescending should be false")
	}
	if len(f.ShowFields) == 0 {
		t.Error("Default ShowFields should not be empty")
	}
}
