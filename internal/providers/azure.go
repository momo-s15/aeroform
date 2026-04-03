package providers

import (
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
)

type AzureProvider struct{}

func (AzureProvider) Name() string { return "azure" }

func (AzureProvider) Region(cfg config.Config) string {
	return cfg.Azure.Location
}

func (AzureProvider) BootstrapSteps(cfg config.Config) []string {
	steps := []string{
		"Validate federated credential / workload identity setup",
		"Validate resource group and subscription permissions",
	}
	if cfg.State.Backend == "azurerm" || cfg.State.Backend == "azureblob" {
		steps = append(steps, "Ensure storage account and container exist for Terraform state")
	}
	return steps
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
