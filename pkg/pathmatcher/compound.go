package pathmatcher

type CompoundPathMatcher struct {
	matchers []PathMatcher
}

func NewCompoundPathMatcher(matchers ...PathMatcher) *CompoundPathMatcher {
	return &CompoundPathMatcher{
		matchers: matchers,
	}
}

func (c *CompoundPathMatcher) Match(path string) bool {
	for _, matcher := range c.matchers {
		if matcher.Match(path) {
			return true
		}
	}

	return false
}

func (c *CompoundPathMatcher) Add(matcher PathMatcher) {
	c.matchers = append(c.matchers, matcher)
}
