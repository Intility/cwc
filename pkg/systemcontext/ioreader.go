package systemcontext

import (
	"fmt"
	"io"
)

type IOReaderContextRetriever struct {
	io.Reader
}

func NewIOReaderContextRetriever(reader io.Reader) *IOReaderContextRetriever {
	return &IOReaderContextRetriever{
		Reader: reader,
	}
}

func (r *IOReaderContextRetriever) RetrieveContext() (string, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("error reading from io.Reader: %w", err)
	}

	return string(bytes), nil
}
