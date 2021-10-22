package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-yaml"
)

const DefaultContextTimeout = 30 * time.Second

const (
	hueUriEnv           string = "HUE_URI"
	hueTokenEnv         string = "HUE_TOKEN"
	magicHeaderKeyEnv   string = "MAGIC_HEADER_KEY"
	magicHeaderValueEnv string = "MAGIC_HEADER_VALUE"
)

type Config struct {
	URI         string             `json:"uri" yaml:"uri"`
	Token       string             `json:"token" yaml:"token"`
	MagicHeader *MagicHeaderConfig `json:"magic_header,omitempty" yaml:"magic_header,omitempty"`
}

type MagicHeaderConfig struct {
	Key    string `json:"key" yaml:"key"`
	Value  string `json:"value" yaml:"value"`
	Header string `json:"header" yaml:"header"`
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

func LoadConfigEnv() (c *Config, err error) {
	c = new(Config)
	c.MagicHeader = new(MagicHeaderConfig)

	c.URI, err = getEnv(hueUriEnv)
	if err != nil {
		c = nil
		return
	}

	c.Token, err = getEnv(hueTokenEnv)
	if err != nil {
		c = nil
		return
	}

	c.MagicHeader.Key, err = getEnv(magicHeaderKeyEnv)
	if err != nil {
		c = nil
		return
	}

	c.MagicHeader.Value, err = getEnv(magicHeaderValueEnv)
	if err != nil {
		c = nil
		return
	}

	return
}

func getEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("%s not set in environment", key)
	}

	return val, nil
}