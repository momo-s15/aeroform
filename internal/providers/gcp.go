package providers

import (
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
)

type GCPProvider struct{}

func (GCPProvider) Name() string { return "gcp" }

func (GCPProvider) Region(cfg config.Config) string {
	return cfg.GCP.Region
}

func (GCPProvider) BootstrapSteps(cfg config.Config) []string {
	steps := []string{
		"Validate Workload Identity provider and service account binding",
		"Validate project-level IAM for Terraform operations",
	}
	if cfg.State.Backend == "gcs" {
		steps = append(steps, "Ensure GCS bucket exists for Terraform state")
	}
	return steps
}

func (GCPProvider) SupportedProTemplates() []string {
	return []string{"vpc", "gke", "cloudsql", "gcs", "cloud-run"}
}

func (GCPProvider) SelectProTemplates(prompt string) []string {
	value := strings.ToLower(prompt)
	templates := []string{"vpc"}

	if strings.Contains(value, "gke") || strings.Contains(value, "kubernetes") {
		templates = append(templates, "gke")
	}
	if strings.Contains(value, "database") || strings.Contains(value, "sql") {
		templates = append(templates, "cloudsql")
	}
	if strings.Contains(value, "bucket") || strings.Contains(value, "storage") {
		templates = append(templates, "gcs")
	}
	if len(templates) == 1 {
		templates = append(templates, "cloud-run")
	}

	return templates
}
