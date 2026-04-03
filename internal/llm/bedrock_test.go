package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type mockInvoker struct {
	response []byte
	err      error
}

func (m *mockInvoker) InvokeModel(_ context.Context, _ *bedrockruntime.InvokeModelInput, _ ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &bedrockruntime.InvokeModelOutput{Body: m.response}, nil
}

func TestBedrockClientName(t *testing.T) {
	c := NewBedrockClient("us-east-1", "", time.Second)
	if c.Name() != "bedrock" {
		t.Fatalf("expected 'bedrock', got %q", c.Name())
	}
}

func TestBedrockClientNotAvailableWithoutRegion(t *testing.T) {
	c := NewBedrockClient("", "", time.Second)
	if c.IsAvailable() {
		t.Fatal("expected not available without region")
	}
}

func TestBedrockDefaultModel(t *testing.T) {
	c := NewBedrockClient("us-east-1", "", time.Second)
	if c.model != "anthropic.claude-3-haiku-20240307-v1:0" {
		t.Fatalf("expected default Claude model, got %q", c.model)
	}
}

func TestBedrockCustomModel(t *testing.T) {
	c := NewBedrockClient("us-east-1", "amazon.titan-text-lite-v1", time.Second)
	if c.model != "amazon.titan-text-lite-v1" {
		t.Fatalf("expected custom model, got %q", c.model)
	}
}

func TestBedrockGenerateClaudeResponse(t *testing.T) {
	resp := claudeResponse{
		Content: []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{
			{Type: "text", Text: `{"templates":["vpc","eks"],"params":{}}`},
		},
		StopReason: "end_turn",
	}
	body, _ := json.Marshal(resp)

	c := NewBedrockClient("us-east-1", "anthropic.claude-3-haiku-20240307-v1:0", 5*time.Second)
	c.invoke = &mockInvoker{response: body}

	result, err := c.Generate("deploy kubernetes")
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if result != `{"templates":["vpc","eks"],"params":{}}` {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestBedrockGenerateTitanResponse(t *testing.T) {
	resp := titanResponse{
		Results: []struct {
			OutputText string `json:"outputText"`
		}{
			{OutputText: `{"templates":["vpc"],"params":{}}`},
		},
	}
	body, _ := json.Marshal(resp)

	c := NewBedrockClient("us-east-1", "amazon.titan-text-lite-v1", 5*time.Second)
	c.invoke = &mockInvoker{response: body}

	result, err := c.Generate("deploy storage")
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if result != `{"templates":["vpc"],"params":{}}` {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestBedrockGenerateInvokeError(t *testing.T) {
	c := NewBedrockClient("us-east-1", "", 5*time.Second)
	c.invoke = &mockInvoker{err: fmt.Errorf("access denied")}

	_, err := c.Generate("test")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBedrockGenerateEmptyClaudeContent(t *testing.T) {
	resp := claudeResponse{Content: nil}
	body, _ := json.Marshal(resp)

	c := NewBedrockClient("us-east-1", "", 5*time.Second)
	c.invoke = &mockInvoker{response: body}

	_, err := c.Generate("test")
	if err == nil {
		t.Fatal("expected error for empty content")
	}
}

func TestBedrockGenerateEmptyTitanResults(t *testing.T) {
	resp := titanResponse{Results: nil}
	body, _ := json.Marshal(resp)

	c := NewBedrockClient("us-east-1", "amazon.titan-text-lite-v1", 5*time.Second)
	c.invoke = &mockInvoker{response: body}

	_, err := c.Generate("test")
	if err == nil {
		t.Fatal("expected error for empty results")
	}
}

func TestBedrockGenerateNoRegion(t *testing.T) {
	c := NewBedrockClient("", "", 5*time.Second)
	_, err := c.Generate("test")
	if err == nil {
		t.Fatal("expected error without region")
	}
}

func TestBedrockBuildClaudeBody(t *testing.T) {
	c := NewBedrockClient("us-east-1", "anthropic.claude-3-haiku-20240307-v1:0", time.Second)
	body, ct, err := c.buildRequestBody("hello")
	if err != nil {
		t.Fatalf("buildRequestBody error: %v", err)
	}
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
	var req claudeRequest
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if req.AnthropicVersion != "bedrock-2023-05-31" {
		t.Fatalf("expected bedrock-2023-05-31, got %q", req.AnthropicVersion)
	}
	if len(req.Messages) != 1 || req.Messages[0].Content != "hello" {
		t.Fatal("expected single user message with 'hello'")
	}
}

func TestBedrockBuildTitanBody(t *testing.T) {
	c := NewBedrockClient("us-east-1", "amazon.titan-text-lite-v1", time.Second)
	body, _, err := c.buildRequestBody("hello")
	if err != nil {
		t.Fatalf("buildRequestBody error: %v", err)
	}
	var req titanRequest
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if req.InputText != "hello" {
		t.Fatalf("expected 'hello', got %q", req.InputText)
	}
}
