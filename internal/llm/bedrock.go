package llm

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type BedrockClient struct {
	region  string
	model   string
	timeout time.Duration
}

func NewBedrockClient(region, model string, timeout time.Duration) *BedrockClient {
	if model == "" {
		model = "anthropic.claude-3-haiku-20240307-v1:0"
	}
	return &BedrockClient{region: region, model: model, timeout: timeout}
}

func (c *BedrockClient) Name() string {
	return "bedrock"
}

func (c *BedrockClient) IsAvailable() bool {
	return strings.TrimSpace(c.region) != ""
}

func (c *BedrockClient) Generate(prompt string) (string, error) {
	if !c.IsAvailable() {
		return "", errors.New("AWS region is required for Bedrock")
	}
	return "", fmt.Errorf("bedrock generation is scaffolded and requires aws-sdk-go-v2 bedrock runtime wiring in a future patch")
}
