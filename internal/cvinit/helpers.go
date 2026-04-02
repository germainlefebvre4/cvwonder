package cvinit

import (
	"fmt"
	"strconv"
	"strings"
)

// splitLines splits a multiline string into a trimmed, non-empty string slice.
func splitLines(s string) []string {
	lines := strings.Split(s, "\n")
	result := make([]string, 0, len(lines))
	for _, l := range lines {
		if t := strings.TrimSpace(l); t != "" {
			result = append(result, t)
		}
	}
	return result
}

// splitComma splits a comma-separated string into a trimmed, non-empty string slice.
func splitComma(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}

// validateCompetencyLevel validates that s is an integer between 1 and 5 inclusive.
func validateCompetencyLevel(s string) error {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return fmt.Errorf("must be a number between 1 and 5")
	}
	if n < 1 || n > 5 {
		return fmt.Errorf("must be between 1 and 5 (got %d)", n)
	}
	return nil
}
