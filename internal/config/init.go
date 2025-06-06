package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./config.yml"
	}

	if _, err := os.Stat(path); err == nil {
		if err := cleanenv.ReadConfig(path, cfg); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("configuration error: %v", err)
	}

	return cfg, nil
}
