package providers

import (
	"bytes"
	"strings"
	"testing"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/shopspring/decimal"
)

func TestForCloud(t *testing.T) {
	tests := []struct {
		cloud string
		want  string
	}{
		{cloud: "aws", want: "aws"},
		{cloud: "azure", want: "azure"},
		{cloud: "gcp", want: "gcp"},
		{cloud: "unknown", want: "aws"},
	}

	for _, testCase := range tests {
		t.Run(testCase.cloud, func(t *testing.T) {
			provider := ForCloud(testCase.cloud)
			if provider.Name() != testCase.want {
				t.Fatalf("ForCloud(%q) = %q, want %q", testCase.cloud, provider.Name(), testCase.want)
			}
		})
	}
}

func TestAWSProviderValidate(t *testing.T) {
	provider := AWSProvider{}

	if err := provider.Validate(config.Config{Cloud: "aws", AWS: config.AWSConfig{Region: "us-east-1", AccountID: "123456789012"}}); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}

	err := provider.Validate(config.Config{Cloud: "aws", AWS: config.AWSConfig{AccountID: "123456789012"}})
	if err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "aws.region") {
		t.Fatalf("Validate() error = %v, want aws.region error", err)
	}
}

func TestAWSProviderGetTemplateDir(t *testing.T) {
	provider := AWSProvider{}

	if got, want := provider.GetTemplateDir(engine.SimpleMode, "static-site"), "templates/simple/aws/static-site"; got != want {
		t.Fatalf("GetTemplateDir(SimpleMode) = %q, want %q", got, want)
	}
	if got, want := provider.GetTemplateDir(engine.ProMode, "eks"), "templates/pro/aws/eks"; got != want {
		t.Fatalf("GetTemplateDir(ProMode) = %q, want %q", got, want)
	}
}

func TestAWSEstimateCost(t *testing.T) {
	est := AWSProvider{}.EstimateCost([]string{"vpc", "eks", "s3-private"})
	if !est.Monthly.GreaterThan(decimal.Zero) {
		t.Fatal("expected positive cost for vpc+eks+s3-private")
	}
	if len(est.Items) != 3 {
		t.Fatalf("expected 3 line items, got %d", len(est.Items))
	}
	if len(est.Lines()) == 0 {
		t.Fatal("expected non-empty Lines()")
	}
}

func TestAzureEstimateCost(t *testing.T) {
	est := AzureProvider{}.EstimateCost([]string{"vnet", "aks"})
	if !est.Monthly.GreaterThan(decimal.Zero) {
		t.Fatal("expected positive cost for vnet+aks")
	}
	if len(est.Items) != 2 {
		t.Fatalf("expected 2 line items, got %d", len(est.Items))
	}
}

func TestGCPEstimateCost(t *testing.T) {
	est := GCPProvider{}.EstimateCost([]string{"vpc", "gke", "cloudsql"})
	if !est.Monthly.GreaterThan(decimal.Zero) {
		t.Fatal("expected positive cost for vpc+gke+cloudsql")
	}
	if len(est.Items) != 3 {
		t.Fatalf("expected 3 line items, got %d", len(est.Items))
	}
}

func TestEstimateCostUnknownTemplate(t *testing.T) {
	est := AWSProvider{}.EstimateCost([]string{"nonexistent"})
	if len(est.Items) != 1 {
		t.Fatalf("expected 1 line item, got %d", len(est.Items))
	}
	if !strings.Contains(est.Items[0].Note, "not yet catalogued") {
		t.Fatalf("expected 'not yet catalogued' note, got %q", est.Items[0].Note)
	}
}

func TestAWSPostDeploy(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.Config{AWS: config.AWSConfig{Region: "us-east-1"}}
	p := AWSProvider{}
	if err := p.PostDeploy(cfg, &buf); err != nil {
		t.Fatalf("PostDeploy error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "post-deploy") {
		t.Fatal("expected post-deploy checklist in output")
	}
	if !strings.Contains(out, "us-east-1") {
		t.Fatal("expected region in EKS instructions")
	}
}

func TestAzurePostDeploy(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.Config{Azure: config.AzureConfig{Location: "eastus"}}
	p := AzureProvider{}
	if err := p.PostDeploy(cfg, &buf); err != nil {
		t.Fatalf("PostDeploy error: %v", err)
	}
	if !strings.Contains(buf.String(), "post-deploy") {
		t.Fatal("expected post-deploy checklist in output")
	}
}

func TestGCPPostDeploy(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.Config{GCP: config.GCPConfig{Region: "us-central1"}}
	p := GCPProvider{}
	if err := p.PostDeploy(cfg, &buf); err != nil {
		t.Fatalf("PostDeploy error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "post-deploy") {
		t.Fatal("expected post-deploy checklist in output")
	}
	if !strings.Contains(out, "us-central1") {
		t.Fatal("expected region in GKE instructions")
	}
}
