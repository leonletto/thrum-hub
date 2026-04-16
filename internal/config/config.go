package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Root      string          `yaml:"root"`
	Embedding EmbeddingConfig `yaml:"embedding"`
}

type EmbeddingConfig struct {
	Backend        string `yaml:"backend"` // ollama | api
	Model          string `yaml:"model"`
	OllamaEndpoint string `yaml:"ollama_endpoint,omitempty"`
	APIBase        string `yaml:"api_base,omitempty"`
	APIKeyEnv      string `yaml:"api_key_env,omitempty"`
	APIModel       string `yaml:"api_model,omitempty"`
}

func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := cfg.applyDefaultsAndValidate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) applyDefaultsAndValidate() error {
	if c.Root == "" {
		return fmt.Errorf("config: root is required")
	}
	switch c.Embedding.Backend {
	case "ollama":
		if c.Embedding.OllamaEndpoint == "" {
			c.Embedding.OllamaEndpoint = "http://localhost:11434"
		}
		if c.Embedding.Model == "" {
			c.Embedding.Model = "nomic-embed-text"
		}
	case "api":
		if c.Embedding.APIBase == "" {
			return fmt.Errorf("config: embedding.api_base required for backend=api")
		}
		if c.Embedding.APIKeyEnv == "" {
			return fmt.Errorf("config: embedding.api_key_env required for backend=api")
		}
	default:
		return fmt.Errorf("config: embedding.backend must be ollama or api (got %q)", c.Embedding.Backend)
	}
	return nil
}
