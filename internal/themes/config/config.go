package theme_config

import (
	"context"
	"fmt"
	"os"
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

func GetThemeConfigFromURL(githubRepo GithubRepo) ThemeConfig {
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
		logrus.Fatal("Error downloading theme.yaml: ", err)
	}

	// Read theme.yaml
	config, err := fileContent.GetContent()
	if err != nil {
		logrus.Fatal("Error reading theme.yaml: ", err)
	}

	// Parse theme.yaml
	themeConfig := ThemeConfig{}
	err = yaml.Unmarshal([]byte(config), &themeConfig)
	if err != nil {
		logrus.Fatal("Error parsing theme.yaml: ", err)
	}

	// Normalize configuration keys to camelCase
	if themeConfig.Configuration != nil {
		themeConfig.Configuration = NormalizeConfigKeys(themeConfig.Configuration)
	}

	return themeConfig
}

func GetThemeConfigFromDir(dir string) ThemeConfig {
	// Read theme.yaml
	config, err := os.ReadFile(dir + "/theme.yaml")

	if err != nil {
		logrus.Panic("Error reading theme.yaml")
	}

	// Parse theme.yaml
	themeConfig := ThemeConfig{}
	err = yaml.Unmarshal(config, &themeConfig)
	if err != nil {
		logrus.Panic("Error parsing theme.yaml")
	}

	// Normalize configuration keys to camelCase
	if themeConfig.Configuration != nil {
		themeConfig.Configuration = NormalizeConfigKeys(themeConfig.Configuration)
	}

	return themeConfig
}

func GetThemeConfigFromThemeName(themeName string) ThemeConfig {
	// This function now expects themeName to already be resolved to the actual directory name
	// The caller should use GetThemeDirectory from themes package first
	return GetThemeConfigFromDir("themes/" + themeName)
}

func (tc *ThemeConfig) VerifyThemeMinimumVersion(cvwonderVersion string) bool {
	// Check if the minimum version is less than or equal to the current version
	if tc.MinimumVersion <= cvwonderVersion {
		return true
	}
	logrus.Error("CV Wonder version: ", cvwonderVersion)
	logrus.Error("Theme minimum version: ", tc.MinimumVersion)
	logrus.Error("The theme minimum version not met. You might encounter issues with this theme.")
	logrus.Error("")
	return false
}

func GetDefaultBranch(githubRepo GithubRepo) string {
	client := utils.GetGitHubClient()
	repo, _, err := client.Repositories.Get(context.TODO(), githubRepo.Owner, githubRepo.Name)
	if err != nil {
		logrus.Fatal("Error getting repository info: ", err)
	}
	return repo.GetDefaultBranch()
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
	merged := make(map[string]interface{})
	if base != nil {
		for k, v := range base {
			merged[k] = v
		}
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
