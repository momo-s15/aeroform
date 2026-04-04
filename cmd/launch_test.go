package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/momo-s15/aeroform/internal/engine"
)

func TestCheckAzureSimpleTemplateFresh(t *testing.T) {
	dir := t.TempDir()
	plan := engine.SimpleLaunchPlan{Provider: "azure", Template: "static-site"}

	if err := os.WriteFile(filepath.Join(dir, "main.tf"), []byte("# legacy template\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := checkAzureSimpleTemplateFresh(dir, plan); err == nil {
		t.Fatal("expected error for template without schema marker")
	}

	if err := os.WriteFile(filepath.Join(dir, "main.tf"), []byte("# aeroform-schema: azure-static-site/2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := checkAzureSimpleTemplateFresh(dir, plan); err != nil {
		t.Fatalf("expected ok with marker: %v", err)
	}

	awsPlan := engine.SimpleLaunchPlan{Provider: "aws", Template: "static-site"}
	if err := checkAzureSimpleTemplateFresh(dir, awsPlan); err != nil {
		t.Fatalf("non-azure should skip check: %v", err)
	}
}
