package config

import (
	"flag"
	"os"

	"github.com/georgg2003/shortener/pkg/flagstruct"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddr       string `mapstructure:"listen_addr" env:"SERVER_ADDRESS" flag:"a" flag_usage:"listen addres"`
	BaseURL          string `mapstructure:"base_url" env:"BASE_URL" flag:"b" flag_usage:"base_url"`
	FileStoragePath  string `mapstructure:"file_storage_path" env:"FILE_STORAGE_PATH" flag:"f" flag_usage:"file storage path"`
	DataBaseDSN      string `mapstructure:"database_dsn" env:"DATABASE_DSN" flag:"d" flag_usage:"database dsn"`
	JWTSecretKey     string `mapstructure:"jwt_secret_key" env:"JWT_SECRET_KEY" flag:"k" flag_usage:"jwt secret key"`
	AuditFile        string `mapstructure:"audit_file" env:"AUDIT_FILE" flag:"audit-file" flag_usage:"file to write audit logs"`
	AuditURL         string `mapstructure:"audit_url" env:"AUDIT_URL" flag:"audit-url" flag_usage:"url to send audit logs"`
	AuditStubEnabled bool   `mapstructure:"audit_stub_enabled" env:"AUDIT_STUB_ENABLED"`
	DebugAddr        string `mapstructure:"debug_addr" env:"DEBUG_ADDR" flag:"debug-addr" flag_usage:"debug listen addres"`
	EnableHTTPS      bool   `mapstructure:"enable_https" env:"ENABLE_HTTPS" flag:"s" flag_usage:"enabled https"`
}

func New() *Config {
	return &Config{
		ListenAddr:      "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "./data.json",
		DataBaseDSN:     "",
		JWTSecretKey:    "secret_key",
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

	if err = yaml.Unmarshal(f, &raw); err != nil {
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

func (c *Config) ReadFromFlags(fs *flag.FlagSet) error {
	return flagstruct.ReadFromFlags(fs, c)
}

func (c *Config) ReadFromEnv() error {
	return cleanenv.ReadEnv(c)
}
