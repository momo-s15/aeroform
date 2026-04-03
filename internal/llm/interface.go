package llm

import "time"

type Client interface {
	Generate(prompt string) (string, error)
	IsAvailable() bool
	Name() string
}

type Config struct {
	Backend string
	Host    string
	Model   string
	Region  string
	APIKey  string
	RoleARN string
	Timeout time.Duration
}

func NewClient(cfg Config) Client {
	if cfg.Backend == "" {
		cfg.Backend = "ollama"
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}

	switch cfg.Backend {
	case "ollama":
		return NewOllamaClient(cfg.Host, cfg.Model, cfg.Timeout)
	case "openai":
		return NewOpenAIClient(cfg.APIKey, cfg.Model, cfg.Timeout)
	case "bedrock":
		return NewBedrockClient(cfg.Region, cfg.Model, cfg.Timeout)
	default:
		return NewOllamaClient(cfg.Host, cfg.Model, cfg.Timeout)
	}
}
