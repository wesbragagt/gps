package detector

import (
	"testing"

	"github.com/wesbragagt/gps/pkg/types"
)

func TestDetect_GoProject(t *testing.T) {
	d := New()

	tree := &types.Directory{
		Path: ".",
		Files: []types.File{
			{Path: "go.mod", Type: "go"},
			{Path: "go.sum", Type: "checksum"},
			{Path: "main.go", Type: "go"},
			{Path: "README.md", Type: "markdown"},
		},
		Subdirs: []types.Directory{
			{
				Path: "cmd",
				Files: []types.File{
					{Path: "cmd/myapp/main.go", Type: "go"},
				},
				Subdirs: []types.Directory{
					{
						Path: "cmd/myapp/internal",
						Files: []types.File{
							{Path: "cmd/myapp/internal/handler.go", Type: "go"},
							{Path: "cmd/myapp/internal/handler_test.go", Type: "go"},
						},
					},
				},
			},
			{
				Path: "pkg",
				Files: []types.File{
					{Path: "pkg/utils/utils.go", Type: "go"},
					{Path: "pkg/utils/utils_test.go", Type: "go"},
				},
			},
		},
	}

	project, err := d.Detect("/test/project", tree)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	// Verify project type
	if project.Type != types.ProjectTypeGo {
		t.Errorf("Expected project type %q, got %q", types.ProjectTypeGo, project.Type)
	}

	// Verify entry points detected
	if len(project.KeyFiles.EntryPoints) == 0 {
		t.Error("Expected entry points to be detected")
	}

	// Verify main.go is an entry point
	foundMain := false
	for _, ep := range project.KeyFiles.EntryPoints {
		if ep == "main.go" {
			foundMain = true
			break
		}
	}
	if !foundMain {
		t.Error("Expected main.go to be detected as entry point")
	}

	// Verify config files
	if len(project.KeyFiles.Configs) == 0 {
		t.Error("Expected config files to be detected")
	}

	// Verify test files
	if len(project.KeyFiles.Tests) == 0 {
		t.Error("Expected test files to be detected")
	}

	// Verify docs
	if len(project.KeyFiles.Docs) == 0 {
		t.Error("Expected doc files to be detected")
	}

	// Verify stats (9 files total)
	if project.Stats.FileCount != 9 {
		t.Errorf("Expected 9 files, got %d", project.Stats.FileCount)
	}
}

func TestDetect_NodeProject(t *testing.T) {
	d := New()

	tree := &types.Directory{
		Path: ".",
		Files: []types.File{
			{Path: "package.json", Type: "json"},
			{Path: "package-lock.json", Type: "json"},
			{Path: "index.js", Type: "javascript"},
			{Path: "README.md", Type: "markdown"},
		},
		Subdirs: []types.Directory{
			{
				Path: "src",
				Files: []types.File{
					{Path: "src/app.js", Type: "javascript"},
					{Path: "src/utils.js", Type: "javascript"},
				},
			},
			{
				Path: "__tests__",
				Files: []types.File{
					{Path: "__tests__/app.test.js", Type: "javascript"},
				},
			},
		},
	}

	project, err := d.Detect("/test/project", tree)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if project.Type != types.ProjectTypeNode {
		t.Errorf("Expected project type %q, got %q", types.ProjectTypeNode, project.Type)
	}

	// Verify entry points
	if len(project.KeyFiles.EntryPoints) == 0 {
		t.Error("Expected entry points to be detected for Node project")
	}
}

func TestDetect_PythonProject(t *testing.T) {
	d := New()

	tree := &types.Directory{
		Path: ".",
		Files: []types.File{
			{Path: "requirements.txt", Type: "text"},
			{Path: "setup.py", Type: "python"},
			{Path: "main.py", Type: "python"},
			{Path: "README.md", Type: "markdown"},
		},
		Subdirs: []types.Directory{
			{
				Path: "mypackage",
				Files: []types.File{
					{Path: "mypackage/__init__.py", Type: "python"},
					{Path: "mypackage/core.py", Type: "python"},
				},
			},
		},
	}

	project, err := d.Detect("/test/project", tree)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if project.Type != types.ProjectTypePython {
		t.Errorf("Expected project type %q, got %q", types.ProjectTypePython, project.Type)
	}
}

func TestDetect_RustProject(t *testing.T) {
	d := New()

	tree := &types.Directory{
		Path: ".",
		Files: []types.File{
			{Path: "Cargo.toml", Type: "toml"},
			{Path: "Cargo.lock", Type: "toml"},
		},
		Subdirs: []types.Directory{
			{
				Path: "src",
				Files: []types.File{
					{Path: "src/main.rs", Type: "rust"},
					{Path: "src/lib.rs", Type: "rust"},
				},
			},
		},
	}

	project, err := d.Detect("/test/project", tree)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if project.Type != types.ProjectTypeRust {
		t.Errorf("Expected project type %q, got %q", types.ProjectTypeRust, project.Type)
	}
}

func TestDetect_JavaProject(t *testing.T) {
	d := New()

	tree := &types.Directory{
		Path: ".",
		Files: []types.File{
			{Path: "pom.xml", Type: "xml"},
		},
		Subdirs: []types.Directory{
			{
				Path: "src/main/java/com/example",
				Files: []types.File{
					{Path: "src/main/java/com/example/Main.java", Type: "java"},
					{Path: "src/main/java/com/example/App.java", Type: "java"},
				},
			},
		},
	}

	project, err := d.Detect("/test/project", tree)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if project.Type != types.ProjectTypeJava {
		t.Errorf("Expected project type %q, got %q", types.ProjectTypeJava, project.Type)
	}
}

func TestDetect_MixedProject(t *testing.T) {
	d := New()

	tree := &types.Directory{
		Path: ".",
		Files: []types.File{
			{Path: "go.mod", Type: "go"},
			{Path: "package.json", Type: "json"},
			{Path: "main.go", Type: "go"},
			{Path: "index.js", Type: "javascript"},
		},
	}

	project, err := d.Detect("/test/project", tree)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if project.Type != types.ProjectTypeMixed {
		t.Errorf("Expected project type %q, got %q", types.ProjectTypeMixed, project.Type)
	}
}

func TestDetect_OtherProject(t *testing.T) {
	d := New()

	tree := &types.Directory{
		Path: ".",
		Files: []types.File{
			{Path: "README.md", Type: "markdown"},
			{Path: "notes.txt", Type: "text"},
		},
	}

	project, err := d.Detect("/test/project", tree)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if project.Type != types.ProjectTypeOther {
		t.Errorf("Expected project type %q, got %q", types.ProjectTypeOther, project.Type)
	}
}

func TestIsGeneratedFile(t *testing.T) {
	d := New()

	tests := []struct {
		path      string
		generated bool
	}{
		{"file.pb.go", true},
		{"file_gen.go", true},
		{"file.min.js", true},
		{"node_modules/package/index.js", true},
		{"dist/bundle.js", true},
		{"build/output.o", true},
		{"main.go", false},
		{"index.js", false},
		{"app.py", false},
		{"src/main.rs", false},
	}

	for _, tt := range tests {
		result := d.isGeneratedFile(tt.path)
		if result != tt.generated {
			t.Errorf("isGeneratedFile(%q) = %v, want %v", tt.path, result, tt.generated)
		}
	}
}

func TestScoreFile(t *testing.T) {
	d := New()

	project := &types.Project{
		Type: types.ProjectTypeGo,
		KeyFiles: types.KeyFiles{
			EntryPoints: []string{"main.go"},
			Configs:     []string{"go.mod"},
			Tests:       []string{"handler_test.go"},
			Docs:        []string{"README.md"},
		},
	}

	tests := []struct {
		path        string
		fileType    string
		isGenerated bool
		minScore    int
		maxScore    int
	}{
		{"main.go", "go", false, 100, 100},       // Entry point
		{"go.mod", "go", false, 90, 90},          // Config
		{"handler.go", "go", false, 70, 80},      // Source file
		{"handler_test.go", "go", false, 50, 50}, // Test file
		{"README.md", "markdown", false, 40, 40}, // Doc
		{"file.pb.go", "go", true, 20, 20},       // Generated
		{"unknown.txt", "text", false, 10, 10},   // Other
	}

	for _, tt := range tests {
		score := d.scoreFile(tt.path, tt.fileType, tt.isGenerated, project)
		if score < tt.minScore || score > tt.maxScore {
			t.Errorf("scoreFile(%q) = %d, want between %d and %d", tt.path, score, tt.minScore, tt.maxScore)
		}
	}
}

func TestCalculateStats(t *testing.T) {
	d := New()

	tree := &types.Directory{
		Path: ".",
		Files: []types.File{
			{Path: "main.go", Type: "go", Size: 100, Lines: 10},
			{Path: "README.md", Type: "markdown", Size: 50, Lines: 5},
		},
		Subdirs: []types.Directory{
			{
				Path: "pkg",
				Files: []types.File{
					{Path: "pkg/utils.go", Type: "go", Size: 200, Lines: 20},
				},
			},
		},
	}

	project := &types.Project{
		Type:  types.ProjectTypeGo,
		Stats: types.Stats{ByType: make(map[string]int)},
	}

	d.calculateStats(tree, &project.Stats)

	if project.Stats.FileCount != 3 {
		t.Errorf("Expected file count 3, got %d", project.Stats.FileCount)
	}

	if project.Stats.TotalSize != 350 {
		t.Errorf("Expected total size 350, got %d", project.Stats.TotalSize)
	}

	if project.Stats.TotalLines != 35 {
		t.Errorf("Expected total lines 35, got %d", project.Stats.TotalLines)
	}

	if project.Stats.ByType["go"] != 2 {
		t.Errorf("Expected 2 go files, got %d", project.Stats.ByType["go"])
	}
}
