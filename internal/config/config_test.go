package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const baseYAML = "cloud: aws\nmode: pro\naws:\n  region: us-east-1\nstate:\n  backend: s3\nllm:\n  backend: ollama\n  model: llama3.2\n"

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestLoadFromFile(t *testing.T) {
	path := writeConfig(t, baseYAML)

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

func TestLoadFromFileDefaultMode(t *testing.T) {
	yaml := "cloud: aws\naws:\n  region: us-east-1\n"
	path := writeConfig(t, yaml)

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	if cfg.Mode != "pro" {
		t.Fatalf("expected default mode 'pro', got %q", cfg.Mode)
	}
}

func TestEnvOverridesCloud(t *testing.T) {
	yaml := "cloud: aws\naws:\n  region: us-east-1\ngcp:\n  region: us-central1\n"
	path := writeConfig(t, yaml)

	t.Setenv("AEROFORM_CLOUD", "gcp")

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	if cfg.Cloud != "gcp" {
		t.Fatalf("expected cloud=gcp from env, got %q", cfg.Cloud)
	}
}

func TestEnvOverridesNestedKey(t *testing.T) {
	path := writeConfig(t, baseYAML)

	t.Setenv("AEROFORM_AWS_REGION", "eu-west-1")

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	if cfg.AWS.Region != "eu-west-1" {
		t.Fatalf("expected aws.region=eu-west-1 from env, got %q", cfg.AWS.Region)
	}
}

func TestEnvOverridesLLMModel(t *testing.T) {
	path := writeConfig(t, baseYAML)

	t.Setenv("AEROFORM_LLM_MODEL", "codellama")

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	if cfg.LLM.Model != "codellama" {
		t.Fatalf("expected llm.model=codellama from env, got %q", cfg.LLM.Model)
	}
}

func TestValidationStillRuns(t *testing.T) {
	yaml := "cloud: invalid\naws:\n  region: us-east-1\n"
	path := writeConfig(t, yaml)

	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected validation error for unsupported cloud")
	}
}

func TestMissingFileReturnsError(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidateCloudRequired(t *testing.T) {
	err := Validate(Config{})
	if err == nil {
		t.Fatal("expected error for empty cloud")
	}
	if !strings.Contains(err.Error(), "cloud") {
		t.Fatalf("expected cloud error, got: %v", err)
	}
}

func TestValidateCloudOneof(t *testing.T) {
	err := Validate(Config{Cloud: "digitalocean"})
	if err == nil {
		t.Fatal("expected error for invalid cloud")
	}
	if !strings.Contains(err.Error(), "must be one of") {
		t.Fatalf("expected oneof error, got: %v", err)
	}
}

func TestValidateModeOneof(t *testing.T) {
	err := Validate(Config{Cloud: "aws", Mode: "turbo", AWS: AWSConfig{Region: "us-east-1"}})
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
	if !strings.Contains(err.Error(), "must be one of") {
		t.Fatalf("expected oneof error, got: %v", err)
	}
}

func TestValidateAWSRegionRequired(t *testing.T) {
	err := Validate(Config{Cloud: "aws"})
	if err == nil {
		t.Fatal("expected error for missing aws.region")
	}
	if !strings.Contains(err.Error(), "aws.region") {
		t.Fatalf("expected aws.region error, got: %v", err)
	}
}

func TestValidateAzureLocationRequired(t *testing.T) {
	err := Validate(Config{Cloud: "azure"})
	if err == nil {
		t.Fatal("expected error for missing azure.location")
	}
	if !strings.Contains(err.Error(), "azure.location") {
		t.Fatalf("expected azure.location error, got: %v", err)
	}
}

func TestValidateGCPRegionRequired(t *testing.T) {
	err := Validate(Config{Cloud: "gcp"})
	if err == nil {
		t.Fatal("expected error for missing gcp.region")
	}
	if !strings.Contains(err.Error(), "gcp.region") {
		t.Fatalf("expected gcp.region error, got: %v", err)
	}
}

func TestValidateLLMBackendOneof(t *testing.T) {
	err := Validate(Config{
		Cloud: "aws",
		AWS:   AWSConfig{Region: "us-east-1"},
		LLM:   LLMConfig{Backend: "banana"},
	})
	if err == nil {
		t.Fatal("expected error for invalid llm.backend")
	}
	if !strings.Contains(err.Error(), "must be one of") {
		t.Fatalf("expected oneof error for llm.backend, got: %v", err)
	}
}

func TestValidateValidAWSConfig(t *testing.T) {
	err := Validate(Config{Cloud: "aws", AWS: AWSConfig{Region: "us-east-1"}})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateValidAzureConfig(t *testing.T) {
	err := Validate(Config{Cloud: "azure", Azure: AzureConfig{Location: "eastus"}})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateValidGCPConfig(t *testing.T) {
	err := Validate(Config{Cloud: "gcp", GCP: GCPConfig{Region: "us-central1"}})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateEmptyModeAllowed(t *testing.T) {
	err := Validate(Config{Cloud: "aws", AWS: AWSConfig{Region: "us-east-1"}})
	if err != nil {
		t.Fatalf("expected no error for empty mode, got: %v", err)
	}
}

func TestValidateEmptyLLMBackendAllowed(t *testing.T) {
	err := Validate(Config{Cloud: "aws", AWS: AWSConfig{Region: "us-east-1"}, LLM: LLMConfig{}})
	if err != nil {
		t.Fatalf("expected no error for empty llm.backend, got: %v", err)
	}
}
