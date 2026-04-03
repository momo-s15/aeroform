//go:build integration

package integration

import (
	"os"
	"testing"

	"github.com/momo-s15/aeroform/internal/cost"
	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/security"
	"github.com/momo-s15/aeroform/internal/terraform"
)

func TestSimpleStaticSitePipeline(t *testing.T) {
	ensureLocalStack(t)

	plan, err := engine.BuildSimpleLaunchPlan(engine.SimpleLaunchInput{
		Prompt:   "a personal website",
		Provider: "aws",
	})
	if err != nil {
		t.Fatalf("BuildSimpleLaunchPlan: %v", err)
	}
	if plan.Template != "static-site" {
		t.Fatalf("expected static-site template, got %q", plan.Template)
	}

	report := security.EvaluateSimplePlan(plan)
	if report.BlockingCount > 0 {
		t.Fatal("security gate should not block static-site")
	}

	workDir := t.TempDir()
	if err := terraform.RenderTemplate(plan.TemplateDir, nil, workDir); err != nil {
		t.Fatalf("RenderTemplate: %v", err)
	}

	vars := map[string]string{
		"project_name": plan.ProjectName,
		"region":       "us-east-1",
	}
	if err := terraform.GenerateTfvars(vars, workDir); err != nil {
		t.Fatalf("GenerateTfvars: %v", err)
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
		t.Fatal("expected resources to add in plan")
	}
}

func TestSimpleCostEstimation(t *testing.T) {
	templates := cost.SupportedSimpleTemplates()
	if len(templates) == 0 {
		t.Fatal("expected supported templates")
	}

	for _, tmpl := range templates {
		est := cost.EstimateForSimpleTemplate("aws", tmpl, false, 20)
		if est.Template != tmpl {
			t.Errorf("template mismatch: got %q, want %q", est.Template, tmpl)
		}
		if len(est.Lines) == 0 {
			t.Errorf("expected cost lines for %q", tmpl)
		}
	}
}

func TestSimpleSecurityGateBlocksInvalid(t *testing.T) {
	plan := engine.SimpleLaunchPlan{
		Provider:    "unsupported-cloud",
		ProjectName: "test-project",
		Template:    "static-site",
	}
	report := security.EvaluateSimplePlan(plan)
	if report.BlockingCount == 0 {
		t.Fatal("expected security gate to block unsupported provider")
	}
}

func ensureLocalStack(t *testing.T) {
	t.Helper()
	if os.Getenv("LOCALSTACK_ENDPOINT") == "" && os.Getenv("AWS_ENDPOINT_URL") == "" {
		t.Setenv("AWS_ENDPOINT_URL", "http://localhost:4566")
	}
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_DEFAULT_REGION", "us-east-1")
}
