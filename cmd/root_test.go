package cmd

import (
	"reflect"
	"testing"
)

func TestParseExtensions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single extension without dot",
			input:    "go",
			expected: []string{".go"},
		},
		{
			name:     "single extension with dot",
			input:    ".go",
			expected: []string{".go"},
		},
		{
			name:     "multiple extensions mixed format",
			input:    "go,.md,yaml",
			expected: []string{".go", ".md", ".yaml"},
		},
		{
			name:     "uppercase normalized to lowercase",
			input:    ".GO,.MD",
			expected: []string{".go", ".md"},
		},
		{
			name:     "empty string returns nil",
			input:    "",
			expected: nil,
		},
		{
			name:     "whitespace trimmed",
			input:    "  go  ,  md  ",
			expected: []string{".go", ".md"},
		},
		{
			name:     "empty items skipped",
			input:    ",,go,,md,",
			expected: []string{".go", ".md"},
		},
		{
			name:     "complex extensions",
			input:    ".test.go,.spec.ts",
			expected: []string{".test.go", ".spec.ts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseExtensions(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseExtensions(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
