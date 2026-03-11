// Package formatter provides test data for formatter tests.
package formatter

import (
	"time"

	"github.com/wesbragagt/gps/pkg/types"
)

// CreateTestProject creates a standard test project for formatter tests.
func CreateTestProject() *types.Project {
	return &types.Project{
		Name: "test-project",
		Type: types.ProjectTypeGo,
		Root: "/test/path",
		Stats: types.Stats{
			FileCount:  5,
			TotalSize:  10000,
			TotalLines: 500,
			ByType:     map[string]int{"go": 5},
		},
		KeyFiles: types.KeyFiles{
			EntryPoints: []string{"main.go"},
			Configs:     []string{"go.mod"},
			Tests:       []string{"main_test.go"},
			Docs:        []string{"README.md"},
		},
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "main.go", Size: 2000, Lines: 100, Type: "go", Importance: 100, ModTime: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)},
				{Path: "go.mod", Size: 200, Lines: 10, Type: "mod", Importance: 90, ModTime: time.Date(2024, 1, 10, 8, 0, 0, 0, time.UTC)},
			},
			Subdirs: []types.Directory{
				{
					Path: "internal",
					Files: []types.File{
						{Path: "internal/util.go", Size: 1500, Lines: 75, Type: "go", Importance: 70, ModTime: time.Date(2024, 1, 12, 14, 0, 0, 0, time.UTC)},
					},
					Subdirs: []types.Directory{
						{
							Path: "handler",
							Files: []types.File{
								{Path: "internal/handler/api.go", Size: 3000, Lines: 150, Type: "go", Importance: 80, ModTime: time.Date(2024, 1, 14, 9, 0, 0, 0, time.UTC)},
							},
						},
					},
				},
			},
		},
	}
}

// CreateMinimalProject creates a minimal project for edge case testing.
func CreateMinimalProject() *types.Project {
	return &types.Project{
		Name: "minimal",
		Type: types.ProjectTypeOther,
		Root: "/minimal",
		Tree: types.Directory{
			Path: ".",
		},
	}
}

// CreateProjectWithBinary creates a project with binary files.
func CreateProjectWithBinary() *types.Project {
	return &types.Project{
		Name: "binary-project",
		Type: types.ProjectTypeGo,
		Root: "/binary",
		Stats: types.Stats{
			FileCount:  2,
			TotalSize:  102400,
			TotalLines: 50,
		},
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "main.go", Size: 2048, Lines: 50, Type: "go", Importance: 100},
				{Path: "image.png", Size: 100000, Lines: 0, Type: "png", Importance: 20, IsBinary: true},
			},
		},
	}
}

// CreateProjectWithMultipleKeyFiles creates a project with multiple key files of each type.
func CreateProjectWithMultipleKeyFiles() *types.Project {
	return &types.Project{
		Name: "multi-keys",
		Type: types.ProjectTypeNode,
		Root: "/multi",
		KeyFiles: types.KeyFiles{
			EntryPoints: []string{"src/index.ts", "src/main.ts", "src/app.ts"},
			Configs:     []string{"package.json", "tsconfig.json", ".eslintrc"},
			Tests:       []string{"test/a.test.ts", "test/b.test.ts"},
			Docs:        []string{"README.md", "CHANGELOG.md"},
		},
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "package.json", Size: 500, Lines: 20, Type: "json", Importance: 100},
			},
		},
	}
}

// CreateLargeProject creates a project with large files for size formatting tests.
func CreateLargeProject() *types.Project {
	return &types.Project{
		Name: "large-project",
		Type: types.ProjectTypeMixed,
		Root: "/large",
		Stats: types.Stats{
			FileCount:  100,
			TotalSize:  5368709120, // 5GB
			TotalLines: 500000,
		},
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "large.dat", Size: 1073741824, Lines: 0, Type: "binary", IsBinary: true}, // 1GB
				{Path: "medium.bin", Size: 52428800, Lines: 0, Type: "binary", IsBinary: true},  // 50MB
				{Path: "small.log", Size: 1536, Lines: 50, Type: "log"},                         // 1.5KB
			},
		},
	}
}

// CreateDeepNestedProject creates a project with deeply nested directories.
func CreateDeepNestedProject() *types.Project {
	return &types.Project{
		Name: "deep-project",
		Type: types.ProjectTypeGo,
		Root: "/deep",
		Tree: types.Directory{
			Path: ".",
			Files: []types.File{
				{Path: "root.go", Size: 100, Lines: 10, Type: "go"},
			},
			Subdirs: []types.Directory{
				{
					Path: "a",
					Files: []types.File{
						{Path: "a/level1.go", Size: 100, Lines: 10, Type: "go"},
					},
					Subdirs: []types.Directory{
						{
							Path: "b",
							Files: []types.File{
								{Path: "a/b/level2.go", Size: 100, Lines: 10, Type: "go"},
							},
							Subdirs: []types.Directory{
								{
									Path: "c",
									Files: []types.File{
										{Path: "a/b/c/level3.go", Size: 100, Lines: 10, Type: "go"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// CreateEmptyProject creates a project with empty fields.
func CreateEmptyProject() *types.Project {
	return &types.Project{
		Name: "empty-project",
		Type: "",
		Root: "",
		Tree: types.Directory{
			Path: ".",
		},
	}
}
