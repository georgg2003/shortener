package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddr      string `mapstructure:"listen_addr" env:"SERVER_ADDRESS"`
	BaseURL         string `mapstructure:"base_url" env:"BASE_URL"`
	FileStoragePath string `mapstructure:"file_storage_path" env:"FILE_STORAGE_PATH"`
	DataBaseDSN     string `mapstructure:"database_dsn" env:"DATABASE_DSN"`
}

func New() *Config {
	return &Config{
		ListenAddr:      "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "./data.json",
		DataBaseDSN:     "user=shortener password=password host=127.0.0.1 database=shortener port=5432", // "postgresql://shortener:password@127.0.0.1:5432/shortener",
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
	fileStoragePath := flag.String("f", "", "file storage path")
	dataBaseDSN := flag.String("d", "", "database dsn")
	flag.Parse()

	if listenAddr != nil && *listenAddr != "" {
		c.ListenAddr = *listenAddr
	}
	if baseURL != nil && *baseURL != "" {
		c.BaseURL = *baseURL
	}
	if fileStoragePath != nil && *fileStoragePath != "" {
		c.FileStoragePath = *fileStoragePath
	}
	if dataBaseDSN != nil && *fileStoragePath != "" {
		c.DataBaseDSN = *dataBaseDSN
	}
}

func (c *Config) ReadFromEnv() error {
	return cleanenv.ReadEnv(c)
}
