package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGenerateTemplateSelection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{
			"response": `{"templates":["static-site"],"params":{"project_name":"demo"}}`,
		})
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "llama3.2", time.Second)
	if !client.IsAvailable() {
		t.Fatalf("expected client to be available")
	}
	raw, err := client.Generate("hello")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	selection, err := ParseTemplateSelection(raw)
	if err != nil {
		t.Fatalf("ParseTemplateSelection() error = %v", err)
	}
	if len(selection.Templates) != 1 || selection.Templates[0] != "static-site" {
		t.Fatalf("unexpected selection %#v", selection)
	}
}

func TestGenerateTemplateSelectionRejectsUnknownTemplate(t *testing.T) {
	selection, err := ParseTemplateSelection(`{"templates":["unknown"],"params":{}}`)
	if err != nil {
		t.Fatalf("ParseTemplateSelection() error = %v", err)
	}
	if len(selection.Templates) != 1 {
		t.Fatalf("unexpected templates")
	}
}
