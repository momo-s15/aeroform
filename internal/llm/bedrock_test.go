package llm

import (
	"testing"
	"time"
)

func TestBedrockClientAvailability(t *testing.T) {
	client := NewBedrockClient("us-east-1", "", time.Second)
	if !client.IsAvailable() {
		t.Fatalf("expected bedrock client to be available with a region")
	}
}
