package pathmatcher

import (
	"fmt"
	"regexp"
)

type RegexPathMatcher struct {
	re *regexp.Regexp
}

func NewRegexPathMatcher(pattern string) (*RegexPathMatcher, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex pattern: %w", err)
	}

	return &RegexPathMatcher{re: re}, nil
}

func (r *RegexPathMatcher) Match(path string) bool {
	return r.re.MatchString(path)
}
