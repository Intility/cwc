package pathmatcher

import "regexp"

type regexPathMatcher struct {
	re *regexp.Regexp
}

func NewRegexPathMatcher(pattern string) (PathMatcher, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &regexPathMatcher{re: re}, nil
}

func (r *regexPathMatcher) Match(path string) bool {
	return r.re.MatchString(path)
}
