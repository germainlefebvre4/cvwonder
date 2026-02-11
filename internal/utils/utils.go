package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	dotignore "github.com/codeglyph/go-dotignore/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"
)

func CheckError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func ReadDirectory(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	CheckError(err)

	result := make([]string, 0)
	for _, file := range files {
		result = append(result, file.Name())
	}

	return result, nil
}

// IgnoreMatcher interface for matching ignore patterns
type IgnoreMatcher interface {
	Matches(path string) (bool, error)
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

// LoadIgnoreMatcher loads a .cvwonderignore file from the specified directory
// Returns nil if the file doesn't exist (no error)
func LoadIgnoreMatcher(rootDir string) (IgnoreMatcher, error) {
	ignorePath := filepath.Join(rootDir, ".cvwonderignore")
	
	// Check if .cvwonderignore exists
	if _, err := os.Stat(ignorePath); os.IsNotExist(err) {
		return nil, nil
	}
	
	// Load and parse the ignore file
	matcher, err := dotignore.NewPatternMatcherFromFile(ignorePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load .cvwonderignore: %w", err)
	}
	
	return matcher, nil
}

// CopyDirectory copies a directory tree from srcDir to dest, excluding files based on optional ignoreMatcher
// If ignoreMatcher is nil, only default exclusions (.git and index.html) apply
func CopyDirectory(scrDir string, dest string) error {
	return CopyDirectoryWithIgnore(scrDir, dest, nil)
}

// CopyDirectoryWithIgnore copies a directory tree from srcDir to dest with optional ignore patterns
// ignoreMatcher: optional matcher for .cvwonderignore patterns (nil to disable)
func CopyDirectoryWithIgnore(scrDir string, dest string, ignoreMatcher IgnoreMatcher) error {
	return copyDirectoryWithIgnoreInternal(scrDir, dest, ignoreMatcher, scrDir)
}

// copyDirectoryWithIgnoreInternal is the internal recursive implementation
// baseDir tracks the original source directory for calculating relative paths for pattern matching
func copyDirectoryWithIgnoreInternal(scrDir string, dest string, ignoreMatcher IgnoreMatcher, baseDir string) error {
	entries, err := os.ReadDir(scrDir)
	if err != nil {
		return err
	}
	// Cleanup the .git/ dir
	if ok, _ := regexp.MatchString(".git", scrDir); ok {
		return nil
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())
		
		// Calculate relative path from base directory for ignore matching
		relPath, err := filepath.Rel(baseDir, sourcePath)
		if err != nil {
			return err
		}

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			// For directories, we need to recurse even if ignored, to check for negation patterns
			// Check if directory itself is ignored
			var dirIgnored bool
			if ignoreMatcher != nil {
				dirIgnored, err = ignoreMatcher.Matches(relPath)
				if err != nil {
					return fmt.Errorf("error checking ignore pattern for %s: %w", relPath, err)
				}
				if dirIgnored {
					logrus.Debugf("Directory %s matched ignore pattern, but checking contents for negations", relPath)
				}
			}
			
			// Always create the directory and recurse to check for negation patterns
			if err := CreateIfNotExists(destPath, 0750); err != nil {
				return err
			}
			if err := copyDirectoryWithIgnoreInternal(sourcePath, destPath, ignoreMatcher, baseDir); err != nil {
				return err
			}
			
			// After recursion, check if the directory is empty (all contents were ignored)
			// If so, remove it and skip chmod
			entries, err := os.ReadDir(destPath)
			if err == nil && len(entries) == 0 {
				os.Remove(destPath)
				continue // Skip the chmod below
			}
		case os.ModeSymlink:
			// Check if symlink should be ignored
			if ignoreMatcher != nil {
				ignored, err := ignoreMatcher.Matches(relPath)
				if err != nil {
					return fmt.Errorf("error checking ignore pattern for %s: %w", relPath, err)
				}
				if ignored {
					logrus.Debugf("Ignoring %s (matched .cvwonderignore pattern)", relPath)
					continue
				}
			}
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			// Check if file should be ignored by .cvwonderignore
			if ignoreMatcher != nil {
				ignored, err := ignoreMatcher.Matches(relPath)
				if err != nil {
					return fmt.Errorf("error checking ignore pattern for %s: %w", relPath, err)
				}
				if ignored {
					logrus.Debugf("Ignoring %s (matched .cvwonderignore pattern)", relPath)
					continue
				}
			}
			
			// Exclude the templated index.html theme file and path with .git/ dir from copying
			if destPath == filepath.Join(dest, "index.html") {
				continue
			}
			if ok, _ := regexp.MatchString(".git", destPath); ok {
				continue
			}
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		fInfo, err := entry.Info()
		if err != nil {
			return err
		}

		isSymlink := fInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, fInfo.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(uint64(time.Now().UnixNano()))
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}
