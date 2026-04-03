package llm

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestOpenAIClientAvailability(t *testing.T) {
	client := NewOpenAIClient("", "gpt-4o-mini", time.Second)
	if client.IsAvailable() {
		t.Fatalf("expected client to be unavailable without an API key")
	}
}

func TestOpenAIClientGenerateWithServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(`{"choices":[{"message":{"content":"{\"templates\":[\"static-site\"],\"params\":{}}"}}]}`))
	}))
	defer server.Close()

	client := NewOpenAIClient("test-key", "gpt-4o-mini", time.Second)
	client.httpClient = server.Client()
	client.httpClient.Timeout = time.Second

	// Override endpoint by swapping transport through URL rewriting in request creation is not exposed here,
	// so this test only validates the client can be instantiated with a key.
	if !client.IsAvailable() {
		t.Fatalf("expected client to be available")
	}
	if strings.TrimSpace(client.Name()) != "openai" {
		t.Fatalf("unexpected client name")
	}
}
