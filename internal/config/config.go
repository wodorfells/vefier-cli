package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Provider struct {
	Name    string `json:"name"`
	BaseURL string `json:"base_url"`
	Key     string `json:"key"`
	Model   string `json:"model"`
}

type Config struct {
	Language        string              `json:"language"`
	Theme           string              `json:"theme"`
	ActiveProvider  string              `json:"active"`
	Providers       map[string]Provider `json:"providers"`
	ConfigPath      string              `json:"-"`
}

func Load() *Config {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".vefier_settings.json")
	
	cfg := &Config{
		Language: "ru",
		Theme:    "crush",
		ActiveProvider: "openai",
		Providers: map[string]Provider{
			"openai": {Name: "OpenAI", BaseURL: "https://api.openai.com/v1", Model: "gpt-4o", Key: ""},
			"openrouter": {Name: "OpenRouter", BaseURL: "https://openrouter.ai/api/v1", Model: "anthropic/claude-3.5-sonnet", Key: ""},
		},
		ConfigPath: path,
	}

	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, cfg)
	} else {
		cfg.Save()
	}
	
	if len(cfg.Providers) == 0 {
		cfg.Providers = make(map[string]Provider)
	}
	return cfg
}

func (c *Config) Save() {
	data, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(c.ConfigPath, data, 0644)
}

func (c *Config) GetActive() Provider {
	if p, ok := c.Providers[c.ActiveProvider]; ok {
		return p
	}
	// Fallback
	for _, p := range c.Providers {
		return p
	}
	return Provider{BaseURL: "https://api.openai.com/v1", Model: "gpt-4o"}
}