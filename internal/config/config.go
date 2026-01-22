package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

const defaultPort = 8423

func DefaultConfig() Config {
	return Config{
		Host: "localhost",
		Port: defaultPort,
	}
}

func (c Config) BaseURL() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "itsyhome", "config.json")
}

func Load() Config {
	path := configPath()
	if path == "" {
		return DefaultConfig()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig()
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig()
	}

	if cfg.Port == 0 {
		cfg.Port = defaultPort
	}
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}

	return cfg
}

func Save(cfg Config) error {
	path := configPath()
	if path == "" {
		return fmt.Errorf("cannot determine config path")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func Path() string {
	return configPath()
}
