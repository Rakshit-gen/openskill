package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// Provider selection
	Provider string `yaml:"provider,omitempty"` // groq, openai, anthropic, ollama

	// Legacy Groq key (for backwards compatibility)
	GroqAPIKey string `yaml:"groq_api_key,omitempty"`

	// Provider-specific API keys
	OpenAIAPIKey    string `yaml:"openai_api_key,omitempty"`
	AnthropicAPIKey string `yaml:"anthropic_api_key,omitempty"`

	// Model settings (per provider)
	Model          string `yaml:"model,omitempty"`           // Default model
	GroqModel      string `yaml:"groq_model,omitempty"`      // Groq-specific model
	OpenAIModel    string `yaml:"openai_model,omitempty"`    // OpenAI-specific model
	AnthropicModel string `yaml:"anthropic_model,omitempty"` // Anthropic-specific model
	OllamaModel    string `yaml:"ollama_model,omitempty"`    // Ollama-specific model

	// Ollama settings
	OllamaEndpoint string `yaml:"ollama_endpoint,omitempty"` // Custom Ollama endpoint
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

// GetProvider returns the active provider
func GetProvider() string {
	// Environment variable takes precedence
	if provider := os.Getenv("OPENSKILL_PROVIDER"); provider != "" {
		return strings.ToLower(provider)
	}

	cfg, err := Load()
	if err != nil || cfg.Provider == "" {
		return "groq" // Default to groq for backwards compatibility
	}

	return strings.ToLower(cfg.Provider)
}

// GetProviderAPIKey returns the API key for a specific provider
func GetProviderAPIKey(provider string) string {
	provider = strings.ToLower(provider)

	// Check environment variables first
	switch provider {
	case "groq":
		if key := os.Getenv("GROQ_API_KEY"); key != "" {
			return key
		}
	case "openai":
		if key := os.Getenv("OPENAI_API_KEY"); key != "" {
			return key
		}
	case "anthropic":
		if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
			return key
		}
	}

	// Fall back to config file
	cfg, err := Load()
	if err != nil {
		return ""
	}

	switch provider {
	case "groq":
		return cfg.GroqAPIKey
	case "openai":
		return cfg.OpenAIAPIKey
	case "anthropic":
		return cfg.AnthropicAPIKey
	}

	return ""
}

// GetProviderModel returns the model for a specific provider
func GetProviderModel(provider string) string {
	provider = strings.ToLower(provider)

	// Check environment variable first
	if model := os.Getenv("OPENSKILL_MODEL"); model != "" {
		return model
	}

	cfg, err := Load()
	if err != nil {
		return getDefaultModel(provider)
	}

	// Check provider-specific model first
	switch provider {
	case "groq":
		if cfg.GroqModel != "" {
			return cfg.GroqModel
		}
	case "openai":
		if cfg.OpenAIModel != "" {
			return cfg.OpenAIModel
		}
	case "anthropic":
		if cfg.AnthropicModel != "" {
			return cfg.AnthropicModel
		}
	case "ollama":
		if cfg.OllamaModel != "" {
			return cfg.OllamaModel
		}
	}

	// Fall back to generic model setting
	if cfg.Model != "" {
		return cfg.Model
	}

	return getDefaultModel(provider)
}

func getDefaultModel(provider string) string {
	switch provider {
	case "groq":
		return "llama-3.3-70b-versatile"
	case "openai":
		return "gpt-4o-mini"
	case "anthropic":
		return "claude-3-5-sonnet-20241022"
	case "ollama":
		return "llama3.2"
	default:
		return "llama-3.3-70b-versatile"
	}
}

// GetOllamaEndpoint returns the Ollama endpoint
func GetOllamaEndpoint() string {
	if endpoint := os.Getenv("OLLAMA_HOST"); endpoint != "" {
		return endpoint + "/api/chat"
	}

	cfg, err := Load()
	if err != nil || cfg.OllamaEndpoint == "" {
		return "http://localhost:11434/api/chat"
	}

	return cfg.OllamaEndpoint
}

// Legacy functions for backwards compatibility
func GetAPIKey() string {
	return GetProviderAPIKey(GetProvider())
}

func GetModel() string {
	return GetProviderModel(GetProvider())
}
