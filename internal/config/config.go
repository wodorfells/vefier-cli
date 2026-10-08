package config

import (
    "os"
    "path/filepath"
)

type Config struct {
	Theme          string
	DefaultModel   string
	ConfigFilePath string
}

func LoadConfig() *Config {
    home, _ := os.UserHomeDir()
	return &Config{
		Theme:          "crush",
		DefaultModel:   "gpt-4o",
        ConfigFilePath: filepath.Join(home, ".vefier_config.json"),
	}
}