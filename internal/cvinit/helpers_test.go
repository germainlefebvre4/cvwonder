package cvinit

import (
"testing"

"github.com/stretchr/testify/assert"
)

func TestSplitLines(t *testing.T) {
tests := []struct {
name     string
input    string
expected []string
}{
{"two lines", "line one\nline two", []string{"line one", "line two"}},
{"empty lines discarded", "line one\n\n\nline two", []string{"line one", "line two"}},
{"only whitespace line discarded", "line one\n   \nline two", []string{"line one", "line two"}},
{"empty string", "", []string{}},
{"single line", "only one", []string{"only one"}},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
result := splitLines(tt.input)
assert.Equal(t, tt.expected, result)
})
}
}

func TestSplitComma(t *testing.T) {
tests := []struct {
name     string
input    string
expected []string
}{
{"standard csv", "Go, Docker, Kubernetes", []string{"Go", "Docker", "Kubernetes"}},
{"no spaces", "Go,Docker,Kubernetes", []string{"Go", "Docker", "Kubernetes"}},
{"extra whitespace", "  Go  ,  Docker  ", []string{"Go", "Docker"}},
{"empty string", "", []string{}},
{"single value", "Go", []string{"Go"}},
{"trailing comma ignored", "Go,Docker,", []string{"Go", "Docker"}},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
result := splitComma(tt.input)
assert.Equal(t, tt.expected, result)
})
}
}

func TestJoinLines_RoundTripsWithSplitLines(t *testing.T) {
tests := [][]string{
{"line one", "line two"},
{"only one"},
{},
}

for _, ss := range tests {
joined := joinLines(ss)
assert.Equal(t, ss, splitLines(joined))
}
}

func TestJoinComma_RoundTripsWithSplitComma(t *testing.T) {
tests := [][]string{
{"Go", "Docker", "Kubernetes"},
{"Go"},
{},
}

for _, ss := range tests {
joined := joinComma(ss)
assert.Equal(t, ss, splitComma(joined))
}
}

func TestConfirmGateTitle(t *testing.T) {
t.Run("no existing data uses add title and defaults false", func(t *testing.T) {
title, def := confirmGateTitle(false, "Add X?", "Review/edit X?")
assert.Equal(t, "Add X?", title)
assert.False(t, def)
})

t.Run("existing data uses resume title and defaults true", func(t *testing.T) {
title, def := confirmGateTitle(true, "Add X?", "Review/edit X?")
assert.Equal(t, "Review/edit X?", title)
assert.True(t, def)
})
}

func TestForcesFirstEntry(t *testing.T) {
assert.True(t, forcesFirstEntry(0))
assert.False(t, forcesFirstEntry(1))
assert.False(t, forcesFirstEntry(5))
}

func TestValidateCompetencyLevel(t *testing.T) {
	t.Run("valid levels from 0 to 100", func(t *testing.T) {
		for _, v := range []string{"0", "1", "50", "80", "99", "100"} {
			assert.NoError(t, validateCompetencyLevel(v), v)
		}
	})

	t.Run("whitespace trimmed", func(t *testing.T) {
		assert.NoError(t, validateCompetencyLevel("  60  "))
	})

	t.Run("non-numeric rejected", func(t *testing.T) {
		err := validateCompetencyLevel("abc")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a number")
	})

	t.Run("above 100 rejected", func(t *testing.T) {
		err := validateCompetencyLevel("101")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "between 0 and 100")
	})

	t.Run("negative rejected", func(t *testing.T) {
		assert.Error(t, validateCompetencyLevel("-1"))
	})
}
