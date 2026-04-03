//go:build e2e

package e2e

import (
	"path/filepath"
	"testing"

	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/security"
	"github.com/momo-s15/aeroform/internal/terraform"
)

func TestSimpleModeDeployStaticSite(t *testing.T) {
	requireAWSCredentials(t)

	plan, err := engine.BuildSimpleLaunchPlan(engine.SimpleLaunchInput{
		Prompt:      "a personal website",
		Provider:    "aws",
		ProjectName: "aeroform-e2e-static",
	})
	if err != nil {
		t.Fatalf("BuildSimpleLaunchPlan: %v", err)
	}

	report := security.EvaluateSimplePlan(plan)
	if report.BlockingCount > 0 {
		t.Fatal("security gate should not block static-site")
	}

	workDir := filepath.Join(t.TempDir(), "deploy")
	if err := terraform.RenderTemplate(plan.TemplateDir, nil, workDir); err != nil {
		t.Fatalf("RenderTemplate: %v", err)
	}

	vars := map[string]string{
		"project_name": plan.ProjectName,
		"region":       awsRegion(),
	}
	if err := terraform.GenerateTfvars(vars, workDir); err != nil {
		t.Fatalf("GenerateTfvars: %v", err)
	}

	if err := terraform.EnsureBinary(); err != nil {
		t.Fatalf("terraform not available: %v", err)
	}

	if err := terraform.Init(workDir); err != nil {
		t.Fatalf("terraform init: %v", err)
	}

	defer func() {
		t.Log("destroying infrastructure...")
		if err := terraform.Destroy(workDir, true); err != nil {
			t.Errorf("terraform destroy failed: %v", err)
		}
	}()

	planResult, err := terraform.Plan(workDir)
	if err != nil {
		t.Fatalf("terraform plan: %v", err)
	}
	if planResult.AddCount == 0 {
		t.Fatal("expected resources to add")
	}
	t.Logf("plan: +%d ~%d -%d", planResult.AddCount, planResult.ChangeCount, planResult.DestroyCount)

	if err := terraform.Apply(workDir, true); err != nil {
		t.Fatalf("terraform apply: %v", err)
	}

	outputs, err := terraform.Output(workDir)
	if err != nil {
		t.Logf("terraform output warning: %v", err)
	} else {
		t.Logf("outputs:\n%s", outputs)
	}
}

func TestSimpleModeDeployLambdaAPI(t *testing.T) {
	requireAWSCredentials(t)

	plan, err := engine.BuildSimpleLaunchPlan(engine.SimpleLaunchInput{
		Prompt:      "a backend api",
		Provider:    "aws",
		ProjectName: "aeroform-e2e-api",
	})
	if err != nil {
		t.Fatalf("BuildSimpleLaunchPlan: %v", err)
	}
	if plan.Template != "lambda-api" {
		t.Fatalf("expected lambda-api, got %q", plan.Template)
	}

	workDir := filepath.Join(t.TempDir(), "deploy")
	if err := terraform.RenderTemplate(plan.TemplateDir, nil, workDir); err != nil {
		t.Fatalf("RenderTemplate: %v", err)
	}

	vars := map[string]string{
		"project_name": plan.ProjectName,
		"region":       awsRegion(),
	}
	if err := terraform.GenerateTfvars(vars, workDir); err != nil {
		t.Fatalf("GenerateTfvars: %v", err)
	}

	if err := terraform.EnsureBinary(); err != nil {
		t.Fatalf("terraform not available: %v", err)
	}

	if err := terraform.Init(workDir); err != nil {
		t.Fatalf("terraform init: %v", err)
	}

	defer func() {
		t.Log("destroying infrastructure...")
		if err := terraform.Destroy(workDir, true); err != nil {
			t.Errorf("terraform destroy failed: %v", err)
		}
	}()

	planResult, err := terraform.Plan(workDir)
	if err != nil {
		t.Fatalf("terraform plan: %v", err)
	}
	t.Logf("plan: +%d ~%d -%d", planResult.AddCount, planResult.ChangeCount, planResult.DestroyCount)

	if err := terraform.Apply(workDir, true); err != nil {
		t.Fatalf("terraform apply: %v", err)
	}

	outputs, err := terraform.Output(workDir)
	if err != nil {
		t.Logf("terraform output warning: %v", err)
	} else {
		t.Logf("outputs:\n%s", outputs)
	}
}
