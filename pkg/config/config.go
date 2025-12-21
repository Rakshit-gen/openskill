package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GroqAPIKey string `yaml:"groq_api_key"`
	Model      string `yaml:"model,omitempty"`
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".openskill"), nil
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func GetAPIKey() string {
	// Environment variable takes precedence
	if key := os.Getenv("GROQ_API_KEY"); key != "" {
		return key
	}

	// Fall back to config file
	cfg, err := Load()
	if err != nil {
		return ""
	}

	return cfg.GroqAPIKey
}

func GetModel() string {
	// Environment variable takes precedence
	if model := os.Getenv("OPENSKILL_MODEL"); model != "" {
		return model
	}

	// Fall back to config file
	cfg, err := Load()
	if err != nil || cfg.Model == "" {
		return "llama-3.3-70b-versatile"
	}

	return cfg.Model
}
