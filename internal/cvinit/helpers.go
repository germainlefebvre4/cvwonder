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

// joinLines joins a string slice into a newline-separated string, the inverse of splitLines.
func joinLines(ss []string) string {
	return strings.Join(ss, "\n")
}

// joinComma joins a string slice into a ", "-separated string, the inverse of splitComma.
func joinComma(ss []string) string {
	return strings.Join(ss, ", ")
}

// confirmGateTitle computes an optional section's confirm-gate title and
// default value based on whether the section already has data (loaded via
// --resume). When it does, the gate defaults to true and reads "Review/edit
// ...?" so Enter proceeds into the section instead of skipping it.
func confirmGateTitle(hasExisting bool, addTitle, resumeTitle string) (string, bool) {
	if hasExisting {
		return resumeTitle, true
	}
	return addTitle, false
}

// forcesFirstEntry reports whether a loop section must collect a mandatory
// first entry before offering "add another?" — true only when the backing
// slice starts empty (i.e. not resumed with existing entries).
func forcesFirstEntry(existingCount int) bool {
	return existingCount == 0
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
