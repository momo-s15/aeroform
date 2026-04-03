package bootstrap

import (
	"strings"
	"testing"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/providers"
)

func TestRunAWS(t *testing.T) {
	report := Run(config.Config{
		Cloud: "aws",
		AWS:   config.AWSConfig{Region: "us-east-1", AccountID: "123456789012"},
	}, Params{Repo: "myorg/myrepo"})

	if report.Cloud != "aws" {
		t.Fatalf("expected aws cloud, got %q", report.Cloud)
	}
	if len(report.Sections) < 3 {
		t.Fatalf("expected at least 3 sections, got %d", len(report.Sections))
	}

	rendered := Render(report)
	if !strings.Contains(rendered, "OIDC") {
		t.Fatal("expected OIDC section in output")
	}
	if !strings.Contains(rendered, "123456789012") {
		t.Fatal("expected account ID in output")
	}
	if !strings.Contains(rendered, "myorg/myrepo") {
		t.Fatal("expected repo in output")
	}
	if !strings.Contains(rendered, "s3") {
		t.Fatal("expected S3 state bucket commands")
	}
}

func TestRunAzure(t *testing.T) {
	report := Run(config.Config{
		Cloud: "azure",
		Azure: config.AzureConfig{Location: "eastus"},
	}, Params{Repo: "myorg/myrepo"})

	if report.Cloud != "azure" {
		t.Fatalf("expected azure cloud, got %q", report.Cloud)
	}
	if len(report.Sections) < 3 {
		t.Fatalf("expected at least 3 sections, got %d", len(report.Sections))
	}

	rendered := Render(report)
	if !strings.Contains(rendered, "federated") {
		t.Fatal("expected federated credential section")
	}
	if !strings.Contains(rendered, "myorg/myrepo") {
		t.Fatal("expected repo in output")
	}
	if !strings.Contains(rendered, "storage") {
		t.Fatal("expected storage account commands")
	}
}

func TestRunGCP(t *testing.T) {
	report := Run(config.Config{
		Cloud: "gcp",
		GCP:   config.GCPConfig{Region: "us-central1"},
	}, Params{Repo: "myorg/myrepo"})

	if report.Cloud != "gcp" {
		t.Fatalf("expected gcp cloud, got %q", report.Cloud)
	}
	if len(report.Sections) < 3 {
		t.Fatalf("expected at least 3 sections, got %d", len(report.Sections))
	}

	rendered := Render(report)
	if !strings.Contains(rendered, "Workload Identity") {
		t.Fatal("expected Workload Identity section")
	}
	if !strings.Contains(rendered, "myorg/myrepo") {
		t.Fatal("expected repo in output")
	}
	if !strings.Contains(rendered, "gcs") || !strings.Contains(rendered, "buckets") {
		t.Fatal("expected GCS state bucket commands")
	}
}

func TestRunWithoutRepo(t *testing.T) {
	report := Run(config.Config{
		Cloud: "aws",
		AWS:   config.AWSConfig{Region: "us-east-1"},
	}, Params{})

	rendered := Render(report)
	if !strings.Contains(rendered, "<owner/repo>") {
		t.Fatal("expected placeholder repo when --repo is not provided")
	}
}

func TestRender(t *testing.T) {
	text := Render(Report{
		Cloud: "aws",
		Repo:  "org/repo",
		Sections: []providers.BootstrapSection{
			{Title: "Step One", Commands: []string{"echo hello"}},
			{Title: "Step Two", Commands: []string{"echo world"}},
		},
	})
	if !strings.Contains(text, "Step One") || !strings.Contains(text, "Step Two") {
		t.Fatalf("expected rendered sections")
	}
	if !strings.Contains(text, "org/repo") {
		t.Fatalf("expected repo in render output")
	}
}
