package pathmatcher

type compoundPathMatcher struct {
	matchers []PathMatcher
}

func NewCompoundPathMatcher(matchers ...PathMatcher) PathMatcher {
	return &compoundPathMatcher{
		matchers: matchers,
	}
}

func (c *compoundPathMatcher) Match(path string) bool {
	for _, matcher := range c.matchers {
		if matcher.Match(path) {
			return true
		}
	}
	return false
}

func (c *compoundPathMatcher) Add(matcher PathMatcher) {
	c.matchers = append(c.matchers, matcher)
}
