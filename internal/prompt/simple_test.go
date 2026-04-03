package prompt

import (
	"strings"
	"testing"
)

func TestBuildSimplePrompt(t *testing.T) {
	prompt := BuildSimplePrompt("a personal website", "aws", "us-east-1", []string{"static-site", "lambda-api"})
	if !strings.Contains(prompt, "static-site") || !strings.Contains(prompt, "lambda-api") {
		t.Fatalf("prompt missing allowed templates")
	}
	if !strings.Contains(prompt, "Respond ONLY with JSON") {
		t.Fatalf("prompt missing JSON instruction")
	}
}
