package gitignore

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ParseGitignore(gitignorePath string) ([]*regexp.Regexp, error) {
	var patterns []*regexp.Regexp

	// no .gitignore file is a valid state
	_, err := os.Stat(gitignorePath)
	if os.IsNotExist(err) {
		return patterns, nil
	}

	// add the gitignore itself to the patterns
	patterns = append(patterns, regexp.MustCompile(globToRegexp(".gitignore")))

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
		pattern, err := regexp.Compile(globToRegexp(line))
		if err != nil {
			return nil, fmt.Errorf("invalid pattern in .gitignore: %s, error: %w", line, err)
		}
		patterns = append(patterns, pattern)
	}
	return patterns, nil
}

func globToRegexp(pattern string) string {
	// Convert simple shell glob patterns to regular expressions for matching
	esc := regexp.QuoteMeta(pattern)
	esc = strings.Replace(esc, "\\*", ".*", -1) // Convert * to .*
	esc = strings.Replace(esc, "\\?", ".", -1)  // Convert ? to .
	return esc
}
