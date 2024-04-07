package tools

type WebExecutor struct {
	urls []string
}

func NewWebExecutor(urls []string) *WebExecutor {
	return &WebExecutor{urls: urls}
}

func (s *WebExecutor) Execute(arguments string) (string, error) {
	return "Web tools are not yet supported", nil
}
