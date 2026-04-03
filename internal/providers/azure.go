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

type AzureProvider struct{}

func (AzureProvider) Name() string { return "azure" }

func (AzureProvider) Validate(cfg config.Config) error {
	if strings.TrimSpace(cfg.Cloud) == "" {
		return fmt.Errorf("cloud is required")
	}
	if strings.TrimSpace(cfg.Azure.Location) == "" {
		return fmt.Errorf("azure.location is required")
	}
	return nil
}

func (AzureProvider) GenerateVars(cfg config.Config, params map[string]string) map[string]string {
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

func (AzureProvider) GetTemplateDir(mode engine.Mode, tmpl string) string {
	modeDir := "simple"
	if mode == engine.ProMode {
		modeDir = "pro"
	}
	return path.Join("templates", modeDir, "azure", tmpl)
}

func (AzureProvider) EstimateCost(templates []string) CostEstimate {
	return cost.EstimateForProTemplates("azure", templates)
}

func (AzureProvider) PostDeploy(cfg config.Config, w io.Writer) error {
	fmt.Fprintln(w, "Azure post-deploy checklist:")
	fmt.Fprintln(w, "  1. Run 'terraform output' in the project directory to see endpoints.")
	fmt.Fprintln(w, "  2. If AKS was deployed: az aks get-credentials --name <cluster> --resource-group <rg>")
	fmt.Fprintln(w, "  3. Verify managed identities and role assignments in the Azure Portal.")
	fmt.Fprintln(w, "  4. Set up Azure Monitor alerts for cost and resource health.")
	return nil
}

func (AzureProvider) Region(cfg config.Config) string {
	return cfg.Azure.Location
}

func (AzureProvider) BootstrapSections(cfg config.Config, repo string) []BootstrapSection {
	location := cfg.Azure.Location
	if location == "" {
		location = "eastus"
	}
	if repo == "" {
		repo = "<owner/repo>"
	}

	appName := "aeroform-github-actions"
	rgName := "rg-aeroform-state"
	saName := "aeroformstate"
	containerName := "tfstate"

	return []BootstrapSection{
		{
			Title: "Create Azure AD app registration + federated credential",
			Commands: []string{
				"# Create the app registration",
				"az ad app create --display-name " + appName,
				"",
				"# Get the app (client) ID — save this for your GitHub secret",
				"APP_ID=$(az ad app list --display-name " + appName + " --query '[0].appId' -o tsv)",
				"",
				"# Create a service principal for the app",
				"az ad sp create --id $APP_ID",
				"",
				"# Get your subscription ID",
				"SUB_ID=$(az account show --query id -o tsv)",
				"",
				"# Assign Contributor role on the subscription",
				"az role assignment create \\",
				"  --assignee $APP_ID \\",
				"  --role Contributor \\",
				"  --scope /subscriptions/$SUB_ID",
				"",
				"# Create federated credential for GitHub Actions",
				`az ad app federated-credential create --id $APP_ID --parameters '{`,
				`  "name": "aeroform-github-main",`,
				`  "issuer": "https://token.actions.githubusercontent.com",`,
				`  "subject": "repo:` + repo + `:ref:refs/heads/main",`,
				`  "audiences": ["api://AzureADTokenExchange"]`,
				`}'`,
			},
		},
		{
			Title: "Create Storage Account for Terraform state",
			Commands: []string{
				"az group create --name " + rgName + " --location " + location,
				"",
				"az storage account create \\",
				"  --name " + saName + " \\",
				"  --resource-group " + rgName + " \\",
				"  --location " + location + " \\",
				"  --sku Standard_LRS \\",
				"  --min-tls-version TLS1_2 \\",
				"  --allow-blob-public-access false",
				"",
				"az storage container create \\",
				"  --name " + containerName + " \\",
				"  --account-name " + saName,
			},
		},
		{
			Title: "Add to your Terraform backend config",
			Commands: []string{
				"# Add this block to your Terraform configuration:",
				`terraform {`,
				`  backend "azurerm" {`,
				`    resource_group_name  = "` + rgName + `"`,
				`    storage_account_name = "` + saName + `"`,
				`    container_name       = "` + containerName + `"`,
				`    key                  = "aeroform/terraform.tfstate"`,
				`  }`,
				`}`,
			},
		},
	}
}

func (AzureProvider) SupportedProTemplates() []string {
	return []string{"vnet", "aks", "cosmos-db", "app-service", "storage", "key-vault"}
}

func (AzureProvider) SelectProTemplates(prompt string) []string {
	value := strings.ToLower(prompt)
	templates := []string{"vnet"}

	if strings.Contains(value, "aks") || strings.Contains(value, "kubernetes") {
		templates = append(templates, "aks")
	}
	if strings.Contains(value, "database") || strings.Contains(value, "cosmos") {
		templates = append(templates, "cosmos-db")
	}
	if strings.Contains(value, "web") || strings.Contains(value, "app") {
		templates = append(templates, "app-service")
	}
	if len(templates) == 1 {
		templates = append(templates, "storage")
	}

	return templates
}
