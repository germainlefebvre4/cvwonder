package theme_config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/goccy/go-yaml"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

type ThemeConfig struct {
	Name           string                 `yaml:"name"`
	Slug           string                 `yaml:"slug"`
	Description    string                 `yaml:"description"`
	Author         string                 `yaml:"author"`
	MinimumVersion string                 `yaml:"minimumVersion"`
	Configuration  map[string]interface{} `yaml:"configuration"`
}

// NormalizeConfigKeys recursively normalizes all map keys to camelCase (lowercases
// the first character of each key). It also coerces map[interface{}]interface{}
// (which goccy/go-yaml may produce in some cases) to map[string]interface{}.
func NormalizeConfigKeys(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(m))
	for k, v := range m {
		normalized := toCamelCase(k)
		result[normalized] = normalizeValue(v)
	}
	return result
}

func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func normalizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		return NormalizeConfigKeys(val)
	case map[interface{}]interface{}:
		converted := make(map[string]interface{}, len(val))
		for k, v2 := range val {
			converted[toCamelCase(fmt.Sprintf("%v", k))] = normalizeValue(v2)
		}
		return converted
	default:
		return v
	}
}

func GetThemeConfigFromURL(githubRepo GithubRepo) (ThemeConfig, error) {
	// Download theme.yaml
	client := utils.GetGitHubClient()

	// Create options with ref if specified
	var opts *github.RepositoryContentGetOptions
	if githubRepo.Ref != "" {
		opts = &github.RepositoryContentGetOptions{
			Ref: githubRepo.Ref,
		}
	}

	fileContent, _, _, err := client.Repositories.GetContents(context.TODO(), githubRepo.Owner, githubRepo.Name, "theme.yaml", opts)
	if err != nil {
		return ThemeConfig{}, fmt.Errorf("downloading theme.yaml: %w", err)
	}

	// Read theme.yaml
	config, err := fileContent.GetContent()
	if err != nil {
		return ThemeConfig{}, fmt.Errorf("reading theme.yaml: %w", err)
	}

	// Parse theme.yaml
	themeConfig := ThemeConfig{}
	err = yaml.Unmarshal([]byte(config), &themeConfig)
	if err != nil {
		return ThemeConfig{}, fmt.Errorf("parsing theme.yaml: %w", err)
	}

	// Normalize configuration keys to camelCase
	if themeConfig.Configuration != nil {
		themeConfig.Configuration = NormalizeConfigKeys(themeConfig.Configuration)
	}

	return themeConfig, nil
}

func GetThemeConfigFromDir(dir string) (ThemeConfig, error) {
	// Read theme.yaml
	config, err := os.ReadFile(dir + "/theme.yaml")

	if err != nil {
		return ThemeConfig{}, fmt.Errorf("reading theme.yaml in %s: %w", dir, err)
	}

	// Parse theme.yaml
	themeConfig := ThemeConfig{}
	err = yaml.Unmarshal(config, &themeConfig)
	if err != nil {
		return ThemeConfig{}, fmt.Errorf("parsing theme.yaml in %s: %w", dir, err)
	}

	// Normalize configuration keys to camelCase
	if themeConfig.Configuration != nil {
		themeConfig.Configuration = NormalizeConfigKeys(themeConfig.Configuration)
	}

	return themeConfig, nil
}

func GetThemeConfigFromThemeName(themeName string) (ThemeConfig, error) {
	// This function now expects themeName to already be resolved to the actual directory name
	// The caller should use GetThemeDirectory from themes package first
	return GetThemeConfigFromDir("themes/" + themeName)
}

// VerifyThemeMinimumVersion reports whether cvwonderVersion satisfies the theme's
// minimum version. Versions are compared numerically (0.9.0 < 0.10.0) and an
// optional "v" prefix is ignored. A missing minimum, or a version that is not
// numeric (a "dev" build, "latest"), cannot be judged and is treated as
// satisfying the requirement.
func (tc *ThemeConfig) VerifyThemeMinimumVersion(cvwonderVersion string) bool {
	if strings.TrimSpace(tc.MinimumVersion) == "" {
		return true
	}
	minimum, ok := parseVersion(tc.MinimumVersion)
	if !ok {
		logrus.Debugf("Theme minimum version %q is not a version number, skipping the check", tc.MinimumVersion)
		return true
	}
	current, ok := parseVersion(cvwonderVersion)
	if !ok {
		logrus.Debugf("CV Wonder version %q is not a version number, skipping the theme minimum version check", cvwonderVersion)
		return true
	}
	if compareVersions(minimum, current) <= 0 {
		return true
	}
	logrus.Error("CV Wonder version: ", cvwonderVersion)
	logrus.Error("Theme minimum version: ", tc.MinimumVersion)
	logrus.Error("The theme minimum version not met. You might encounter issues with this theme.")
	logrus.Error("")
	return false
}

// version is a parsed X.Y.Z[-prerelease] version; build metadata is dropped.
type version struct {
	numbers    [3]int
	prerelease string
}

// parseVersion parses "1", "1.2", "1.2.3", "v1.2.3", "1.2.3-rc.1" and
// "1.2.3+build". It returns false for anything else (e.g. "dev").
func parseVersion(raw string) (version, bool) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "v"), "V")
	if i := strings.IndexByte(s, '+'); i >= 0 {
		s = s[:i]
	}
	var v version
	if i := strings.IndexByte(s, '-'); i >= 0 {
		v.prerelease = s[i+1:]
		s = s[:i]
	}
	parts := strings.Split(s, ".")
	if s == "" || len(parts) > len(v.numbers) {
		return version{}, false
	}
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 {
			return version{}, false
		}
		v.numbers[i] = n
	}
	return v, true
}

// compareVersions returns -1, 0 or 1. A pre-release sorts before its release.
func compareVersions(a, b version) int {
	for i := range a.numbers {
		if a.numbers[i] != b.numbers[i] {
			if a.numbers[i] < b.numbers[i] {
				return -1
			}
			return 1
		}
	}
	switch {
	case a.prerelease == b.prerelease:
		return 0
	case a.prerelease == "":
		return 1
	case b.prerelease == "":
		return -1
	default:
		return strings.Compare(a.prerelease, b.prerelease)
	}
}

func GetDefaultBranch(githubRepo GithubRepo) (string, error) {
	client := utils.GetGitHubClient()
	repo, _, err := client.Repositories.Get(context.TODO(), githubRepo.Owner, githubRepo.Name)
	if err != nil {
		return "", fmt.Errorf("getting repository info: %w", err)
	}
	return repo.GetDefaultBranch(), nil
}

// DeepMerge merges src into dst at the leaf level. Values in src override dst.
// Nested maps are merged recursively so sibling keys in dst are preserved.
func DeepMerge(dst, src map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(dst))
	for k, v := range dst {
		result[k] = v
	}
	for k, v := range src {
		if srcMap, ok := v.(map[string]interface{}); ok {
			if dstMap, ok := result[k].(map[string]interface{}); ok {
				result[k] = DeepMerge(dstMap, srcMap)
				continue
			}
		}
		result[k] = v
	}
	return result
}

// ParseConfigOverrides parses a slice of "key=value" strings (where key may use
// dot-notation for nesting) and deep-merges the resulting overrides on top of base.
// Keys are camelCase-normalized. Values are auto-coerced via YAML parsing.
func ParseConfigOverrides(overrides []string, base map[string]interface{}) (map[string]interface{}, error) {
	// Start with a copy of base (may be nil)
	merged := make(map[string]interface{}, len(base))
	for k, v := range base {
		merged[k] = v
	}

	for _, override := range overrides {
		// Split on the first '=' only
		idx := -1
		for i, c := range override {
			if c == '=' {
				idx = i
				break
			}
		}
		if idx < 0 {
			return nil, fmt.Errorf("invalid --config value %q: expected key=value format", override)
		}
		rawKey := override[:idx]
		rawValue := override[idx+1:]

		// Coerce value via YAML parsing
		var coercedValue interface{}
		if err := yaml.Unmarshal([]byte(rawValue), &coercedValue); err != nil {
			coercedValue = rawValue
		}

		// Build a nested map from the dot-notation key, with camelCase normalization
		segments := splitDotKey(rawKey)
		if len(segments) == 0 {
			return nil, fmt.Errorf("invalid --config value %q: key must contain at least one non-empty segment", override)
		}
		overrideMap := buildNestedMap(segments, coercedValue)
		merged = DeepMerge(merged, overrideMap)
	}

	return merged, nil
}

// splitDotKey splits a dot-notation key into camelCase-normalized segments.
func splitDotKey(key string) []string {
	var segments []string
	current := ""
	for _, c := range key {
		if c == '.' {
			if current != "" {
				segments = append(segments, toCamelCase(current))
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		segments = append(segments, toCamelCase(current))
	}
	return segments
}

// buildNestedMap constructs a nested map[string]interface{} from a slice of key
// segments and a leaf value.
func buildNestedMap(segments []string, value interface{}) map[string]interface{} {
	if len(segments) == 1 {
		return map[string]interface{}{segments[0]: value}
	}
	return map[string]interface{}{segments[0]: buildNestedMap(segments[1:], value)}
}
