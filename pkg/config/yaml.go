package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type YamlMarshaller struct{}

func (y *YamlMarshaller) Marshal(data interface{}) ([]byte, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling data: %w", err)
	}

	return bytes, nil
}

func (y *YamlMarshaller) Unmarshal(data []byte, out interface{}) error {
	err := yaml.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("error unmarshalling data: %w", err)
	}

	return nil
}
