//go:build e2e

package e2e

import (
	"path/filepath"
	"testing"

	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/providers"
	"github.com/momo-s15/aeroform/internal/terraform"
)

func TestProModeDeployVPC(t *testing.T) {
	requireAWSCredentials(t)

	provider := providers.ForCloud("aws")
	templateDir := provider.GetTemplateDir(engine.ProMode, "vpc")

	workDir := filepath.Join(t.TempDir(), "pro-vpc")
	sources := []terraform.ProTemplateSource{
		{Name: "vpc", Dir: templateDir},
	}
	vars := map[string]string{
		"aws_region":     awsRegion(),
		"aws_account_id": awsAccountID(t),
		"cloud":          "aws",
		"mode":           "pro",
	}

	if err := terraform.RenderProProject(sources, vars, workDir); err != nil {
		t.Fatalf("RenderProProject: %v", err)
	}

	if err := terraform.EnsureBinary(); err != nil {
		t.Fatalf("terraform not available: %v", err)
	}

	if err := terraform.Init(workDir); err != nil {
		t.Fatalf("terraform init: %v", err)
	}

	defer func() {
		t.Log("destroying VPC infrastructure...")
		if err := terraform.Destroy(workDir, true); err != nil {
			t.Errorf("terraform destroy failed: %v", err)
		}
	}()

	planResult, err := terraform.Plan(workDir)
	if err != nil {
		t.Fatalf("terraform plan: %v", err)
	}
	if planResult.AddCount == 0 {
		t.Fatal("expected VPC resources to add")
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
