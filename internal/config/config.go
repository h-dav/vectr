package config

import (
	"log/slog"
	"path/filepath"

	"github.com/h-dav/envconfig"
)

type Config struct {
	Service struct {
		Version string `env:"VERSION"`
	} `env:"prefix=SERVICE_"`
	MinimumLogLevel string `env:"MINIMUM_LOG_LEVEL"`
	HTTP            struct {
		Port string `env:"PORT"`
	} `env:"prefix=HTTP_"`
	Database struct {
		Host     string `env:"HOST"`
		Name     string `env:"NAME"`
		Username string `env:"USERNAME"`
		Password string `env:"PASSWORD"`
		Port     string `env:"PORT"`
	} `env:"prefix=DATABASE_"`
}

func New() (Config, error) {
	var cfg Config
	if err := envconfig.SetPopulate(filepath.Join("internal", "config", "default.env"), &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func AsSlogLevel(level string) slog.Level {
	switch {
	case level == "DEBUG":
		return slog.LevelDebug
	case level == "INFO":
		return slog.LevelInfo
	case level == "ERROR":
		return slog.LevelError
	case level == "WARN":
		return slog.LevelWarn
	}
	return slog.LevelInfo
}
