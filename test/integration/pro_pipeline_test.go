//go:build integration

package integration

import (
	"testing"

	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/providers"
	"github.com/momo-s15/aeroform/internal/security"
	"github.com/momo-s15/aeroform/internal/terraform"
)

func TestProVPCTemplatePlan(t *testing.T) {
	ensureLocalStack(t)

	provider := providers.ForCloud("aws")
	templateDir := provider.GetTemplateDir(engine.ProMode, "vpc")

	workDir := t.TempDir()
	sources := []terraform.ProTemplateSource{
		{Name: "vpc", Dir: templateDir},
	}
	vars := map[string]string{
		"aws_region":     "us-east-1",
		"aws_account_id": "000000000000",
		"cloud":          "aws",
		"mode":           "pro",
	}

	if err := terraform.RenderProProject(sources, vars, workDir); err != nil {
		t.Fatalf("RenderProProject: %v", err)
	}

	if err := terraform.EnsureBinary(); err != nil {
		t.Skipf("terraform not available: %v", err)
	}

	if err := terraform.Init(workDir); err != nil {
		t.Fatalf("terraform init: %v", err)
	}

	planResult, err := terraform.Plan(workDir)
	if err != nil {
		t.Fatalf("terraform plan: %v", err)
	}
	if planResult.AddCount == 0 {
		t.Fatal("expected VPC resources to add")
	}
}

func TestProSecurityGateBlocking(t *testing.T) {
	badDir := t.TempDir()
	report := security.RunProSecurityScan(badDir)

	for _, line := range report.SummaryLines() {
		t.Log(line)
	}
}

func TestProCostEstimate(t *testing.T) {
	provider := providers.ForCloud("aws")
	est := provider.EstimateCost([]string{"vpc", "eks", "rds-private"})
	if est.Monthly.IsZero() {
		t.Fatal("expected non-zero cost for vpc+eks+rds-private")
	}
	lines := est.Lines()
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 cost lines, got %d", len(lines))
	}
}
