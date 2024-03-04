package pathmatcher

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/emilkje/cwc/pkg/errors"
)

type gitignorePathMatcher struct {
	patterns []string
}

func NewGitignorePathMatcher(gitignorePath string) (PathMatcher, error) {
	patterns, err := parseGitignore(gitignorePath)

	if err != nil {
		return nil, err
	}

	return &gitignorePathMatcher{patterns: patterns}, nil
}

// Match accepts a relative path to a file (e.g. "foo/bar.txt") and returns true if the path is matched by any of the patterns in the .gitignore file.
func (g *gitignorePathMatcher) Match(path string) bool {
	for _, pattern := range g.patterns {
		if match(pattern, path) {
			return true
		}
	}
	return false
}

func parseGitignore(gitignorePath string) (patterns []string, err error) {

	// no .gitignore file is a valid state
	_, err = os.Stat(gitignorePath)
	if os.IsNotExist(err) {
		return patterns, errors.FileNotExistError{FileName: gitignorePath}
	}

	data, err := os.ReadFile(gitignorePath) // #nosec
	if err != nil {
		return nil, fmt.Errorf("error reading .gitignore: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		patterns = append(patterns, line)
	}
	return patterns, nil
}

const dblAsterisks = "**"

// Match matches patterns in the same manner that gitignore does.
// Reference https://git-scm.com/docs/gitignore.
// example: pattern is [Bb]uild/ and value is Build/foo.js
func match(pattern, value string) bool {

	// empty lines should be ignored
	if strings.TrimSpace(pattern) == "" {
		return false
	}

	// comments should be ignored
	if strings.HasPrefix(strings.TrimSpace(pattern), "#") {
		return false
	}

	// handle negation
	if pattern[0] == '!' {
		return !match(pattern[1:], value)
	}

	if pattern == "bin/" && value == "bin/cwc" {
		fmt.Println("hello")
	}

	// Placeholder sequences
	const (
		leadingDirsPlaceholder  = "__LEADING_DIR__"
		anyDirPlaceholder       = "__ANY_DIR__"
		trailingDirsPlaceholder = "__TRAILING_DIRS__"
	)

	// Handle special "**" cases
	if strings.HasPrefix(pattern, dblAsterisks+"/") {
		pattern = strings.Replace(pattern, dblAsterisks+"/", leadingDirsPlaceholder, 1)
	}

	pattern = strings.Replace(pattern, "/"+dblAsterisks+"/", anyDirPlaceholder, -1)
	if strings.HasSuffix(pattern, "/"+dblAsterisks) {
		pattern = strings.TrimSuffix(pattern, "/"+dblAsterisks) + trailingDirsPlaceholder
	}

	pattern = strings.ReplaceAll(pattern, "?", ".")
	pattern = strings.ReplaceAll(pattern, "/", "\\/")
	pattern = strings.ReplaceAll(pattern, ".", "\\.")
	pattern = strings.ReplaceAll(pattern, "*", "[^/]*")

	// replace placeholders
	pattern = strings.ReplaceAll(pattern, leadingDirsPlaceholder, ".*")
	pattern = strings.ReplaceAll(pattern, anyDirPlaceholder, ".*")
	pattern = strings.ReplaceAll(pattern, trailingDirsPlaceholder, ".*")

	// compile regext and match against path
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false
	}

	// replace

	return matched
}
