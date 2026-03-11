// Package formatter provides output formatters for project structure.
package formatter

import (
	"encoding/json"

	"github.com/wesbragagt/gps/pkg/types"
)

// JsonFormatter formats projects as JSON.
type JsonFormatter struct {
	// PrettyPrint indents output with 2 spaces when true.
	PrettyPrint bool

	// Compact produces minimal output with no whitespace.
	Compact bool
}

// NewJsonFormatter creates a new JSON formatter with default settings.
func NewJsonFormatter() *JsonFormatter {
	return &JsonFormatter{
		PrettyPrint: true,
		Compact:     false,
	}
}

// jsonOutput wraps the project for JSON output.
type jsonOutput struct {
	Project *types.Project `json:"project"`
}

// Format converts a Project to JSON format string.
func (f *JsonFormatter) Format(project *types.Project) (string, error) {
	if project == nil {
		return "{}", nil
	}

	output := jsonOutput{
		Project: project,
	}

	var data []byte
	var err error

	// Compact mode overrides pretty print
	if f.Compact {
		data, err = json.Marshal(output)
	} else if f.PrettyPrint {
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		data, err = json.Marshal(output)
	}

	if err != nil {
		return "", err
	}

	return string(data), nil
}
