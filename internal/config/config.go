package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Cloud string      `yaml:"cloud"`
	Mode  string      `yaml:"mode"`
	AWS   AWSConfig   `yaml:"aws"`
	Azure AzureConfig `yaml:"azure"`
	GCP   GCPConfig   `yaml:"gcp"`
	LLM   LLMConfig   `yaml:"llm"`
	State StateConfig `yaml:"state"`
}

type AWSConfig struct {
	Region    string `yaml:"region"`
	AccountID string `yaml:"account_id"`
}

type AzureConfig struct {
	Location string `yaml:"location"`
}

type GCPConfig struct {
	Region string `yaml:"region"`
}

type StateConfig struct {
	Backend string `yaml:"backend"`
}

type LLMConfig struct {
	Backend string `yaml:"backend"`
	Model   string `yaml:"model"`
	Host    string `yaml:"host"`
	Region  string `yaml:"region"`
	APIKey  string `yaml:"api_key"`
	RoleARN string `yaml:"role_arn"`
}

func LoadFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if cfg.Mode == "" {
		cfg.Mode = "pro"
	}

	if err := Validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
