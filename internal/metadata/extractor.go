// Package metadata provides file metadata extraction functionality.
//
// This package handles extraction of metadata from files including:
//   - File size and modification time
//   - Line count for text files
//   - File type/language detection (based on extension and filename)
//   - Binary file detection
//
// Example usage:
//
//	extractor := metadata.NewExtractor()
//	file, err := extractor.Extract("/path/to/file.go", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("%s: %d lines, type=%s\n", file.Path, file.Lines, file.Type)
package metadata

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wesbragagt/gps/pkg/types"
)

// Extractor extracts metadata from files.
type Extractor struct {
	// MaxBytesToCheck is the maximum bytes to read for binary detection.
	MaxBytesToCheck int64
}

// NewExtractor creates a new Extractor with default settings.
func NewExtractor() *Extractor {
	return &Extractor{
		MaxBytesToCheck: 512,
	}
}

// Extract extracts metadata from a file and returns a populated File struct.
func (e *Extractor) Extract(filePath string, info os.FileInfo) (*types.File, error) {
	if info == nil {
		var err error
		info, err = os.Stat(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file %s: %w", filePath, err)
		}
	}

	file := &types.File{
		Path:    filePath,
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	// Detect if binary
	file.IsBinary = e.isBinaryFile(filePath, info)

	// Detect file type
	file.Type = e.detectFileType(filePath)

	// Count lines for text files
	if !file.IsBinary {
		lines, err := e.countLines(filePath)
		if err != nil {
			// Log warning but don't fail
			file.Lines = 0
		} else {
			file.Lines = lines
		}
	}

	return file, nil
}

// fileTypeMap maps file extensions to their types.
var fileTypeMap = map[string]string{
	// Programming languages
	".go":     "go",
	".js":     "javascript",
	".jsx":    "javascript",
	".ts":     "typescript",
	".tsx":    "typescript",
	".py":     "python",
	".pyw":    "python",
	".java":   "java",
	".kt":     "kotlin",
	".kts":    "kotlin",
	".rs":     "rust",
	".c":      "c",
	".h":      "c",
	".cpp":    "cpp",
	".cc":     "cpp",
	".cxx":    "cpp",
	".hpp":    "cpp",
	".hxx":    "cpp",
	".rb":     "ruby",
	".php":    "php",
	".swift":  "swift",
	".m":      "objective-c",
	".mm":     "objective-cpp",
	".cs":     "csharp",
	".scala":  "scala",
	".sc":     "scala",
	".ex":     "elixir",
	".exs":    "elixir",
	".erl":    "erlang",
	".hs":     "haskell",
	".lhs":    "haskell",
	".lua":    "lua",
	".r":      "r",
	".R":      "r",
	".pl":     "perl",
	".pm":     "perl",
	".sh":     "shell",
	".bash":   "shell",
	".zsh":    "shell",
	".ps1":    "powershell",
	".psm1":   "powershell",
	".vb":     "visualbasic",
	".vbs":    "visualbasic",
	".dart":   "dart",
	".groovy": "groovy",
	".gvy":    "groovy",
	".clj":    "clojure",
	".cljs":   "clojure",
	".cljc":   "clojure",
	".lisp":   "lisp",
	".lsp":    "lisp",
	".scm":    "scheme",
	".ss":     "scheme",
	".f90":    "fortran",
	".f95":    "fortran",
	".f03":    "fortran",
	".f":      "fortran",
	".for":    "fortran",
	".pas":    "pascal",
	".pp":     "pascal",
	".sql":    "sql",
	".asm":    "assembly",
	".s":      "assembly",
	".v":      "verilog",
	".vh":     "verilog",
	".vhd":    "vhdl",
	".vhdl":   "vhdl",

	// Markup and data formats
	".md":         "markdown",
	".markdown":   "markdown",
	".html":       "html",
	".htm":        "html",
	".xhtml":      "html",
	".xml":        "xml",
	".xsl":        "xml",
	".xslt":       "xml",
	".svg":        "svg",
	".yaml":       "yaml",
	".yml":        "yaml",
	".json":       "json",
	".jsonc":      "json",
	".json5":      "json",
	".toml":       "toml",
	".ini":        "ini",
	".cfg":        "config",
	".conf":       "config",
	".config":     "config",
	".env":        "env",
	".properties": "properties",

	// Stylesheets
	".css":  "css",
	".scss": "scss",
	".sass": "sass",
	".less": "less",
	".styl": "stylus",

	// Templates
	".tmpl":       "template",
	".template":   "template",
	".mustache":   "mustache",
	".hbs":        "handlebars",
	".handlebars": "handlebars",
	".ejs":        "ejs",
	".pug":        "pug",
	".jade":       "pug",
	".twig":       "twig",

	// Build and config files
	".gradle": "gradle",
	".mk":     "makefile",
	".cmake":  "cmake",

	// Documentation
	".rst":      "restructuredtext",
	".adoc":     "asciidoc",
	".asciidoc": "asciidoc",
	".tex":      "latex",
	".org":      "org",
	".txt":      "text",
	".text":     "text",

	// Other common files
	".csv":  "csv",
	".tsv":  "tsv",
	".log":  "log",
	".lock": "lock",
	".sum":  "checksum",

	// Web
	".vue":    "vue",
	".svelte": "svelte",
	".astro":  "astro",

	// Container/DevOps
	".dockerfile":    "dockerfile",
	".containerfile": "dockerfile",

	// Graphviz
	".dot": "dot",
	".gv":  "dot",

	// Other
	".proto":   "protobuf",
	".thrift":  "thrift",
	".graphql": "graphql",
	".gql":     "graphql",
}

// specialFileNames maps special filenames (without extension) to their types.
var specialFileNames = map[string]string{
	"dockerfile":       "dockerfile",
	"containerfile":    "dockerfile",
	"makefile":         "makefile",
	"gemfile":          "ruby",
	"rakefile":         "ruby",
	"brewfile":         "ruby",
	"podfile":          "ruby",
	"vagrantfile":      "ruby",
	"procfile":         "procfile",
	"jenkinsfile":      "groovy",
	"justfile":         "just",
	"makefile.PL":      "perl",
	"buildfile":        "ruby",
	"pkgfile":          "shell",
	"PKGBUILD":         "shell",
	".gitignore":       "gitignore",
	".gitattributes":   "git",
	".gitmodules":      "git",
	".gitkeep":         "git",
	".dockerignore":    "dockerignore",
	".editorconfig":    "editorconfig",
	".eslintrc":        "eslint",
	".prettierrc":      "prettier",
	".babelrc":         "babel",
	".npmrc":           "npm",
	".yarnrc":          "yarn",
	".pypirc":          "pypi",
	".pip.conf":        "pip",
	".cargo.toml":      "cargo",
	".rustfmt.toml":    "rust",
	".clang-format":    "clang",
	".gitlab-ci.yml":   "gitlab-ci",
	".travis.yml":      "travis",
	".circleci":        "circleci",
	"license":          "license",
	"license.md":       "license",
	"license.txt":      "license",
	"copying":          "license",
	"copying.md":       "license",
	"copying.txt":      "license",
	"readme":           "readme",
	"readme.md":        "readme",
	"readme.txt":       "readme",
	"changelog":        "changelog",
	"changelog.md":     "changelog",
	"changelog.txt":    "changelog",
	"authors":          "authors",
	"authors.md":       "authors",
	"authors.txt":      "authors",
	"contributors":     "contributors",
	"contributors.md":  "contributors",
	"contributors.txt": "contributors",
	"news":             "news",
	"news.md":          "news",
	"history":          "history",
	"history.md":       "history",
	"todo":             "todo",
	"todo.md":          "todo",
	"notice":           "notice",
	"notice.md":        "notice",
}

// detectFileType detects the file type/language from the file extension or name.
func (e *Extractor) detectFileType(path string) string {
	base := filepath.Base(path)
	lowerBase := strings.ToLower(base)

	// Check special filenames first (case-insensitive)
	if t, ok := specialFileNames[lowerBase]; ok {
		return t
	}

	// Check extension (case-insensitive)
	ext := strings.ToLower(filepath.Ext(path))
	if ext != "" {
		if t, ok := fileTypeMap[ext]; ok {
			return t
		}
	}

	// Handle files without extension
	if ext == "" {
		// Check if it's a dotfile
		if strings.HasPrefix(base, ".") {
			return "config"
		}
	}

	return "unknown"
}

// countLines counts the number of lines in a text file.
func (e *Extractor) countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to scan file: %w", err)
	}

	return lineCount, nil
}

// isBinaryFile detects if a file is binary by checking for null bytes.
func (e *Extractor) isBinaryFile(filePath string, info os.FileInfo) bool {
	// Directories are not binary files
	if info.IsDir() {
		return false
	}

	// Empty files are not binary
	if info.Size() == 0 {
		return false
	}

	// Check file extension first for known binary types
	ext := strings.ToLower(filepath.Ext(filePath))
	if isKnownBinaryExtension(ext) {
		return true
	}

	// Open file and check for null bytes
	file, err := os.Open(filePath)
	if err != nil {
		// If we can't open it, assume it might be binary
		return true
	}
	defer file.Close()

	// Read first N bytes
	bytesToRead := e.MaxBytesToCheck
	if info.Size() < bytesToRead {
		bytesToRead = info.Size()
	}

	buf := make([]byte, bytesToRead)
	n, err := io.ReadFull(file, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return true
	}

	// Check for null bytes (common indicator of binary files)
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true
		}
	}

	return false
}

// binaryExtensions contains extensions of known binary file types.
var binaryExtensions = map[string]bool{
	// Images
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".bmp":  true,
	".ico":  true,
	".webp": true,
	".tiff": true,
	".tif":  true,
	".psd":  true,
	".ai":   true,
	".eps":  true,
	".svg":  false, // SVG is text-based

	// Audio
	".mp3":  true,
	".mp4":  true,
	".wav":  true,
	".ogg":  true,
	".flac": true,
	".aac":  true,
	".m4a":  true,
	".wma":  true,
	".aiff": true,

	// Video
	".avi":  true,
	".mkv":  true,
	".mov":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
	".m4v":  true,
	".mpeg": true,
	".mpg":  true,

	// Archives
	".zip":  true,
	".tar":  true,
	".gz":   true,
	".bz2":  true,
	".xz":   true,
	".7z":   true,
	".rar":  true,
	".tgz":  true,
	".tbz2": true,

	// Executables
	".exe":   true,
	".dll":   true,
	".so":    true,
	".dylib": true,
	".a":     true,
	".lib":   true,
	".o":     true,
	".obj":   true,

	// Documents
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,

	// Fonts
	".ttf":   true,
	".otf":   true,
	".woff":  true,
	".woff2": true,
	".eot":   true,

	// Database
	".db":      true,
	".sqlite":  true,
	".sqlite3": true,

	// Other
	".pyc":   true,
	".pyo":   true,
	".pyd":   true,
	".class": true,
	".jar":   true,
	".war":   true,
	".ear":   true,
	".swf":   true,
	".iso":   true,
	".dmg":   true,
	".deb":   true,
	".rpm":   true,
	".msi":   true,
}

// isKnownBinaryExtension checks if the extension is a known binary type.
func isKnownBinaryExtension(ext string) bool {
	isBinary, ok := binaryExtensions[ext]
	return ok && isBinary
}

// HumanizeSize converts bytes to human-readable format.
func HumanizeSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.1fTB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.1fGB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1fMB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1fKB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

// ExtractSimple is a convenience function that extracts metadata without requiring an Extractor instance.
func ExtractSimple(filePath string) (*types.File, error) {
	return NewExtractor().Extract(filePath, nil)
}

// GetModTime returns the modification time of a file.
func GetModTime(filePath string) (time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to stat file: %w", err)
	}
	return info.ModTime(), nil
}
