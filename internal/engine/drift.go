package engine

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/terraform"
)

type DriftChange struct {
	Address    string
	ChangeType string
}

type DriftReport struct {
	HasDrift bool
	Changes  []DriftChange
	Summary  []string
	Add      int
	Change   int
	Destroy  int
}

// BuildDriftReport runs terraform plan -detailed-exitcode against an existing
// project directory and returns a structured drift report.
func BuildDriftReport(cfg config.Config, prompt string) (DriftReport, error) {
	workDir := resolveDriftWorkDir(prompt)

	if _, err := os.Stat(filepath.Join(workDir, ".terraform")); os.IsNotExist(err) {
		return DriftReport{}, fmt.Errorf("no initialized Terraform state found in %s — run 'aeroform generate' first", workDir)
	}

	result, err := terraform.PlanDetailed(workDir)
	if err != nil {
		return DriftReport{}, fmt.Errorf("drift check failed: %w", err)
	}

	report := DriftReport{
		HasDrift: result.HasDrift,
		Add:      result.Plan.AddCount,
		Change:   result.Plan.ChangeCount,
		Destroy:  result.Plan.DestroyCount,
	}

	if !result.HasDrift {
		report.Summary = []string{
			fmt.Sprintf("cloud: %s", cfg.Cloud),
			fmt.Sprintf("project: %s", workDir),
			"drift status: No drift detected. Infrastructure matches desired state.",
		}
		return report, nil
	}

	for _, c := range result.Changes {
		report.Changes = append(report.Changes, DriftChange{
			Address:    c.Address,
			ChangeType: c.ChangeType,
		})
	}

	summary := []string{
		fmt.Sprintf("cloud: %s", cfg.Cloud),
		fmt.Sprintf("project: %s", workDir),
		fmt.Sprintf("drift status: %d to add, %d to change, %d to destroy",
			report.Add, report.Change, report.Destroy),
		"",
		"Drifted resources:",
	}
	for _, c := range report.Changes {
		summary = append(summary, fmt.Sprintf("  %-8s %s", c.ChangeType, c.Address))
	}
	summary = append(summary, "", "Run 'aeroform generate' with the same prompt to reconcile.")

	report.Summary = summary
	return report, nil
}

func resolveDriftWorkDir(prompt string) string {
	slug := Slugify(prompt)
	proDir := filepath.Join(".aeroform", "pro", slug)
	if info, err := os.Stat(proDir); err == nil && info.IsDir() {
		return proDir
	}
	simpleDir := filepath.Join(".aeroform", "projects", slug)
	if info, err := os.Stat(simpleDir); err == nil && info.IsDir() {
		return simpleDir
	}
	return proDir
}
