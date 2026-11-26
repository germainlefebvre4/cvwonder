package ui

import (
	"fmt"

	"github.com/germainlefebvre4/cvwonder/internal/validator"
)

// ReportValidation renders a validation result: a passing result is a single
// success line; errors are listed with line, field, issue and suggestion,
// followed by any warnings. When final is true the outcome is reported as a
// Result (kept by --quiet), otherwise as a Step of a larger command.
func ReportValidation(r Reporter, res *validator.ValidationResult, final bool) {
	emit := r.Step
	if final {
		emit = r.Result
	}

	switch {
	case res.Valid && len(res.Warnings) == 0:
		emit(Success, "Validation passed: the CV YAML file is valid")
		return
	case res.Valid:
		emit(Warning, fmt.Sprintf("Validation passed with %s", plural(len(res.Warnings), "warning")))
	default:
		emit(Failure, fmt.Sprintf("Validation failed: %s", plural(len(res.Errors), "error")))
		for i, e := range res.Errors {
			r.Detail(fmt.Sprintf("Error %d:", i+1))
			if e.LineNumber > 0 {
				r.Detail(fmt.Sprintf("  Line: %d", e.LineNumber))
			}
			r.Detail(fmt.Sprintf("  Field: %s", e.Field), fmt.Sprintf("  Issue: %s", e.Message))
			if e.Suggestion != "" {
				r.Detail(fmt.Sprintf("  Suggestion: %s", e.Suggestion))
			}
		}
		if len(res.Warnings) > 0 {
			r.Detail("Warnings:")
		}
	}

	for i, w := range res.Warnings {
		r.Detail(fmt.Sprintf("Warning %d:", i+1), fmt.Sprintf("  Field: %s", w.Field), fmt.Sprintf("  Message: %s", w.Message))
		if w.Suggestion != "" {
			r.Detail(fmt.Sprintf("  Suggestion: %s", w.Suggestion))
		}
	}
}

func plural(n int, word string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}
