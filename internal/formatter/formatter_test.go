package formatter

import (
	"strings"
	"testing"

	"github.com/wesbragagt/gps/pkg/types"
)

// Formatter interface - all formatters must implement this
type testFormatter interface {
	Format(project *types.Project) (string, error)
}

// TestAllFormattersImplementInterface verifies all formatters implement the interface
func TestAllFormattersImplementInterface(t *testing.T) {
	// This test ensures compile-time interface compliance
	var _ testFormatter = NewToonFormatter()
	var _ testFormatter = NewJsonFormatter()
	var _ testFormatter = NewTreeFormatter()
	var _ testFormatter = NewFlatFormatter()
}

func TestAllFormatters_NilProject(t *testing.T) {
	formatters := []struct {
		name string
		fmt  testFormatter
	}{
		{"TOON", NewToonFormatter()},
		{"JSON", NewJsonFormatter()},
		{"Tree", NewTreeFormatter()},
		{"Flat", NewFlatFormatter()},
	}

	for _, tt := range formatters {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.fmt.Format(nil)
			if err != nil {
				t.Errorf("Format(nil) returned error: %v", err)
			}
			// Just verify it doesn't crash - each formatter has different nil handling
			t.Logf("Format(nil) returned: %q", output)
		})
	}
}

func TestAllFormatters_ValidOutput(t *testing.T) {
	project := CreateTestProject()

	formatters := []struct {
		name              string
		fmt               testFormatter
		shouldContainName bool
	}{
		{"TOON", NewToonFormatter(), true},
		{"JSON", NewJsonFormatter(), true},
		{"Tree", NewTreeFormatter(), true},
		{"Flat", NewFlatFormatter(), false}, // Flat only outputs file paths
	}

	for _, tt := range formatters {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.fmt.Format(project)
			if err != nil {
				t.Fatalf("Format failed: %v", err)
			}

			if output == "" {
				t.Error("Output should not be empty")
			}

			// Verify project name is present (except for Flat which only shows files)
			if tt.shouldContainName && !strings.Contains(output, "test-project") {
				t.Error("Output should contain project name")
			}

			// All formatters should contain at least one file
			if !strings.Contains(output, "main.go") && !strings.Contains(output, "go.mod") {
				t.Error("Output should contain file names")
			}
		})
	}
}

func TestAllFormatters_ContainSameLogicalContent(t *testing.T) {
	project := CreateTestProject()

	// Get all outputs
	toonOut, _ := NewToonFormatter().Format(project)
	jsonOut, _ := NewJsonFormatter().Format(project)
	treeOut, _ := NewTreeFormatter().Format(project)
	flatOut, _ := NewFlatFormatter().Format(project)

	// All should contain the main files
	files := []string{"main.go", "go.mod", "util.go", "api.go"}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			if !strings.Contains(toonOut, file) {
				t.Errorf("TOON missing file: %s", file)
			}
			if !strings.Contains(jsonOut, file) {
				t.Errorf("JSON missing file: %s", file)
			}
			if !strings.Contains(treeOut, file) {
				t.Errorf("Tree missing file: %s", file)
			}
			if !strings.Contains(flatOut, file) {
				t.Errorf("Flat missing file: %s", file)
			}
		})
	}
}

func TestAllFormatters_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		project *types.Project
	}{
		{"minimal", CreateMinimalProject()},
		{"empty", CreateEmptyProject()},
		{"binary", CreateProjectWithBinary()},
		{"deep nested", CreateDeepNestedProject()},
		{"large", CreateLargeProject()},
	}

	formatters := []struct {
		name string
		fmt  testFormatter
	}{
		{"TOON", NewToonFormatter()},
		{"JSON", NewJsonFormatter()},
		{"Tree", NewTreeFormatter()},
		{"Flat", NewFlatFormatter()},
	}

	for _, tt := range tests {
		for _, f := range formatters {
			t.Run(tt.name+"_"+f.name, func(t *testing.T) {
				output, err := f.fmt.Format(tt.project)
				if err != nil {
					t.Errorf("Format failed: %v", err)
				}
				// Just verify it doesn't crash
				t.Logf("Output length: %d", len(output))
			})
		}
	}
}

func TestOutputSizeComparison(t *testing.T) {
	project := CreateTestProject()

	toonOut, _ := NewToonFormatter().Format(project)
	jsonPrettyOut, _ := NewJsonFormatter().Format(project)
	jsonCompactOut, _ := (&JsonFormatter{Compact: true}).Format(project)
	treeOut, _ := NewTreeFormatter().Format(project)
	flatOut, _ := NewFlatFormatter().Format(project)

	t.Log("=== Output Size Comparison ===")
	t.Logf("TOON:           %d bytes", len(toonOut))
	t.Logf("JSON (pretty):  %d bytes", len(jsonPrettyOut))
	t.Logf("JSON (compact): %d bytes", len(jsonCompactOut))
	t.Logf("Tree:           %d bytes", len(treeOut))
	t.Logf("Flat:           %d bytes", len(flatOut))

	// TOON should generally be more compact than pretty JSON
	if len(toonOut) > len(jsonPrettyOut) {
		t.Log("Note: TOON output is larger than pretty JSON (may vary by project size)")
	}
}

func TestTokenReductionClaims(t *testing.T) {
	// Compare approximate token efficiency
	project := CreateTestProject()

	toonOut, _ := NewToonFormatter().Format(project)
	jsonOut, _ := (&JsonFormatter{Compact: true}).Format(project)

	// Rough token estimation: ~4 chars per token on average
	toonTokens := len(toonOut) / 4
	jsonTokens := len(jsonOut) / 4

	t.Logf("Estimated tokens (TOON): %d", toonTokens)
	t.Logf("Estimated tokens (JSON): %d", jsonTokens)

	if jsonTokens > 0 {
		reduction := float64(jsonTokens-toonTokens) / float64(jsonTokens) * 100
		t.Logf("Token reduction: %.1f%%", reduction)
	}
}

func TestFormatterConsistency(t *testing.T) {
	// Run formatter multiple times to ensure consistent output
	project := CreateTestProject()

	formatters := []struct {
		name string
		fmt  testFormatter
	}{
		{"TOON", NewToonFormatter()},
		{"JSON", NewJsonFormatter()},
		{"Tree", NewTreeFormatter()},
		{"Flat", NewFlatFormatter()},
	}

	for _, f := range formatters {
		t.Run(f.name, func(t *testing.T) {
			output1, err1 := f.fmt.Format(project)
			output2, err2 := f.fmt.Format(project)

			if err1 != nil || err2 != nil {
				t.Fatalf("Format errors: %v, %v", err1, err2)
			}

			if output1 != output2 {
				t.Error("Formatter output is not consistent across calls")
			}
		})
	}
}

func TestAllFormatters_WithStats(t *testing.T) {
	project := CreateTestProject()

	// All formatters should include stats in some form
	formatters := []struct {
		name          string
		fmt           testFormatter
		shouldContain []string
	}{
		{"TOON", NewToonFormatter(), []string{"files:", "size:", "lines:"}},
		{"JSON", NewJsonFormatter(), []string{"file_count", "total_size", "total_lines"}},
		{"Tree", NewTreeFormatter(), []string{"file", "KB"}}, // Summary line
		{"Flat", NewFlatFormatter(), []string{}},             // Flat doesn't show project stats
	}

	for _, tt := range formatters {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.fmt.Format(project)
			if err != nil {
				t.Fatalf("Format failed: %v", err)
			}

			for _, c := range tt.shouldContain {
				if !strings.Contains(output, c) {
					t.Errorf("Output should contain %q", c)
				}
			}
		})
	}
}

func TestAllFormatters_WithKeyFiles(t *testing.T) {
	project := CreateTestProject()

	formatters := []struct {
		name string
		fmt  testFormatter
	}{
		{"TOON", NewToonFormatter()},
		{"JSON", NewJsonFormatter()},
		{"Tree", NewTreeFormatter()},
		{"Flat", NewFlatFormatter()},
	}

	for _, tt := range formatters {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.fmt.Format(project)
			if err != nil {
				t.Fatalf("Format failed: %v", err)
			}

			// Key files should be present in TOON and JSON
			if tt.name == "TOON" {
				if !strings.Contains(output, "keyfiles") {
					t.Error("TOON should contain keyfiles section")
				}
			}
			if tt.name == "JSON" {
				if !strings.Contains(output, "key_files") {
					t.Error("JSON should contain key_files")
				}
			}
		})
	}
}

func TestAllFormatters_WithDifferentProjectTypes(t *testing.T) {
	projectTypes := []struct {
		name string
		typ  types.ProjectType
	}{
		{"go", types.ProjectTypeGo},
		{"node", types.ProjectTypeNode},
		{"python", types.ProjectTypePython},
		{"rust", types.ProjectTypeRust},
		{"java", types.ProjectTypeJava},
		{"mixed", types.ProjectTypeMixed},
		{"other", types.ProjectTypeOther},
	}

	for _, tt := range projectTypes {
		t.Run(tt.name, func(t *testing.T) {
			project := &types.Project{
				Name: "type-test",
				Type: tt.typ,
				Tree: types.Directory{Path: "."},
			}

			// All formatters should handle all project types
			formatters := []testFormatter{
				NewToonFormatter(),
				NewJsonFormatter(),
				NewTreeFormatter(),
				NewFlatFormatter(),
			}

			for _, f := range formatters {
				output, err := f.Format(project)
				if err != nil {
					t.Errorf("Format failed for type %s: %v", tt.name, err)
				}
				if output == "" && tt.typ != "" {
					t.Errorf("Empty output for type %s", tt.name)
				}
			}
		})
	}
}

func BenchmarkFormatters(b *testing.B) {
	project := CreateTestProject()

	benchmarks := []struct {
		name string
		fmt  testFormatter
	}{
		{"TOON", NewToonFormatter()},
		{"JSON", NewJsonFormatter()},
		{"JSON_Compact", &JsonFormatter{Compact: true}},
		{"Tree", NewTreeFormatter()},
		{"Flat", NewFlatFormatter()},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bm.fmt.Format(project)
			}
		})
	}
}
