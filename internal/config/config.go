package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/georgg2003/shortener/pkg/flagstruct"
	"github.com/georgg2003/shortener/pkg/utils"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert/yaml"
)

var (
	errUnknownFormat     = errors.New("unknown format")
	errFileWithoutFormat = errors.New("config file doesn't have any format")
)

type Config struct {
	ConfigPath       string `env:"CONFIG" flag:"config" flag_usage:"path to config file"`
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

// Читает конфиг из файла.
func (c *Config) ReadFromFile(fname string) (err error) {
	if fname == "" {
		return nil
	}

	parts := strings.Split(fname, ".")
	lp := len(parts)
	if lp < 2 {
		return errFileWithoutFormat
	}
	format := parts[lp-1]

	if format != "yaml" && format != "json" {
		return errUnknownFormat
	}

	f, err := os.ReadFile(fname)
	if err != nil {
		return utils.ErrWrap(err, fmt.Sprintf("failed to read config %s", fname))
	}

	var raw any
	switch format {
	case "yaml":
		err = yaml.Unmarshal(f, &raw)
	case "json":
		err = json.Unmarshal(f, &raw)
	default:
		err = errUnknownFormat
	}

	if err != nil {
		return utils.ErrWrap(err, "failed to unmarshall data")
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

func New() *Config {
	return &Config{
		ListenAddr:      "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
		DataBaseDSN:     "",
		JWTSecretKey:    "secret_key",
	}
}

func Load() (*Config, error) {
	cfg := New()

	// Предварительная инициализация, чтобы достать путь до конфига.
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	_ = cfg.ReadFromEnv()
	_ = cfg.ReadFromFlags(fs)
	_ = fs.Parse(os.Args[1:])

	if cfg.ConfigPath != "" {
		if err := cfg.ReadFromFile(cfg.ConfigPath); err != nil {
			return nil, fmt.Errorf("file config error: %w", err)
		}
	}

	if err := cfg.ReadFromEnv(); err != nil {
		return nil, fmt.Errorf("env config error: %w", err)
	}

	if err := cfg.ReadFromFlags(flag.CommandLine); err != nil {
		return nil, fmt.Errorf("flag config error: %w", err)
	}

	if !flag.Parsed() {
		flag.Parse()
	}

	return cfg, nil
}
