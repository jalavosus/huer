package utils

import (
	"github.com/goccy/go-yaml"
)

type (
	MarshalFunc   func(any) ([]byte, error)
	UnmarshalFunc func([]byte, any) error
)

func MarshalYaml(data any) ([]byte, error) {
	return yaml.Marshal((interface{})(data))
}

func UnmarshalYAML(data []byte, out any) error {
	return yaml.Unmarshal(data, (interface{})(out))
}