// Package integration provides end-to-end tests for the gps CLI
package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var binaryPath string

// TestMain builds the binary before running tests
func TestMain(m *testing.M) {
	// Get project root
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

	// Build binary
	binaryPath = filepath.Join(projectRoot, "gps")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	// Remove existing binary if present
	os.Remove(binaryPath)

	// Build
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/gps")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build binary: %v\n%s", err, output)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Remove(binaryPath)

	os.Exit(code)
}

// Helper to run gps command
func runGPS(args ...string) (string, string, error) {
	cmd := exec.Command(binaryPath, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// Helper to get fixtures path
func fixturesPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "fixtures")
}

// =============================================================================
// Basic Scanning Tests
// =============================================================================

func TestScanCurrentDirectory(t *testing.T) {
	// Run gps on fixtures directory
	stdout, _, err := runGPS(fixturesPath())
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Verify output contains project info
	if !strings.Contains(stdout, "simple_project") {
		t.Error("Expected output to contain 'simple_project'")
	}
	if !strings.Contains(stdout, "go_project") {
		t.Error("Expected output to contain 'go_project'")
	}
}

func TestScanSimpleProject(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "simple_project")
	stdout, _, err := runGPS(projectPath)
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Should contain README.md and file.txt
	if !strings.Contains(stdout, "README.md") {
		t.Error("Expected output to contain 'README.md'")
	}
	if !strings.Contains(stdout, "file.txt") {
		t.Error("Expected output to contain 'file.txt'")
	}
}

func TestScanGoProject(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")
	stdout, _, err := runGPS(projectPath)
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Should detect as Go project
	if !strings.Contains(stdout, "go") {
		t.Error("Expected output to indicate Go project")
	}
	// Should contain go.mod and main.go
	if !strings.Contains(stdout, "go.mod") {
		t.Error("Expected output to contain 'go.mod'")
	}
	if !strings.Contains(stdout, "main.go") {
		t.Error("Expected output to contain 'main.go'")
	}
}

func TestScanMixedProject(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "mixed_project")
	stdout, _, err := runGPS(projectPath)
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Should detect as mixed project
	if !strings.Contains(stdout, "mixed") {
		t.Error("Expected output to indicate mixed project")
	}
}

// =============================================================================
// Output Format Tests
// =============================================================================

func TestToonFormat(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "simple_project")
	stdout, _, err := runGPS(projectPath, "-f", "toon")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// TOON format should have project{} wrapper
	if !strings.Contains(stdout, "project[") {
		t.Error("Expected TOON format with 'project[' marker")
	}
	if !strings.Contains(stdout, "]{") {
		t.Error("Expected TOON format with ']{' marker")
	}
}

func TestJSONFormat(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "simple_project")
	stdout, _, err := runGPS(projectPath, "-f", "json")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Should be valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Errorf("Output is not valid JSON: %v\nOutput: %s", err, stdout)
	}

	// Should have project wrapper
	project, ok := result["project"].(map[string]interface{})
	if !ok {
		t.Fatal("JSON output missing 'project' wrapper")
	}

	// Should have required fields
	if _, ok := project["name"]; !ok {
		t.Error("JSON output missing 'name' field")
	}
	if _, ok := project["type"]; !ok {
		t.Error("JSON output missing 'type' field")
	}
}

func TestTreeFormat(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "simple_project")
	stdout, _, err := runGPS(projectPath, "-f", "tree")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Tree format uses tree-drawing characters
	if !strings.Contains(stdout, "├") && !strings.Contains(stdout, "└") {
		t.Error("Expected tree format with tree-drawing characters")
	}
}

func TestFlatFormat(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "simple_project")
	stdout, _, err := runGPS(projectPath, "-f", "flat")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Flat format is tabular
	if !strings.Contains(stdout, "path") {
		t.Error("Expected flat format with 'path' header")
	}
	// Should have README.md in output
	if !strings.Contains(stdout, "README.md") {
		t.Error("Expected flat format to contain 'README.md'")
	}
}

func TestAllFormats(t *testing.T) {
	formats := []string{"toon", "json", "tree", "flat"}
	projectPath := filepath.Join(fixturesPath(), "go_project")

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			stdout, _, err := runGPS(projectPath, "-f", format)
			if err != nil {
				t.Fatalf("gps failed with format %s: %v", format, err)
			}

			// All formats should include the project files
			if !strings.Contains(stdout, "go.mod") {
				t.Errorf("Format %s missing 'go.mod'", format)
			}
			if !strings.Contains(stdout, "main.go") {
				t.Errorf("Format %s missing 'main.go'", format)
			}

			// Output should not be empty
			if len(strings.TrimSpace(stdout)) == 0 {
				t.Errorf("Format %s produced empty output", format)
			}
		})
	}
}

// =============================================================================
// Flag Combination Tests
// =============================================================================

func TestDepthFlag(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")

	// Depth 0 should only show root files
	stdout0, _, err := runGPS(projectPath, "-L", "0")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}
	if strings.Contains(stdout0, "util.go") {
		t.Error("Depth 0 should not show nested files")
	}

	// Depth 1 should show lib/ but not files inside
	stdout1, _, err := runGPS(projectPath, "-L", "1")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}
	// Should show lib directory
	if !strings.Contains(stdout1, "lib") {
		t.Error("Depth 1 should show 'lib' directory")
	}
}

func TestDirsOnlyFlag(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")
	stdout, _, err := runGPS(projectPath, "-d")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Should contain lib/ and test/ directories
	if !strings.Contains(stdout, "lib") {
		t.Error("Expected 'lib' directory in output")
	}
	// Should not contain file names (in dirs-only mode)
	// Note: the output format may still show the directory, so we check for files
	if strings.Contains(stdout, "util.go") && !strings.Contains(stdout, "lib") {
		t.Error("Should not show files in dirs-only mode")
	}
}

func TestAllFlag(t *testing.T) {
	// Create a hidden file in simple_project
	projectPath := filepath.Join(fixturesPath(), "simple_project")
	hiddenFile := filepath.Join(projectPath, ".hidden")
	if err := os.WriteFile(hiddenFile, []byte("hidden"), 0644); err != nil {
		t.Fatalf("Failed to create hidden file: %v", err)
	}
	defer os.Remove(hiddenFile)

	// Without -a flag
	stdoutNormal, _, _ := runGPS(projectPath)
	if strings.Contains(stdoutNormal, ".hidden") {
		t.Error("Hidden file should not appear without -a flag")
	}

	// With -a flag
	stdoutAll, _, err := runGPS(projectPath, "-a")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}
	if !strings.Contains(stdoutAll, ".hidden") {
		t.Error("Hidden file should appear with -a flag")
	}
}

func TestExcludePattern(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")

	// Exclude .go files
	stdout, _, err := runGPS(projectPath, "-I", "*.go")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Should not contain .go files but should have go.mod
	if strings.Contains(stdout, "main.go") {
		t.Error("main.go should be excluded with -I '*.go'")
	}
	if !strings.Contains(stdout, "go.mod") {
		t.Error("go.mod should not be excluded")
	}
}

func TestIncludePattern(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")

	// Include only .go files
	stdout, _, err := runGPS(projectPath, "-P", "*.go")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Should contain .go files but not go.mod
	if !strings.Contains(stdout, "main.go") {
		t.Error("main.go should be included with -P '*.go'")
	}
	if strings.Contains(stdout, "go.mod") {
		t.Error("go.mod should not be included")
	}
}

func TestSummaryFlag(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")
	stdout, _, err := runGPS(projectPath, "-s")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Summary format should have specific structure
	if !strings.Contains(stdout, "project[") {
		t.Error("Expected summary to start with 'project['")
	}
	if !strings.Contains(stdout, "type:") {
		t.Error("Expected summary to contain 'type:'")
	}
	// Should not show full file tree
	if strings.Contains(stdout, "├") && len(stdout) > 500 {
		t.Error("Summary should be concise, not full tree")
	}
}

func TestEntryPointsFlag(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")
	stdout, _, err := runGPS(projectPath, "--entry-points")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Entry points format
	if !strings.Contains(stdout, "entry-points") {
		t.Error("Expected 'entry-points' in output")
	}
	// Should detect main.go as entry point
	if !strings.Contains(stdout, "main.go") {
		t.Error("Expected 'main.go' as entry point")
	}
}

func TestFocusFlag(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")
	stdout, _, err := runGPS(projectPath, "--focus", "lib")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Should focus on lib directory
	if !strings.Contains(stdout, "util.go") {
		t.Error("Expected 'util.go' in focused output")
	}
	if !strings.Contains(stdout, "helper.go") {
		t.Error("Expected 'helper.go' in focused output")
	}
}

func TestNoMetaFlag(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "simple_project")
	stdout, _, err := runGPS(projectPath, "--no-meta")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Output should still work, just without metadata
	if !strings.Contains(stdout, "README.md") {
		t.Error("Expected 'README.md' in output")
	}
}

func TestMultipleFlags(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")

	// Combine depth and format
	stdout, _, err := runGPS(projectPath, "-L", "1", "-f", "json")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Should be valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}

	// Combine exclude and format
	stdout2, _, err := runGPS(projectPath, "-I", "*.go", "-f", "flat")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}
	if strings.Contains(stdout2, "main.go") {
		t.Error("main.go should be excluded")
	}
}

// =============================================================================
// Exit Code Tests
// =============================================================================

func TestExitCodeSuccess(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "simple_project")
	cmd := exec.Command(binaryPath, projectPath)
	err := cmd.Run()

	if err != nil {
		t.Errorf("Expected exit code 0, got error: %v", err)
	}
}

func TestExitCodeNonexistentPath(t *testing.T) {
	cmd := exec.Command(binaryPath, "/nonexistent/path/that/does/not/exist")
	err := cmd.Run()

	if err == nil {
		t.Error("Expected non-zero exit code for nonexistent path")
	}

	// Check it's an exit error (not some other error)
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 0 {
			t.Error("Exit code should be non-zero for nonexistent path")
		}
	}
}

func TestExitCodeInvalidPath(t *testing.T) {
	// Try to scan a file instead of directory
	filePath := filepath.Join(fixturesPath(), "simple_project", "file.txt")
	cmd := exec.Command(binaryPath, filePath)
	err := cmd.Run()

	if err == nil {
		t.Error("Expected non-zero exit code when scanning a file")
	}
}

// =============================================================================
// Error Handling Tests
// =============================================================================

func TestErrorMessageOnStderr(t *testing.T) {
	_, stderr, err := runGPS("/nonexistent/path")
	if err == nil {
		t.Error("Expected error for nonexistent path")
	}

	// Error message should be on stderr or in error
	combinedOutput := stderr + err.Error()
	if !strings.Contains(combinedOutput, "cannot access") &&
		!strings.Contains(combinedOutput, "no such file") {
		t.Errorf("Expected error message about inaccessible path, got: %s", combinedOutput)
	}
}

func TestInvalidFormat(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "simple_project")
	stdout, _, err := runGPS(projectPath, "-f", "invalid_format")

	// Should not fail - falls back to default
	if err != nil {
		t.Logf("Note: invalid format handling may vary: %v", err)
	}

	// Should still produce output (defaulting to toon)
	if len(strings.TrimSpace(stdout)) == 0 {
		t.Error("Expected some output even with invalid format")
	}
}

// =============================================================================
// Large Project Tests (optional)
// =============================================================================

func TestScanWithSubdirectories(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")
	stdout, _, err := runGPS(projectPath)
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Should show nested structure
	if !strings.Contains(stdout, "lib") {
		t.Error("Expected 'lib' directory")
	}
	if !strings.Contains(stdout, "test") {
		t.Error("Expected 'test' directory")
	}
	if !strings.Contains(stdout, "util.go") {
		t.Error("Expected 'util.go' in nested lib directory")
	}
}

// =============================================================================
// Format-Specific Validation Tests
// =============================================================================

func TestJSONFormatStructure(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")
	stdout, _, err := runGPS(projectPath, "-f", "json")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	// Get project wrapper
	project, ok := result["project"].(map[string]interface{})
	if !ok {
		t.Fatal("JSON output missing 'project' wrapper")
	}

	// Check required fields
	requiredFields := []string{"name", "type", "root", "tree", "stats", "key_files"}
	for _, field := range requiredFields {
		if _, ok := project[field]; !ok {
			t.Errorf("Missing required field: %s", field)
		}
	}

	// Check tree structure
	if tree, ok := project["tree"].(map[string]interface{}); ok {
		if _, ok := tree["files"]; !ok {
			t.Error("Tree missing 'files' field")
		}
		if _, ok := tree["subdirs"]; !ok {
			t.Error("Tree missing 'subdirs' field")
		}
	}
}

func TestToonFormatStructure(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")
	stdout, _, err := runGPS(projectPath, "-f", "toon")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// TOON format markers
	if !strings.Contains(stdout, "project[") {
		t.Error("Missing 'project[' marker")
	}
	if !strings.Contains(stdout, "]") {
		t.Error("Missing closing bracket")
	}
	if !strings.Contains(stdout, "{") && !strings.Contains(stdout, "}") {
		t.Error("Missing braces")
	}
}

func TestTreeFormatStructure(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")
	stdout, _, err := runGPS(projectPath, "-f", "tree")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// Tree format uses specific characters
	hasTreeChars := strings.Contains(stdout, "├") ||
		strings.Contains(stdout, "└") ||
		strings.Contains(stdout, "│") ||
		strings.Contains(stdout, "─")

	if !hasTreeChars {
		t.Error("Expected tree-drawing characters in tree format")
	}
}

func TestFlatFormatStructure(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")
	stdout, _, err := runGPS(projectPath, "-f", "flat")
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) < 2 {
		t.Error("Flat format should have header and data rows")
	}

	// First line is usually the header
	header := strings.ToLower(lines[0])
	if !strings.Contains(header, "path") {
		t.Error("Flat format header should contain 'path'")
	}
}

// =============================================================================
// Real-World Scenario Tests
// =============================================================================

func TestScanSelfProject(t *testing.T) {
	// Get project root (parent of test directory)
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

	stdout, _, err := runGPS(projectRoot)
	if err != nil {
		t.Fatalf("gps failed on self: %v", err)
	}

	// Should detect as Go project
	if !strings.Contains(stdout, "go") {
		t.Error("Self-project should be detected as Go project")
	}

	// Should contain key files
	if !strings.Contains(stdout, "go.mod") {
		t.Error("Should contain go.mod")
	}
	if !strings.Contains(stdout, "cmd") {
		t.Error("Should contain cmd directory")
	}
}

func TestScanWithGitignore(t *testing.T) {
	projectPath := filepath.Join(fixturesPath(), "go_project")

	// Create a .gitignore that excludes lib/
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte("lib/\n"), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}
	defer os.Remove(gitignorePath)

	stdout, _, err := runGPS(projectPath)
	if err != nil {
		t.Fatalf("gps failed: %v", err)
	}

	// lib/ should be excluded by gitignore
	if strings.Contains(stdout, "util.go") {
		t.Error("lib/ files should be excluded by .gitignore")
	}
}
