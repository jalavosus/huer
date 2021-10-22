package config

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/goccy/go-yaml"
)

const DefaultContextTimeout = 30 * time.Second

type Config struct {
	URI   string `json:"uri" yaml:"uri"`
	Token string `json:"token" yaml:"token"`
}

func LoadConfig(path string) (*Config, error) {
	var conf Config

	ext := filepath.Ext(path)

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if ext == ".json" {
		raw, err = yaml.JSONToYAML(raw)
		if err != nil {
			return nil, err
		}
	}

	if err = yaml.Unmarshal(raw, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}