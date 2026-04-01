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

func TestValidateCompetencyLevel(t *testing.T) {
t.Run("valid levels 1 to 5", func(t *testing.T) {
for _, v := range []string{"1", "2", "3", "4", "5"} {
assert.NoError(t, validateCompetencyLevel(v))
}
})

t.Run("whitespace trimmed", func(t *testing.T) {
assert.NoError(t, validateCompetencyLevel("  3  "))
})

t.Run("non-numeric rejected", func(t *testing.T) {
err := validateCompetencyLevel("abc")
assert.Error(t, err)
assert.Contains(t, err.Error(), "must be a number")
})

t.Run("zero rejected", func(t *testing.T) {
err := validateCompetencyLevel("0")
assert.Error(t, err)
})

t.Run("six rejected", func(t *testing.T) {
err := validateCompetencyLevel("6")
assert.Error(t, err)
})

t.Run("negative rejected", func(t *testing.T) {
err := validateCompetencyLevel("-1")
assert.Error(t, err)
})
}
