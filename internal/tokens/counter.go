// Package tokens provides token counting functionality.
//
// This package supports both approximate and accurate token counting
// to help AI agents understand token usage in outputs.
package tokens

import (
	"fmt"
)

// Counter defines the interface for token counting.
type Counter interface {
	Count(text string) int
}

// ApproximateCounter uses a simple heuristic: 4 chars ≈ 1 token.
// This provides a fast approximation without external dependencies.
type ApproximateCounter struct{}

// Count returns the approximate token count.
// Uses the rule: (len(text) + 3) / 4 for rounding up.
func (c *ApproximateCounter) Count(text string) int {
	return (len(text) + 3) / 4
}

// NewCounter creates a token counter based on the tokenizer type.
// Supported types: "approx" (default, no dependencies)
// Future: "tiktoken" for accurate GPT token counting
func NewCounter(tokenizer string) (Counter, error) {
	switch tokenizer {
	case "approx", "":
		return &ApproximateCounter{}, nil
	case "tiktoken":
		// Tiktoken support would require external dependency
		// For now, return error with helpful message
		return nil, fmt.Errorf("tiktoken tokenizer not available: use 'approx' or add github.com/pkoukk/tiktoken-go dependency")
	default:
		return nil, fmt.Errorf("unknown tokenizer: %s (use 'approx')", tokenizer)
	}
}
