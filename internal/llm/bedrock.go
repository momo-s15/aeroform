package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockClient struct {
	region  string
	model   string
	roleARN string
	timeout time.Duration
	invoke  bedrockInvoker
}

type bedrockInvoker interface {
	InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error)
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
	if strings.TrimSpace(c.region) == "" {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.loadAWSConfig(ctx)
	return err == nil
}

func (c *BedrockClient) Generate(prompt string) (string, error) {
	if strings.TrimSpace(c.region) == "" {
		return "", fmt.Errorf("bedrock: aws region is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	invoker, err := c.getInvoker(ctx)
	if err != nil {
		return "", fmt.Errorf("bedrock: failed to create client: %w", err)
	}

	body, contentType, err := c.buildRequestBody(prompt)
	if err != nil {
		return "", fmt.Errorf("bedrock: failed to build request: %w", err)
	}

	out, err := invoker.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.model),
		ContentType: aws.String(contentType),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return "", fmt.Errorf("bedrock: InvokeModel failed: %w", err)
	}

	return c.parseResponse(out.Body)
}

func (c *BedrockClient) loadAWSConfig(ctx context.Context) (aws.Config, error) {
	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(c.region),
	}
	if c.roleARN != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("", "", ""),
		))
	}
	return awsconfig.LoadDefaultConfig(ctx, opts...)
}

func (c *BedrockClient) getInvoker(ctx context.Context) (bedrockInvoker, error) {
	if c.invoke != nil {
		return c.invoke, nil
	}
	cfg, err := c.loadAWSConfig(ctx)
	if err != nil {
		return nil, err
	}
	return bedrockruntime.NewFromConfig(cfg), nil
}

func (c *BedrockClient) buildRequestBody(prompt string) ([]byte, string, error) {
	if strings.Contains(c.model, "anthropic.claude") {
		return c.buildClaudeBody(prompt)
	}
	if strings.Contains(c.model, "amazon.titan") {
		return c.buildTitanBody(prompt)
	}
	return c.buildClaudeBody(prompt)
}

// Claude Messages API format
type claudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version"`
	MaxTokens        int             `json:"max_tokens"`
	Temperature      float64         `json:"temperature"`
	Messages         []claudeMessage `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
}

func (c *BedrockClient) buildClaudeBody(prompt string) ([]byte, string, error) {
	req := claudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        4096,
		Temperature:      0,
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
	}
	body, err := json.Marshal(req)
	return body, "application/json", err
}

// Amazon Titan format
type titanRequest struct {
	InputText            string         `json:"inputText"`
	TextGenerationConfig titanGenConfig `json:"textGenerationConfig"`
}

type titanGenConfig struct {
	MaxTokenCount int     `json:"maxTokenCount"`
	Temperature   float64 `json:"temperature"`
}

type titanResponse struct {
	Results []struct {
		OutputText string `json:"outputText"`
	} `json:"results"`
}

func (c *BedrockClient) buildTitanBody(prompt string) ([]byte, string, error) {
	req := titanRequest{
		InputText: prompt,
		TextGenerationConfig: titanGenConfig{
			MaxTokenCount: 4096,
			Temperature:   0,
		},
	}
	body, err := json.Marshal(req)
	return body, "application/json", err
}

func (c *BedrockClient) parseResponse(body []byte) (string, error) {
	if strings.Contains(c.model, "amazon.titan") {
		var resp titanResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return "", fmt.Errorf("bedrock: failed to parse Titan response: %w", err)
		}
		if len(resp.Results) == 0 {
			return "", fmt.Errorf("bedrock: Titan returned no results")
		}
		return resp.Results[0].OutputText, nil
	}

	var resp claudeResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("bedrock: failed to parse Claude response: %w", err)
	}
	var parts []string
	for _, block := range resp.Content {
		if block.Type == "text" {
			parts = append(parts, block.Text)
		}
	}
	if len(parts) == 0 {
		return "", fmt.Errorf("bedrock: Claude returned no text content")
	}
	return strings.Join(parts, ""), nil
}
