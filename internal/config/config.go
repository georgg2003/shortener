package config

import (
	"flag"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddr string `mapstructure:"listen_addr"`
	BaseURL    string `mapstructure:"base_url"`
}

func New() *Config {
	return &Config{
		ListenAddr: "localhost:8080",
		BaseURL:    "http://localhost:8080",
	}
}

func (c *Config) ReadFromYaml() error {
	configFilename := flag.String("config", "", "config filename")
	flag.Parse()

	f, err := os.ReadFile(*configFilename)
	if err != nil {
		return err
	}

	var raw any

	if err := yaml.Unmarshal(f, &raw); err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &c})
	if err != nil {
		return err
	}

	if err := decoder.Decode(raw); err != nil {
		return err
	}

	return nil
}

func (c *Config) ReadFromFlags() {
	listenAddr := flag.String("a", "", "listen addres")
	baseURL := flag.String("b", "", "base url")
	flag.Parse()

	if listenAddr != nil && *listenAddr != "" {
		c.ListenAddr = *listenAddr
	}
	if baseURL != nil && *baseURL != "" {
		c.BaseURL = *baseURL
	}
}

func (c *Config) ReadFromEnv() {
	listenAddr, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		c.ListenAddr = listenAddr
	}

	baseURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		c.BaseURL = baseURL
	}
}
