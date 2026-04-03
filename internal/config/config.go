package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Cloud string      `yaml:"cloud" mapstructure:"cloud" validate:"required,oneof=aws azure gcp"`
	Mode  string      `yaml:"mode"  mapstructure:"mode"  validate:"omitempty,oneof=simple pro"`
	AWS   AWSConfig   `yaml:"aws"   mapstructure:"aws"`
	Azure AzureConfig `yaml:"azure" mapstructure:"azure"`
	GCP   GCPConfig   `yaml:"gcp"   mapstructure:"gcp"`
	LLM   LLMConfig   `yaml:"llm"   mapstructure:"llm"`
	State StateConfig `yaml:"state" mapstructure:"state"`
}

type AWSConfig struct {
	Region    string `yaml:"region"     mapstructure:"region"     validate:"required_if=Cloud aws"`
	AccountID string `yaml:"account_id" mapstructure:"account_id"`
}

type AzureConfig struct {
	Location string `yaml:"location" mapstructure:"location" validate:"required_if=Cloud azure"`
}

type GCPConfig struct {
	Region string `yaml:"region" mapstructure:"region" validate:"required_if=Cloud gcp"`
}

type StateConfig struct {
	Backend string `yaml:"backend" mapstructure:"backend"`
}

type LLMConfig struct {
	Backend string `yaml:"backend"  mapstructure:"backend"  validate:"omitempty,oneof=ollama openai bedrock"`
	Model   string `yaml:"model"    mapstructure:"model"`
	Host    string `yaml:"host"     mapstructure:"host"`
	Region  string `yaml:"region"   mapstructure:"region"`
	APIKey  string `yaml:"api_key"  mapstructure:"api_key"`
	RoleARN string `yaml:"role_arn" mapstructure:"role_arn"`
}

// LoadFromFile reads config.yaml and applies ENV var overrides.
// Priority: ENV vars (AEROFORM_*) > YAML values > defaults.
func LoadFromFile(path string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.SetEnvPrefix("AEROFORM")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	bindEnvKeys(v)

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
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

// bindEnvKeys explicitly binds dotted config keys to AEROFORM_ env vars
// so that nested values like aws.region are overridden by AEROFORM_AWS_REGION.
func bindEnvKeys(v *viper.Viper) {
	keys := []string{
		"cloud",
		"mode",
		"aws.region",
		"aws.account_id",
		"azure.location",
		"gcp.region",
		"llm.backend",
		"llm.model",
		"llm.host",
		"llm.region",
		"llm.api_key",
		"llm.role_arn",
		"state.backend",
	}
	for _, k := range keys {
		_ = v.BindEnv(k)
	}
}
