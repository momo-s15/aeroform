package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "config.yaml")
	content := []byte("cloud: aws\nmode: pro\naws:\n  region: us-east-1\nstate:\n  backend: s3\nllm:\n  backend: ollama\n  model: llama3.2\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	if cfg.Cloud != "aws" || cfg.AWS.Region != "us-east-1" {
		t.Fatalf("unexpected config %#v", cfg)
	}
	if cfg.LLM.Backend != "ollama" {
		t.Fatalf("expected llm backend to load")
	}
}
