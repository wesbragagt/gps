// Package detector provides project type detection and intelligent traversal.
//
// The detector analyzes a project structure to determine:
//   - Project type (Go, Node.js, Python, Rust, Java, Mixed, Other)
//   - Entry points (main.go, index.js, etc.)
//   - Configuration files (go.mod, package.json, etc.)
//   - Test files
//   - Documentation files
//   - File importance scores
//
// Example usage:
//
//	det := detector.New()
//	project, err := det.Detect("/path/to/project", tree)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Project type: %s\n", project.Type)
//	fmt.Printf("Entry points: %v\n", project.KeyFiles.EntryPoints)
package detector

import (
	"path/filepath"
	"strings"

	"github.com/wesbragagt/gps/pkg/types"
)

// Detector analyzes project structure and determines project metadata.
type Detector struct {
	// MaxFileSizeForGeneratedCheck limits bytes read for generated file detection.
	MaxFileSizeForGeneratedCheck int64
}

// New creates a new Detector with default settings.
func New() *Detector {
	return &Detector{
		MaxFileSizeForGeneratedCheck: 1024,
	}
}

// Detect analyzes a project structure and returns a populated Project.
func (d *Detector) Detect(root string, tree *types.Directory) (*types.Project, error) {
	project := &types.Project{
		Name:     filepath.Base(root),
		Root:     root,
		Type:     types.ProjectTypeOther,
		Stats:    types.Stats{ByType: make(map[string]int)},
		KeyFiles: types.KeyFiles{},
		Tree:     *tree,
	}

	// Detect project type
	detectedTypes := d.detectProjectTypes(tree)
	project.Type = d.resolveProjectType(detectedTypes)

	// Collect all files for analysis
	files := d.collectFiles(tree, "")

	// Detect key files
	d.detectEntryPoints(files, project)
	d.detectConfigFiles(files, project)
	d.detectTestFiles(files, project)
	d.detectDocFiles(files, project)

	// Score importance and detect generated files
	d.scoreAndMarkFiles(tree, project)

	// Calculate statistics
	d.calculateStats(tree, &project.Stats)

	return project, nil
}

// projectTypeRules defines characteristic files for each project type.
type projectTypeRule struct {
	configFiles   []string
	sourceFiles   []string
	sourceExt     []string
	directoryName []string
}

var projectTypeRules = map[types.ProjectType]projectTypeRule{
	types.ProjectTypeGo: {
		configFiles:   []string{"go.mod", "go.sum", "go.work"},
		sourceFiles:   []string{},
		sourceExt:     []string{".go"},
		directoryName: []string{},
	},
	types.ProjectTypeNode: {
		configFiles:   []string{"package.json", "package-lock.json", "yarn.lock", ".npmrc", ".yarnrc"},
		sourceFiles:   []string{},
		sourceExt:     []string{".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs"},
		directoryName: []string{"node_modules"},
	},
	types.ProjectTypePython: {
		configFiles:   []string{"requirements.txt", "setup.py", "pyproject.toml", "Pipfile", "Pipfile.lock", "setup.cfg"},
		sourceFiles:   []string{},
		sourceExt:     []string{".py", ".pyw"},
		directoryName: []string{"__pycache__", ".venv", "venv"},
	},
	types.ProjectTypeRust: {
		configFiles:   []string{"Cargo.toml", "Cargo.lock"},
		sourceFiles:   []string{},
		sourceExt:     []string{".rs"},
		directoryName: []string{},
	},
	types.ProjectTypeJava: {
		configFiles:   []string{"pom.xml", "build.gradle", "build.gradle.kts", "settings.gradle", "settings.gradle.kts"},
		sourceFiles:   []string{},
		sourceExt:     []string{".java", ".kt", ".kts"},
		directoryName: []string{},
	},
}

// detectProjectTypes returns all detected project types.
func (d *Detector) detectProjectTypes(tree *types.Directory) []types.ProjectType {
	var detected []types.ProjectType

	for pt, rule := range projectTypeRules {
		if d.matchesProjectType(tree, rule) {
			detected = append(detected, pt)
		}
	}

	return detected
}

// matchesProjectType checks if a directory tree matches a project type rule.
func (d *Detector) matchesProjectType(tree *types.Directory, rule projectTypeRule) bool {
	// Check for config files
	for _, cf := range rule.configFiles {
		if d.hasFile(tree, cf) {
			return true
		}
	}

	// Check for source file extensions
	for _, ext := range rule.sourceExt {
		if d.hasFilesWithExt(tree, ext) {
			return true
		}
	}

	// Check for characteristic directories
	for _, dir := range rule.directoryName {
		if d.hasDirectory(tree, dir) {
			return true
		}
	}

	return false
}

// hasFile checks if a file exists in the tree.
func (d *Detector) hasFile(tree *types.Directory, name string) bool {
	for _, f := range tree.Files {
		if f.Path == name || filepath.Base(f.Path) == name {
			return true
		}
	}
	for _, subdir := range tree.Subdirs {
		if d.hasFile(&subdir, name) {
			return true
		}
	}
	return false
}

// hasFilesWithExt checks if files with extension exist in the tree.
func (d *Detector) hasFilesWithExt(tree *types.Directory, ext string) bool {
	for _, f := range tree.Files {
		if strings.HasSuffix(f.Path, ext) {
			return true
		}
	}
	for _, subdir := range tree.Subdirs {
		if d.hasFilesWithExt(&subdir, ext) {
			return true
		}
	}
	return false
}

// hasDirectory checks if a directory exists in the tree.
func (d *Detector) hasDirectory(tree *types.Directory, name string) bool {
	for _, subdir := range tree.Subdirs {
		if filepath.Base(subdir.Path) == name {
			return true
		}
		if d.hasDirectory(&subdir, name) {
			return true
		}
	}
	return false
}

// resolveProjectType determines the final project type from detected types.
func (d *Detector) resolveProjectType(detected []types.ProjectType) types.ProjectType {
	if len(detected) == 0 {
		return types.ProjectTypeOther
	}
	if len(detected) == 1 {
		return detected[0]
	}
	return types.ProjectTypeMixed
}

// fileInfo holds collected file information for analysis.
type fileInfo struct {
	path string
	typ  string
}

// collectFiles gathers all files from the tree.
func (d *Detector) collectFiles(tree *types.Directory, prefix string) []fileInfo {
	var files []fileInfo

	for _, f := range tree.Files {
		files = append(files, fileInfo{path: f.Path, typ: f.Type})
	}

	for _, subdir := range tree.Subdirs {
		files = append(files, d.collectFiles(&subdir, subdir.Path)...)
	}

	return files
}

// entryPointPatterns defines entry point file patterns per project type.
var entryPointPatterns = map[types.ProjectType][]struct {
	pattern    string
	exactMatch bool
	pathSuffix string
}{
	types.ProjectTypeGo: {
		{pathSuffix: "main.go"},
		{pathSuffix: "/cmd/"},
	},
	types.ProjectTypeNode: {
		{pattern: "index.js", exactMatch: true},
		{pattern: "index.ts", exactMatch: true},
		{pattern: "server.js", exactMatch: true},
		{pattern: "server.ts", exactMatch: true},
		{pattern: "app.js", exactMatch: true},
		{pattern: "app.ts", exactMatch: true},
		{pattern: "main.js", exactMatch: true},
		{pattern: "main.ts", exactMatch: true},
		{pathSuffix: "src/index.js"},
		{pathSuffix: "src/index.ts"},
	},
	types.ProjectTypePython: {
		{pattern: "__main__.py", exactMatch: true},
		{pattern: "main.py", exactMatch: true},
		{pattern: "app.py", exactMatch: true},
		{pattern: "run.py", exactMatch: true},
		{pattern: "wsgi.py", exactMatch: true},
		{pattern: "asgi.py", exactMatch: true},
	},
	types.ProjectTypeRust: {
		{pathSuffix: "src/main.rs"},
	},
	types.ProjectTypeJava: {
		{pattern: "Main.java", exactMatch: false},
		{pattern: "Application.java", exactMatch: false},
	},
}

// detectEntryPoints finds entry point files.
func (d *Detector) detectEntryPoints(files []fileInfo, project *types.Project) {
	patterns, ok := entryPointPatterns[project.Type]
	if !ok {
		return
	}

	for _, f := range files {
		for _, p := range patterns {
			if d.matchesEntryPattern(f.path, p) {
				project.KeyFiles.EntryPoints = append(project.KeyFiles.EntryPoints, f.path)
				break
			}
		}
	}
}

// matchesEntryPattern checks if a file path matches an entry point pattern.
func (d *Detector) matchesEntryPattern(path string, p struct {
	pattern    string
	exactMatch bool
	pathSuffix string
}) bool {
	if p.pathSuffix != "" {
		return strings.HasSuffix(path, p.pathSuffix) || strings.Contains(path, p.pathSuffix)
	}
	if p.exactMatch {
		return filepath.Base(path) == p.pattern
	}
	return strings.Contains(filepath.Base(path), p.pattern)
}

// configPatterns defines config file patterns per project type.
var configPatterns = map[types.ProjectType][]string{
	types.ProjectTypeGo:     {"go.mod", "go.sum", "go.work", "go.work.sum"},
	types.ProjectTypeNode:   {"package.json", "package-lock.json", "yarn.lock", ".npmrc", ".yarnrc", ".yarnrc.yml"},
	types.ProjectTypePython: {"requirements.txt", "setup.py", "pyproject.toml", "Pipfile", "Pipfile.lock", "setup.cfg", "tox.ini"},
	types.ProjectTypeRust:   {"Cargo.toml", "Cargo.lock", ".cargo/config.toml"},
	types.ProjectTypeJava:   {"pom.xml", "build.gradle", "build.gradle.kts", "settings.gradle", "settings.gradle.kts", "gradle.properties"},
}

// genericConfigPatterns are config files applicable to any project.
var genericConfigPatterns = []string{
	".env", ".env.local", ".env.development", ".env.production",
	"docker-compose.yml", "docker-compose.yaml", "docker-compose.override.yml",
	"Dockerfile", "Dockerfile.prod", "Dockerfile.dev",
	"Makefile", "makefile", "Makefile.PL",
	".gitignore", ".gitattributes", ".gitmodules",
	".editorconfig", ".prettierrc", ".eslintrc",
	"renovate.json", "dependabot.yml",
	".github/workflows",
}

// detectConfigFiles finds configuration files.
func (d *Detector) detectConfigFiles(files []fileInfo, project *types.Project) {
	// Project-specific configs
	if patterns, ok := configPatterns[project.Type]; ok {
		for _, f := range files {
			for _, p := range patterns {
				if filepath.Base(f.path) == p {
					project.KeyFiles.Configs = append(project.KeyFiles.Configs, f.path)
					break
				}
			}
		}
	}

	// Generic configs
	for _, f := range files {
		for _, p := range genericConfigPatterns {
			if filepath.Base(f.path) == p || strings.HasPrefix(f.path, p) {
				project.KeyFiles.Configs = append(project.KeyFiles.Configs, f.path)
				break
			}
		}
	}
}

// testPatterns defines test file patterns.
var testPatterns = []struct {
	suffix  string
	prefix  string
	exact   string
	dirName string
}{
	{suffix: "_test.go"},
	{suffix: ".test.js"},
	{suffix: ".test.ts"},
	{suffix: ".spec.js"},
	{suffix: ".spec.ts"},
	{prefix: "test_"},
	{prefix: "tests/"},
	{prefix: "__tests__/"},
	{prefix: "spec/"},
	{prefix: "specs/"},
	{exact: "test"},
	{exact: "tests"},
	{dirName: "test"},
	{dirName: "tests"},
	{dirName: "__tests__"},
	{dirName: "spec"},
	{dirName: "specs"},
}

// detectTestFiles finds test files and directories.
func (d *Detector) detectTestFiles(files []fileInfo, project *types.Project) {
	for _, f := range files {
		for _, p := range testPatterns {
			if p.suffix != "" && strings.HasSuffix(f.path, p.suffix) {
				project.KeyFiles.Tests = append(project.KeyFiles.Tests, f.path)
				break
			}
			if p.prefix != "" && (strings.HasPrefix(filepath.Base(f.path), p.prefix) || strings.Contains(f.path, p.prefix)) {
				project.KeyFiles.Tests = append(project.KeyFiles.Tests, f.path)
				break
			}
			if p.exact != "" && filepath.Base(f.path) == p.exact {
				project.KeyFiles.Tests = append(project.KeyFiles.Tests, f.path)
				break
			}
		}
	}
}

// docPatterns defines documentation file patterns.
var docPatterns = []struct {
	prefix    string
	exact     string
	dirPrefix string
}{
	{prefix: "README"},
	{prefix: "readme"},
	{exact: "README.md"},
	{exact: "README.txt"},
	{exact: "CHANGELOG.md"},
	{exact: "CHANGELOG.txt"},
	{prefix: "CHANGELOG"},
	{exact: "AUTHORS"},
	{exact: "AUTHORS.md"},
	{exact: "CONTRIBUTING.md"},
	{exact: "CONTRIBUTORS"},
	{exact: "LICENSE"},
	{exact: "LICENSE.md"},
	{exact: "LICENSE.txt"},
	{exact: "COPYING"},
	{dirPrefix: "docs/"},
	{dirPrefix: "doc/"},
}

// detectDocFiles finds documentation files.
func (d *Detector) detectDocFiles(files []fileInfo, project *types.Project) {
	for _, f := range files {
		for _, p := range docPatterns {
			base := filepath.Base(f.path)
			if p.prefix != "" && strings.HasPrefix(base, p.prefix) {
				project.KeyFiles.Docs = append(project.KeyFiles.Docs, f.path)
				break
			}
			if p.exact != "" && base == p.exact {
				project.KeyFiles.Docs = append(project.KeyFiles.Docs, f.path)
				break
			}
			if p.dirPrefix != "" && strings.HasPrefix(f.path, p.dirPrefix) {
				project.KeyFiles.Docs = append(project.KeyFiles.Docs, f.path)
				break
			}
		}
	}
}

// generatedPatterns defines patterns for generated files.
var generatedFilePatterns = []struct {
	suffix      string
	prefix      string
	contains    string
	pathContain string
}{
	{suffix: "_gen.go"},
	{suffix: ".pb.go"},
	{suffix: ".pb.gw.go"},
	{suffix: ".min.js"},
	{suffix: ".min.css"},
	{suffix: ".generated.go"},
	{prefix: "generated_"},
	{contains: "zz_generated"},
	{pathContain: "node_modules/"},
	{pathContain: "dist/"},
	{pathContain: "build/"},
	{pathContain: "target/"},
	{pathContain: ".venv/"},
	{pathContain: "venv/"},
	{pathContain: "__pycache__/"},
	{pathContain: "vendor/"},
}

// isGeneratedFile checks if a file is auto-generated based on path patterns.
func (d *Detector) isGeneratedFile(path string) bool {
	for _, p := range generatedFilePatterns {
		if p.suffix != "" && strings.HasSuffix(path, p.suffix) {
			return true
		}
		if p.prefix != "" && strings.HasPrefix(filepath.Base(path), p.prefix) {
			return true
		}
		if p.contains != "" && strings.Contains(filepath.Base(path), p.contains) {
			return true
		}
		if p.pathContain != "" && strings.Contains(path, p.pathContain) {
			return true
		}
	}
	return false
}

// scoreAndMarkFiles scores importance and marks generated files.
func (d *Detector) scoreAndMarkFiles(tree *types.Directory, project *types.Project) {
	d.scoreDirectory(tree, project)
}

// scoreDirectory recursively scores files in a directory.
func (d *Detector) scoreDirectory(dir *types.Directory, project *types.Project) {
	for i := range dir.Files {
		file := &dir.Files[i]

		// Check if generated
		file.IsGenerated = d.isGeneratedFile(file.Path)

		// Score based on file role
		file.Importance = d.scoreFile(file.Path, file.Type, file.IsGenerated, project)
	}

	for i := range dir.Subdirs {
		d.scoreDirectory(&dir.Subdirs[i], project)
	}
}

// scoreFile calculates importance score for a file.
func (d *Detector) scoreFile(path, fileType string, isGenerated bool, project *types.Project) int {
	// Generated files get low scores
	if isGenerated {
		return 20
	}

	// Entry points get highest score
	for _, ep := range project.KeyFiles.EntryPoints {
		if path == ep {
			return 100
		}
	}

	// Config files get high score
	for _, cf := range project.KeyFiles.Configs {
		if path == cf {
			return 90
		}
	}

	// Core source files based on project type
	if d.isCoreSourceFile(path, project.Type) {
		return 80
	}

	// Test files
	for _, tf := range project.KeyFiles.Tests {
		if path == tf {
			return 50
		}
	}

	// Documentation
	for _, df := range project.KeyFiles.Docs {
		if path == df {
			return 40
		}
	}

	// Other source files
	if d.isSourceFile(fileType, project.Type) {
		return 70
	}

	// Default
	return 10
}

// isCoreSourceFile checks if a file is a core source file.
func (d *Detector) isCoreSourceFile(path string, pt types.ProjectType) bool {
	switch pt {
	case types.ProjectTypeGo:
		return strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go")
	case types.ProjectTypeNode:
		return (strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".ts")) &&
			!strings.Contains(path, "test") && !strings.Contains(path, "spec")
	case types.ProjectTypePython:
		return strings.HasSuffix(path, ".py") && !strings.HasPrefix(filepath.Base(path), "test_")
	case types.ProjectTypeRust:
		return strings.HasSuffix(path, ".rs") && !strings.Contains(path, "test")
	case types.ProjectTypeJava:
		return (strings.HasSuffix(path, ".java") || strings.HasSuffix(path, ".kt")) &&
			!strings.Contains(path, "Test")
	}
	return false
}

// isSourceFile checks if a file is a source file for the project type.
func (d *Detector) isSourceFile(fileType string, pt types.ProjectType) bool {
	switch pt {
	case types.ProjectTypeGo:
		return fileType == "go"
	case types.ProjectTypeNode:
		return fileType == "javascript" || fileType == "typescript"
	case types.ProjectTypePython:
		return fileType == "python"
	case types.ProjectTypeRust:
		return fileType == "rust"
	case types.ProjectTypeJava:
		return fileType == "java" || fileType == "kotlin"
	}
	return false
}

// calculateStats calculates project statistics.
func (d *Detector) calculateStats(tree *types.Directory, stats *types.Stats) {
	d.walkForStats(tree, stats)
}

// walkForStats walks the tree to calculate statistics.
func (d *Detector) walkForStats(dir *types.Directory, stats *types.Stats) {
	for _, file := range dir.Files {
		stats.FileCount++
		stats.TotalSize += file.Size
		stats.TotalLines += file.Lines

		if stats.ByType == nil {
			stats.ByType = make(map[string]int)
		}
		stats.ByType[file.Type]++
	}

	for _, subdir := range dir.Subdirs {
		d.walkForStats(&subdir, stats)
	}
}
