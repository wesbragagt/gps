package metadata

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wesbragagt/gps/pkg/types"
)

func TestDetectFileType(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"main.go", "go"},
		{"app.js", "javascript"},
		{"app.jsx", "javascript"},
		{"app.ts", "typescript"},
		{"app.tsx", "typescript"},
		{"script.py", "python"},
		{"Main.java", "java"},
		{"lib.rs", "rust"},
		{"main.c", "c"},
		{"main.cpp", "cpp"},
		{"app.rb", "ruby"},
		{"index.php", "php"},
		{"README.md", "readme"}, // Special filename takes precedence
		{"index.html", "html"},
		{"config.yaml", "yaml"},
		{"data.json", "json"},
		{"config.toml", "toml"},
		{"Dockerfile", "dockerfile"},
		{"Makefile", "makefile"},
		{"style.css", "css"},
		{"app.vue", "vue"},
		{"unknown.xyz", "unknown"},
		{".gitignore", "gitignore"},
		{"LICENSE", "license"},
		{"README", "readme"},
	}

	e := NewExtractor()

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := e.detectFileType(tt.path)
			if result != tt.expected {
				t.Errorf("detectFileType(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestHumanizeSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0B"},
		{100, "100B"},
		{512, "512B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{1048576, "1.0MB"},
		{1572864, "1.5MB"},
		{1073741824, "1.0GB"},
		{1610612736, "1.5GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := HumanizeSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("HumanizeSize(%d) = %q, want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestIsBinaryFile(t *testing.T) {
	// Create temp directory for test files
	tmpDir := t.TempDir()

	e := NewExtractor()

	// Test text file
	textFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(textFile, []byte("Hello, World!\nThis is a text file.\n"), 0644); err != nil {
		t.Fatalf("Failed to create text file: %v", err)
	}

	info, _ := os.Stat(textFile)
	if e.isBinaryFile(textFile, info) {
		t.Error("Text file detected as binary")
	}

	// Test binary file (contains null bytes)
	binaryFile := filepath.Join(tmpDir, "test.bin")
	binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x05}
	if err := os.WriteFile(binaryFile, binaryData, 0644); err != nil {
		t.Fatalf("Failed to create binary file: %v", err)
	}

	info, _ = os.Stat(binaryFile)
	if !e.isBinaryFile(binaryFile, info) {
		t.Error("Binary file not detected as binary")
	}

	// Test known binary extension
	pngFile := filepath.Join(tmpDir, "test.png")
	if err := os.WriteFile(pngFile, []byte("fake png"), 0644); err != nil {
		t.Fatalf("Failed to create png file: %v", err)
	}

	info, _ = os.Stat(pngFile)
	if !e.isBinaryFile(pngFile, info) {
		t.Error("PNG file not detected as binary by extension")
	}
}

func TestCountLines(t *testing.T) {
	tmpDir := t.TempDir()

	e := NewExtractor()

	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{"empty", "", 0},
		{"single line", "Hello", 1},
		{"two lines", "Hello\nWorld", 2},
		{"trailing newline", "Hello\nWorld\n", 2},
		{"multiple empty lines", "\n\n\n", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(tmpDir, tt.name+".txt")
			if err := os.WriteFile(file, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}

			lines, err := e.countLines(file)
			if err != nil {
				t.Fatalf("countLines failed: %v", err)
			}
			if lines != tt.expected {
				t.Errorf("countLines() = %d, want %d", lines, tt.expected)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tmpDir := t.TempDir()

	e := NewExtractor()

	// Create a test Go file
	goFile := filepath.Join(tmpDir, "test.go")
	content := `package main

func main() {
	println("Hello")
}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create Go file: %v", err)
	}

	info, err := os.Stat(goFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	file, err := e.Extract(goFile, info)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if file.Path != goFile {
		t.Errorf("Path = %q, want %q", file.Path, goFile)
	}

	if file.Type != "go" {
		t.Errorf("Type = %q, want %q", file.Type, "go")
	}

	if file.IsBinary {
		t.Error("Go file detected as binary")
	}

	if file.Lines != 5 {
		t.Errorf("Lines = %d, want %d", file.Lines, 5)
	}

	if file.Size == 0 {
		t.Error("Size should not be 0")
	}
}

func TestExtractWithNilInfo(t *testing.T) {
	tmpDir := t.TempDir()

	e := NewExtractor()

	// Create a test file
	file := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Extract with nil info (should stat the file itself)
	result, err := e.Extract(file, nil)
	if err != nil {
		t.Fatalf("Extract with nil info failed: %v", err)
	}

	if result.Size == 0 {
		t.Error("Size should not be 0")
	}
}

func TestExtractNonExistent(t *testing.T) {
	e := NewExtractor()

	_, err := e.Extract("/nonexistent/file.txt", nil)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestExtractSimple(t *testing.T) {
	tmpDir := t.TempDir()

	file := filepath.Join(tmpDir, "simple.py")
	content := "#!/usr/bin/env python3\nprint('hello')\n"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	result, err := ExtractSimple(file)
	if err != nil {
		t.Fatalf("ExtractSimple failed: %v", err)
	}

	if result.Type != "python" {
		t.Errorf("Type = %q, want %q", result.Type, "python")
	}

	if result.IsBinary {
		t.Error("Python file detected as binary")
	}
}

func TestIsKnownBinaryExtension(t *testing.T) {
	tests := []struct {
		ext      string
		expected bool
	}{
		{".png", true},
		{".jpg", true},
		{".exe", true},
		{".zip", true},
		{".pdf", true},
		{".txt", false},
		{".go", false},
		{".md", false},
		{".svg", false}, // SVG is text-based
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			result := isKnownBinaryExtension(tt.ext)
			if result != tt.expected {
				t.Errorf("isKnownBinaryExtension(%q) = %v, want %v", tt.ext, result, tt.expected)
			}
		})
	}
}

// Test types.File struct is properly populated
func TestFileStructFields(t *testing.T) {
	tmpDir := t.TempDir()

	e := NewExtractor()

	file := filepath.Join(tmpDir, "test.js")
	content := "console.log('test');\n"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	result, err := e.Extract(file, nil)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Verify all expected fields are set
	if result.Path == "" {
		t.Error("Path not set")
	}
	if result.Size == 0 {
		t.Error("Size not set")
	}
	if result.Lines == 0 {
		t.Error("Lines not set")
	}
	if result.Type == "" {
		t.Error("Type not set")
	}
	if result.ModTime.IsZero() {
		t.Error("ModTime not set")
	}
	// IsBinary is already false by default, but should be explicitly set

	// Type check (compile-time assertion)
	var _ *types.File = result
}
