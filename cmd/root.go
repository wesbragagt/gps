// Package cmd provides the CLI commands for gps.
//
// The root command handles project scanning and formatting.
// It supports tree-compatible flags (-L, -d, -a, -I, -P) and
// AI-optimized flags (-f, --summary, --entry-points, --focus).
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wesbragagt/gps/internal/detector"
	"github.com/wesbragagt/gps/internal/formatter"
	"github.com/wesbragagt/gps/internal/metadata"
	"github.com/wesbragagt/gps/internal/scanner"
	"github.com/wesbragagt/gps/internal/tokens"
	"github.com/wesbragagt/gps/pkg/types"
)

var (
	// Tree-compatible flags
	flagDepth      int
	flagDirsOnly   bool
	flagAll        bool
	flagExclude    string
	flagPattern    string
	flagExt        string
	flagExcludeExt string

	// AI-optimized flags
	flagFormat        string
	flagMeta          bool
	flagNoMeta        bool
	flagProjectInfo   bool
	flagSummary       bool
	flagEntryPoints   bool
	flagSmartTraverse bool
	flagFocus         string

	// Token counting flags
	flagTokens    bool
	flagCompare   bool
	flagTokenizer string

	// Target path
	flagPath string
)

var rootCmd = &cobra.Command{
	Use:   "gps [path]",
	Short: "Project structure navigator for AI agents",
	Long: `gps (go-project-gps) provides a tree-like experience optimized for AI agents.

If you know how to use 'tree', you know how to use 'gps'.

Examples:
  gps                    # Show current directory (TOON format)
  gps -L 2               # Limit depth to 2
  gps -d                 # Directories only
  gps -a                 # Include hidden files
  gps -I "node_modules"  # Exclude pattern
  gps -e go,md           # Only Go and Markdown files
  gps --exclude-ext log,tmp  # Exclude log and tmp files
  gps -f json            # JSON output format
  gps --focus src        # Focus on subdirectory`,
	Args: cobra.MaximumNArgs(1),
	RunE: runScan,
}

func init() {
	// Tree-compatible flags
	rootCmd.Flags().IntVarP(&flagDepth, "depth", "L", -1, "Max depth of traversal (-1 for unlimited)")
	rootCmd.Flags().BoolVarP(&flagDirsOnly, "dirs-only", "d", false, "List directories only")
	rootCmd.Flags().BoolVarP(&flagAll, "all", "a", false, "Include hidden files")
	rootCmd.Flags().StringVarP(&flagExclude, "exclude", "I", "", "Exclude files matching pattern")
	rootCmd.Flags().StringVarP(&flagPattern, "pattern", "P", "", "Include only files matching pattern")
	rootCmd.Flags().StringVarP(&flagExt, "ext", "e", "", "Include only files with these extensions (comma-separated)")
	rootCmd.Flags().StringVar(&flagExcludeExt, "exclude-ext", "", "Exclude files with these extensions (comma-separated)")

	// AI-optimized flags
	rootCmd.Flags().StringVarP(&flagFormat, "format", "f", "toon", "Output format: toon, json, tree, flat")
	rootCmd.Flags().BoolVar(&flagMeta, "meta", true, "Include file metadata")
	rootCmd.Flags().BoolVar(&flagNoMeta, "no-meta", false, "Disable metadata (overrides --meta)")
	rootCmd.Flags().BoolVar(&flagProjectInfo, "project-info", true, "Include project detection info")
	rootCmd.Flags().BoolVarP(&flagSummary, "summary", "s", false, "High-level project overview only")
	rootCmd.Flags().BoolVar(&flagEntryPoints, "entry-points", false, "Show only entry points")
	rootCmd.Flags().BoolVar(&flagSmartTraverse, "smart-traverse", true, "Intelligent traversal based on project type")
	rootCmd.Flags().StringVar(&flagFocus, "focus", "", "Focus on specific subdirectory")

	// Token counting flags
	rootCmd.Flags().BoolVarP(&flagTokens, "tokens", "t", false, "Show token count in output")
	rootCmd.Flags().BoolVar(&flagCompare, "compare", false, "Compare token counts across all formats")
	rootCmd.Flags().StringVar(&flagTokenizer, "tokenizer", "approx", "Tokenizer type: approx")

	// Initialize Viper for config file support
	initConfig()
}

// initConfig sets up Viper for config file handling
func initConfig() {
	// Set default values
	viper.SetDefault("format", "toon")
	viper.SetDefault("depth", -1)
	viper.SetDefault("all", false)
	viper.SetDefault("dirs-only", false)
	viper.SetDefault("meta", true)
	viper.SetDefault("project-info", true)
	viper.SetDefault("summary", false)
	viper.SetDefault("entry-points", false)
	viper.SetDefault("smart-traverse", true)
	viper.SetDefault("focus", "")
	viper.SetDefault("exclude", []string{})
	viper.SetDefault("pattern", []string{})
	viper.SetDefault("ext", []string{})
	viper.SetDefault("exclude-ext", []string{})

	// Config file settings
	viper.SetConfigName(".gps")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	// Bind flags to Viper (flags override config)
	viper.BindPFlag("format", rootCmd.Flags().Lookup("format"))
	viper.BindPFlag("depth", rootCmd.Flags().Lookup("depth"))
	viper.BindPFlag("dirs-only", rootCmd.Flags().Lookup("dirs-only"))
	viper.BindPFlag("all", rootCmd.Flags().Lookup("all"))
	viper.BindPFlag("exclude", rootCmd.Flags().Lookup("exclude"))
	viper.BindPFlag("pattern", rootCmd.Flags().Lookup("pattern"))
	viper.BindPFlag("meta", rootCmd.Flags().Lookup("meta"))
	viper.BindPFlag("project-info", rootCmd.Flags().Lookup("project-info"))
	viper.BindPFlag("summary", rootCmd.Flags().Lookup("summary"))
	viper.BindPFlag("entry-points", rootCmd.Flags().Lookup("entry-points"))
	viper.BindPFlag("smart-traverse", rootCmd.Flags().Lookup("smart-traverse"))
	viper.BindPFlag("focus", rootCmd.Flags().Lookup("focus"))
	viper.BindPFlag("ext", rootCmd.Flags().Lookup("ext"))
	viper.BindPFlag("exclude-ext", rootCmd.Flags().Lookup("exclude-ext"))

	// Read config file (ignore if not found)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "warning: error reading config file: %v\n", err)
		}
	}
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetOutput sets the output streams
func SetOutput(stdout, stderr *os.File) {
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
}

// Formatter interface for output formatting
type Formatter interface {
	Format(project *types.Project) (string, error)
}

// runScan is the main scan workflow
func runScan(cmd *cobra.Command, args []string) error {
	// Handle --no-meta override
	if flagNoMeta {
		viper.Set("meta", false)
	}

	// Determine target path
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	// Apply --focus flag
	focus := viper.GetString("focus")
	if focus != "" {
		targetPath = filepath.Join(targetPath, focus)
	}

	// Clean and validate path
	targetPath = filepath.Clean(targetPath)
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("cannot access %q: %w", targetPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", targetPath)
	}

	// Get exclude patterns from Viper (could be from config file or flag)
	excludePatterns := viper.GetStringSlice("exclude")
	if flagExclude != "" {
		// Flag takes precedence - parse and override
		excludePatterns = parsePatterns(flagExclude)
	}

	includePatterns := viper.GetStringSlice("pattern")
	if flagPattern != "" {
		includePatterns = parsePatterns(flagPattern)
	}

	includeExtensions := viper.GetStringSlice("ext")
	if flagExt != "" {
		includeExtensions = parseExtensions(flagExt)
	}

	excludeExtensions := viper.GetStringSlice("exclude-ext")
	if flagExcludeExt != "" {
		excludeExtensions = parseExtensions(flagExcludeExt)
	}

	// Create scanner config from Viper (flags > config > defaults)
	config := &scanner.Config{
		MaxDepth:          viper.GetInt("depth"),
		IncludeHidden:     viper.GetBool("all"),
		ExcludePatterns:   excludePatterns,
		IncludePatterns:   includePatterns,
		IncludeExtensions: includeExtensions,
		ExcludeExtensions: excludeExtensions,
		RespectGitignore:  true,
	}

	// Create and run scanner
	s := scanner.New(config)
	tree, err := s.Walk(targetPath)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Report non-fatal scanner errors to stderr
	for _, scanErr := range s.Errors() {
		fmt.Fprintf(os.Stderr, "warning: %v\n", scanErr)
	}

	// Extract metadata if enabled
	if viper.GetBool("meta") {
		extractMetadata(tree, targetPath)
	}

	// Detect project info
	det := detector.New()
	project, err := det.Detect(targetPath, tree)
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	// Handle --entry-points mode
	if viper.GetBool("entry-points") {
		output := formatEntryPoints(&project.KeyFiles)
		fmt.Println(output)
		return nil
	}

	// Handle --summary mode
	if viper.GetBool("summary") {
		output := formatSummary(project)
		fmt.Println(output)
		return nil
	}

	// Handle token counting modes
	if flagTokens || flagCompare {
		counter, err := tokens.NewCounter(flagTokenizer)
		if err != nil {
			return err
		}

		if flagCompare {
			return compareAllFormats(project, counter)
		}

		// Show token count for current format
		f := getFormatter(viper.GetString("format"))
		output, err := f.Format(project)
		if err != nil {
			return fmt.Errorf("formatting failed: %w", err)
		}
		tokenCount := counter.Count(output)

		fmt.Print(output)
		fmt.Fprintf(os.Stderr, "\n---\nTokens: %d (approx)\nBytes: %d\n", tokenCount, len(output))
		return nil
	}

	// Get formatter and format output
	f := getFormatter(viper.GetString("format"))
	output, err := f.Format(project)
	if err != nil {
		return fmt.Errorf("formatting failed: %w", err)
	}

	// Print output
	fmt.Print(output)

	return nil
}

// parsePatterns converts comma-separated patterns to slice
func parsePatterns(patterns string) []string {
	if patterns == "" {
		return nil
	}

	var result []string
	for _, p := range strings.Split(patterns, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// parseExtensions converts comma-separated extensions to normalized slice.
// Ensures each extension has a leading dot and is lowercase.
func parseExtensions(extensions string) []string {
	if extensions == "" {
		return nil
	}

	var result []string
	for _, ext := range strings.Split(extensions, ",") {
		ext = strings.TrimSpace(ext)
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		result = append(result, strings.ToLower(ext))
	}
	return result
}

// extractMetadata recursively extracts metadata for all files
func extractMetadata(dir *types.Directory, rootPath string) {
	extractor := metadata.NewExtractor()

	for i := range dir.Files {
		fullPath := filepath.Join(rootPath, dir.Files[i].Path)
		meta, err := extractor.Extract(fullPath, nil)
		if err != nil {
			// Log warning to stderr, continue
			continue
		}
		// Preserve path and modtime from scan, update other fields
		dir.Files[i].Lines = meta.Lines
		dir.Files[i].Type = meta.Type
		dir.Files[i].IsBinary = meta.IsBinary
	}

	for i := range dir.Subdirs {
		extractMetadata(&dir.Subdirs[i], rootPath)
	}
}

// getFormatter returns the appropriate formatter for the format name
func getFormatter(format string) Formatter {
	switch format {
	case "toon":
		return formatter.NewToonFormatter()
	case "json":
		return &formatter.JsonFormatter{
			PrettyPrint: true,
			Compact:     false,
		}
	case "tree":
		return &formatter.TreeFormatter{
			ShowMetadata: viper.GetBool("meta"),
			Colorize:     true,
			ShowSummary:  viper.GetBool("project-info"),
		}
	case "flat":
		return &formatter.FlatFormatter{
			ShowHeader:     true,
			SortBy:         "path",
			SortDescending: false,
			ShowFields:     []string{"path", "size", "lines", "type", "importance"},
		}
	default:
		// Default to TOON
		return formatter.NewToonFormatter()
	}
}

// formatEntryPoints formats entry points output
func formatEntryPoints(kf *types.KeyFiles) string {
	var sb strings.Builder
	sb.WriteString("entry-points{\n")

	if len(kf.EntryPoints) > 0 {
		if len(kf.EntryPoints) == 1 {
			sb.WriteString("  main: ")
			sb.WriteString(kf.EntryPoints[0])
			sb.WriteString("\n")
		} else {
			for _, ep := range kf.EntryPoints {
				sb.WriteString("  - ")
				sb.WriteString(ep)
				sb.WriteString("\n")
			}
		}
	} else {
		sb.WriteString("  # No entry points detected\n")
	}

	sb.WriteString("}")
	return sb.String()
}

// formatSummary formats summary output
func formatSummary(project *types.Project) string {
	var sb strings.Builder

	sb.WriteString("project[")
	sb.WriteString(project.Name)
	sb.WriteString("]{\n")

	sb.WriteString("  type: ")
	sb.WriteString(string(project.Type))
	sb.WriteString("\n")

	if project.Stats.FileCount > 0 {
		sb.WriteString("  files: ")
		sb.WriteString(fmt.Sprintf("%d", project.Stats.FileCount))
		sb.WriteString("\n")
	}

	if project.Stats.TotalSize > 0 {
		sb.WriteString("  size: ")
		sb.WriteString(humanizeSize(project.Stats.TotalSize))
		sb.WriteString("\n")
	}

	if len(project.KeyFiles.EntryPoints) > 0 {
		sb.WriteString("  entry: ")
		sb.WriteString(project.KeyFiles.EntryPoints[0])
		sb.WriteString("\n")
	}

	if len(project.KeyFiles.Tests) > 0 {
		sb.WriteString("  tests: ")
		sb.WriteString(fmt.Sprintf("%d files", len(project.KeyFiles.Tests)))
		sb.WriteString("\n")
	}

	if len(project.KeyFiles.Docs) > 0 {
		sb.WriteString("  docs: ")
		// Find README
		for _, doc := range project.KeyFiles.Docs {
			if strings.Contains(strings.ToLower(doc), "readme") {
				sb.WriteString(doc)
				break
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("}")
	return sb.String()
}

// humanizeSize converts bytes to human-readable format
func humanizeSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
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

// compareAllFormats outputs a comparison table of all format token counts
func compareAllFormats(project *types.Project, counter tokens.Counter) error {
	formats := []struct {
		name string
		fmt  Formatter
	}{
		{"TOON", formatter.NewToonFormatter()},
		{"JSON", &formatter.JsonFormatter{PrettyPrint: true}},
		{"JSON (compact)", &formatter.JsonFormatter{Compact: true}},
		{"Tree", &formatter.TreeFormatter{ShowMetadata: true, Colorize: false}},
		{"Flat", &formatter.FlatFormatter{}},
	}

	type result struct {
		name   string
		tokens int
		bytes  int
	}
	results := make([]result, len(formats))

	for i, f := range formats {
		output, err := f.fmt.Format(project)
		if err != nil {
			return fmt.Errorf("format %s failed: %w", f.name, err)
		}
		results[i] = result{
			name:   f.name,
			tokens: counter.Count(output),
			bytes:  len(output),
		}
	}

	// Print comparison table
	fmt.Println("Format          Tokens  Bytes   vs TOON")
	fmt.Println("-------         ------  -----   -------")

	baseTokens := results[0].tokens
	for _, r := range results {
		diff := ""
		if r.tokens != baseTokens && baseTokens > 0 {
			pct := float64(r.tokens-baseTokens) / float64(baseTokens) * 100
			if pct > 0 {
				diff = fmt.Sprintf("+%.0f%%", pct)
			} else {
				diff = fmt.Sprintf("%.0f%%", pct)
			}
		} else if r.tokens == baseTokens {
			diff = "baseline"
		}
		fmt.Printf("%-16s %6d  %5d   %s\n", r.name, r.tokens, r.bytes, diff)
	}

	fmt.Fprintf(os.Stderr, "\nNote: Token counts are approximate (4 chars ≈ 1 token)\n")
	return nil
}
