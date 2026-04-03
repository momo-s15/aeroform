package providers

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/cost"
	"github.com/momo-s15/aeroform/internal/engine"
)

type GCPProvider struct{}

func (GCPProvider) Name() string { return "gcp" }

func (GCPProvider) Validate(cfg config.Config) error {
	if strings.TrimSpace(cfg.Cloud) == "" {
		return fmt.Errorf("cloud is required")
	}
	if strings.TrimSpace(cfg.GCP.Region) == "" {
		return fmt.Errorf("gcp.region is required")
	}
	return nil
}

func (GCPProvider) GenerateVars(cfg config.Config, params map[string]string) map[string]string {
	vars := map[string]string{
		"cloud":          cfg.Cloud,
		"mode":           cfg.Mode,
		"aws_region":     cfg.AWS.Region,
		"aws_account_id": cfg.AWS.AccountID,
		"azure_location": cfg.Azure.Location,
		"gcp_region":     cfg.GCP.Region,
		"state_backend":  cfg.State.Backend,
		"llm_backend":    cfg.LLM.Backend,
		"llm_model":      cfg.LLM.Model,
		"llm_host":       cfg.LLM.Host,
		"llm_region":     cfg.LLM.Region,
		"llm_api_key":    cfg.LLM.APIKey,
		"llm_role_arn":   cfg.LLM.RoleARN,
	}
	for key, value := range params {
		vars[key] = value
	}
	return vars
}

func (GCPProvider) GetTemplateDir(mode engine.Mode, tmpl string) string {
	modeDir := "simple"
	if mode == engine.ProMode {
		modeDir = "pro"
	}
	return path.Join("templates", modeDir, "gcp", tmpl)
}

func (GCPProvider) EstimateCost(templates []string) CostEstimate {
	return cost.EstimateForProTemplates("gcp", templates)
}

func (GCPProvider) PostDeploy(cfg config.Config, w io.Writer) error {
	fmt.Fprintln(w, "GCP post-deploy checklist:")
	fmt.Fprintln(w, "  1. Run 'terraform output' in the project directory to see endpoints.")
	fmt.Fprintln(w, "  2. If GKE was deployed: gcloud container clusters get-credentials <cluster> --region "+cfg.GCP.Region)
	fmt.Fprintln(w, "  3. Verify Workload Identity bindings and IAM policies in the GCP Console.")
	fmt.Fprintln(w, "  4. Set up Cloud Monitoring alerting policies for budget and health.")
	return nil
}

func (GCPProvider) Region(cfg config.Config) string {
	return cfg.GCP.Region
}

func (GCPProvider) BootstrapSections(cfg config.Config, repo string) []BootstrapSection {
	region := cfg.GCP.Region
	if region == "" {
		region = "us-central1"
	}
	if repo == "" {
		repo = "<owner/repo>"
	}

	poolName := "aeroform-github-pool"
	providerName := "aeroform-github-provider"
	saName := "aeroform-deployer"
	bucketName := "aeroform-state-${PROJECT_ID}"

	return []BootstrapSection{
		{
			Title: "Create Workload Identity pool + provider for GitHub Actions",
			Commands: []string{
				"PROJECT_ID=$(gcloud config get-value project)",
				"",
				"gcloud iam workload-identity-pools create " + poolName + " \\",
				"  --location=global \\",
				`  --display-name="Aeroform GitHub Actions pool"`,
				"",
				"gcloud iam workload-identity-pools providers create-oidc " + providerName + " \\",
				"  --location=global \\",
				"  --workload-identity-pool=" + poolName + " \\",
				"  --issuer-uri=https://token.actions.githubusercontent.com \\",
				"  --attribute-mapping=google.subject=assertion.sub,attribute.repository=assertion.repository \\",
				`  --attribute-condition="assertion.repository=='` + repo + `'"`,
			},
		},
		{
			Title: "Create service account + IAM binding",
			Commands: []string{
				"gcloud iam service-accounts create " + saName + " \\",
				`  --display-name="Aeroform Terraform deployer"`,
				"",
				"# Grant Editor role (scope down for production)",
				"gcloud projects add-iam-policy-binding $PROJECT_ID \\",
				"  --member=serviceAccount:" + saName + "@$PROJECT_ID.iam.gserviceaccount.com \\",
				"  --role=roles/editor",
				"",
				"# Allow GitHub Actions to impersonate the service account",
				"gcloud iam service-accounts add-iam-policy-binding \\",
				"  " + saName + "@$PROJECT_ID.iam.gserviceaccount.com \\",
				"  --member=\"principalSet://iam.googleapis.com/projects/$PROJECT_NUMBER/locations/global/workloadIdentityPools/" + poolName + "/attribute.repository/" + repo + "\" \\",
				"  --role=roles/iam.workloadIdentityUser",
				"",
				"# Get the project number (needed for GitHub Actions config)",
				"gcloud projects describe $PROJECT_ID --format='value(projectNumber)'",
			},
		},
		{
			Title: "Create GCS bucket for Terraform state",
			Commands: []string{
				"gcloud storage buckets create gs://" + bucketName + " \\",
				"  --location=" + region + " \\",
				"  --uniform-bucket-level-access \\",
				"  --public-access-prevention=enforced",
				"",
				"gcloud storage buckets update gs://" + bucketName + " --versioning",
			},
		},
		{
			Title: "Add to your Terraform backend config",
			Commands: []string{
				"# Add this block to your Terraform configuration:",
				`terraform {`,
				`  backend "gcs" {`,
				`    bucket = "` + bucketName + `"`,
				`    prefix = "aeroform/terraform"`,
				`  }`,
				`}`,
			},
		},
	}
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
