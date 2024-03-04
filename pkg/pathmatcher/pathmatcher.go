package pathmatcher

type PathMatcher interface {
	Match(path string) bool
}
