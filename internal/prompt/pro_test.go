package prompt

import (
	"strings"
	"testing"
)

func TestBuildProPrompt(t *testing.T) {
	text := BuildProPrompt("resilient eks cluster", "aws", "us-east-1", "prod", []string{"vpc", "eks", "rds-private"})
	if !strings.Contains(text, "Available pro templates") {
		t.Fatalf("expected template section")
	}
	if !strings.Contains(text, "Respond ONLY with JSON") {
		t.Fatalf("expected JSON-only response instruction")
	}
	if !strings.Contains(text, "cloud=aws") {
		t.Fatalf("expected cloud context")
	}
}
