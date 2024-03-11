package pathmatcher

import (
	"bytes"
	"fmt"
	"github.com/emilkje/cwc/pkg/errors"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type GitignorePathMatcher struct {
	ignoredPaths []string
}

func NewGitignorePathMatcher() (*GitignorePathMatcher, error) {
	matcher := &GitignorePathMatcher{
		ignoredPaths: make([]string, 0),
	}

	err := matcher.gitLsFiles()
	return matcher, err
}

func (g *GitignorePathMatcher) Match(path string) bool {
	return slices.Contains(g.ignoredPaths, path)
}

func (g *GitignorePathMatcher) Any() bool {
	return len(g.ignoredPaths) > 0
}

func (g *GitignorePathMatcher) gitLsFiles() error {
	// git ls-files -o --exclude-standard
	buf := new(bytes.Buffer)
	cmd := exec.Command("git", "ls-files", "-o", "--ignored", "--exclude-standard")
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found in") {
			return errors.GitNotInstalledError{Message: "git not found in PATH"}
		}

		return fmt.Errorf("error running git ls-files: %w", err)
	}

	// create a slice of ignored paths and remove the last empty string
	ignored := strings.Split(buf.String(), "\n")
	ignored = ignored[:len(ignored)-1]

	g.ignoredPaths = append(g.ignoredPaths, ignored...)

	return nil
}
