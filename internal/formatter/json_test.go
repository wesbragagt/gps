package formatter

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/wesbragagt/gps/pkg/types"
)

func TestJsonFormatter_Format_PrettyPrint(t *testing.T) {
	f := NewJsonFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify pretty print has indentation
	if !strings.Contains(output, "  ") {
		t.Error("Pretty print should have indentation")
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}
}

func TestJsonFormatter_Format_Compact(t *testing.T) {
	f := &JsonFormatter{Compact: true}
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Compact output should not have pretty indentation
	if strings.Contains(output, "\n  ") {
		t.Error("Compact output should not have indentation")
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}
}

func TestJsonFormatter_Format_NoPrettyNoCompact(t *testing.T) {
	f := &JsonFormatter{PrettyPrint: false, Compact: false}
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should fall back to compact
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}
}

func TestJsonFormatter_Format_ValidJSON(t *testing.T) {
	tests := []struct {
		name    string
		project *types.Project
	}{
		{"standard project", CreateTestProject()},
		{"minimal project", CreateMinimalProject()},
		{"project with binary", CreateProjectWithBinary()},
		{"large project", CreateLargeProject()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewJsonFormatter()
			output, err := f.Format(tt.project)
			if err != nil {
				t.Fatalf("Format failed: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal([]byte(output), &result); err != nil {
				t.Errorf("Output is not valid JSON: %v\nOutput: %s", err, output)
			}
		})
	}
}

func TestJsonFormatter_Format_AllFields(t *testing.T) {
	f := NewJsonFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify all top-level fields are present
	requiredFields := []string{
		`"name"`,
		`"type"`,
		`"root"`,
		`"stats"`,
		`"key_files"`,
		`"tree"`,
	}

	for _, field := range requiredFields {
		if !strings.Contains(output, field) {
			t.Errorf("Missing field: %s", field)
		}
	}

	// Verify stats fields
	statsFields := []string{
		`"file_count"`,
		`"total_size"`,
		`"total_lines"`,
		`"by_type"`,
	}

	for _, field := range statsFields {
		if !strings.Contains(output, field) {
			t.Errorf("Missing stats field: %s", field)
		}
	}

	// Verify tree structure
	treeFields := []string{
		`"path"`,
		`"files"`,
		`"subdirs"`,
	}

	for _, field := range treeFields {
		if !strings.Contains(output, field) {
			t.Errorf("Missing tree field: %s", field)
		}
	}
}

func TestJsonFormatter_Format_JSONStructure(t *testing.T) {
	f := NewJsonFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Parse and verify structure
	var result jsonOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify project data
	if result.Project.Name != "test-project" {
		t.Errorf("Project name = %s, want test-project", result.Project.Name)
	}
	if result.Project.Type != types.ProjectTypeGo {
		t.Errorf("Project type = %s, want go", result.Project.Type)
	}
	if result.Project.Stats.FileCount != 5 {
		t.Errorf("File count = %d, want 5", result.Project.Stats.FileCount)
	}
}

func TestJsonFormatter_Format_RoundTrip(t *testing.T) {
	f := NewJsonFormatter()
	original := CreateTestProject()

	output, err := f.Format(original)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Unmarshal back
	var result jsonOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify round-trip preserves data
	if result.Project.Name != original.Name {
		t.Errorf("Name not preserved: got %s, want %s", result.Project.Name, original.Name)
	}
	if result.Project.Type != original.Type {
		t.Errorf("Type not preserved: got %s, want %s", result.Project.Type, original.Type)
	}
	if result.Project.Stats.FileCount != original.Stats.FileCount {
		t.Errorf("FileCount not preserved: got %d, want %d", result.Project.Stats.FileCount, original.Stats.FileCount)
	}
}

func TestJsonFormatter_Format_NilProject(t *testing.T) {
	f := NewJsonFormatter()
	output, err := f.Format(nil)
	if err != nil {
		t.Fatalf("Format(nil) failed: %v", err)
	}
	if output != "{}" {
		t.Errorf("Format(nil) = %q, want {}", output)
	}
}

func TestJsonFormatter_Format_EmptyProject(t *testing.T) {
	f := NewJsonFormatter()
	project := &types.Project{}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should still be valid JSON
	var result jsonOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Empty project output is not valid JSON: %v", err)
	}
}

func TestJsonFormatter_SizeComparison(t *testing.T) {
	project := CreateTestProject()

	// JSON pretty
	jsonPretty := NewJsonFormatter()
	jsonPrettyOut, _ := jsonPretty.Format(project)

	// JSON compact
	jsonCompact := &JsonFormatter{Compact: true}
	jsonCompactOut, _ := jsonCompact.Format(project)

	// TOON
	toon := NewToonFormatter()
	toonOut, _ := toon.Format(project)

	t.Logf("Size comparison:")
	t.Logf("  TOON:           %d bytes", len(toonOut))
	t.Logf("  JSON (pretty):  %d bytes", len(jsonPrettyOut))
	t.Logf("  JSON (compact): %d bytes", len(jsonCompactOut))

	// Compact should be smaller than pretty
	if len(jsonCompactOut) >= len(jsonPrettyOut) {
		t.Error("Compact JSON should be smaller than pretty JSON")
	}

	// TOON should generally be more compact than pretty JSON
	// (This may not always be true for very small projects)
}

func TestJsonFormatter_SpecialCharacters(t *testing.T) {
	f := NewJsonFormatter()
	project := &types.Project{
		Name: "test\"with\"quotes",
		Type: types.ProjectTypeGo,
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "file with spaces.go", Size: 100, Lines: 10, Type: "go"},
				{Path: "file\twith\ttabs.go", Size: 100, Lines: 10, Type: "go"},
			},
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should be valid JSON
	var result jsonOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output with special chars is not valid JSON: %v", err)
	}
}

func TestJsonFormatter_UnicodePath(t *testing.T) {
	f := NewJsonFormatter()
	project := &types.Project{
		Name: "unicode-test",
		Type: types.ProjectTypeOther,
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "日本語.go", Size: 100, Lines: 10, Type: "go"},
				{Path: "emoji_🎉.txt", Size: 50, Lines: 5, Type: "txt"},
			},
		},
	}

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should be valid JSON
	var result jsonOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output with unicode is not valid JSON: %v", err)
	}
}

func TestJsonFormatter_NestedDirectories(t *testing.T) {
	f := NewJsonFormatter()
	project := CreateDeepNestedProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	var result jsonOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify nested structure
	if len(result.Project.Tree.Subdirs) == 0 {
		t.Fatal("Expected nested subdirectories")
	}
}

func TestJsonFormatter_ByTypeStats(t *testing.T) {
	f := NewJsonFormatter()
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	var result jsonOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify ByType map
	if result.Project.Stats.ByType == nil {
		t.Error("ByType should not be nil")
	}
	if result.Project.Stats.ByType["go"] != 5 {
		t.Errorf("ByType[go] = %d, want 5", result.Project.Stats.ByType["go"])
	}
}

func TestJsonFormatter_NoPrettyNoCompactFallsBack(t *testing.T) {
	// Test the fallback path when neither PrettyPrint nor Compact is set
	f := &JsonFormatter{PrettyPrint: false, Compact: false}
	project := CreateTestProject()

	output, err := f.Format(project)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Should produce valid JSON
	var result jsonOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output should be valid JSON: %v", err)
	}
}
