package config

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddr string `mapstructure:"listen_addr"`
	BaseURL    string `mapstructure:"base_url"`
}

func ReadFromYaml(filename string) (*Config, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c Config
	var raw interface{}

	if err := yaml.Unmarshal(f, &raw); err != nil {
		return nil, err
	}

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &c})
	if err := decoder.Decode(raw); err != nil {
		return nil, err
	}

	return &c, nil
}
