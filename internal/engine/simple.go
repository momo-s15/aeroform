package engine

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/momo-s15/aeroform/internal/cost"
	"github.com/momo-s15/aeroform/internal/llm"
	"github.com/momo-s15/aeroform/internal/prompt"
)

type SimpleLaunchInput struct {
	Prompt        string
	Provider      string
	ProjectName   string
	CustomDomain  string
	GCPProjectID  string // required when Provider is gcp (Terraform google provider project)
	AzureLocation string // Azure region (e.g. canadacentral) when Provider is azure
}

type SimpleLaunchPlan struct {
	Provider      string
	ProjectName   string
	Prompt        string
	Template      string
	TemplateDir   string
	CustomDomain  string
	GCPProjectID  string
	AzureLocation string
	Cost          cost.Estimate
	Summary       []string
}

func BuildSimpleLaunchPlan(input SimpleLaunchInput) (SimpleLaunchPlan, error) {
	return BuildSimpleLaunchPlanWithClient(input, nil)
}

func BuildSimpleLaunchPlanWithClient(input SimpleLaunchInput, client llm.Client) (SimpleLaunchPlan, error) {
	provider := strings.ToLower(strings.TrimSpace(input.Provider))
	if provider == "" {
		provider = "aws"
	}
	supportedProviders := map[string]bool{"aws": true, "azure": true, "gcp": true}
	if !supportedProviders[provider] {
		return SimpleLaunchPlan{}, fmt.Errorf("unsupported provider %q; choose aws, azure, or gcp", provider)
	}

	available := cost.SupportedSimpleTemplatesForProvider(provider)
	template := SelectSimpleTemplateForProvider(input.Prompt, provider)

	if client != nil {
		selectionPrompt := prompt.BuildSimplePrompt(input.Prompt, provider, "us-east-1", available)
		selection, err := llm.GenerateTemplateSelection(client, selectionPrompt, available)
		if err == nil && len(selection.Templates) > 0 {
			template = selection.Templates[0]
		}
	}

	projectName := strings.TrimSpace(input.ProjectName)
	if projectName == "" {
		projectName = Slugify(input.Prompt)
	}
	if projectName == "" {
		projectName = "aeroform-project"
	}

	templateDir := filepath.ToSlash(filepath.Join("templates", "simple", provider, template))
	costEstimate := cost.EstimateForSimpleTemplate(provider, template, strings.TrimSpace(input.CustomDomain) != "", 20)

	plan := SimpleLaunchPlan{
		Provider:      provider,
		ProjectName:   projectName,
		Prompt:        strings.TrimSpace(input.Prompt),
		Template:      template,
		TemplateDir:   templateDir,
		CustomDomain:  strings.TrimSpace(input.CustomDomain),
		GCPProjectID:  strings.TrimSpace(input.GCPProjectID),
		AzureLocation: strings.TrimSpace(strings.ToLower(input.AzureLocation)),
		Cost:          costEstimate,
		Summary: []string{
			fmt.Sprintf("provider: %s", provider),
			fmt.Sprintf("project: %s", projectName),
			fmt.Sprintf("template: %s", template),
			fmt.Sprintf("template_dir: %s", templateDir),
		},
	}
	if provider == "gcp" && plan.GCPProjectID != "" {
		plan.Summary = append(plan.Summary, "gcp_project_id: "+plan.GCPProjectID)
	}
	if provider == "azure" && plan.AzureLocation != "" {
		plan.Summary = append(plan.Summary, "azure_region: "+plan.AzureLocation)
	}
	return plan, nil
}

func DefaultSimpleLLMClient() llm.Client {
	return llm.NewClient(llm.Config{Backend: "ollama", Host: "http://localhost:11434", Model: "llama3.2"})
}

func SelectSimpleTemplate(prompt string) string {
	return SelectSimpleTemplateForProvider(prompt, "aws")
}

func SelectSimpleTemplateForProvider(userPrompt, provider string) string {
	value := strings.ToLower(userPrompt)

	switch provider {
	case "azure":
		return selectAzureTemplate(value)
	case "gcp":
		return selectGCPTemplate(value)
	default:
		return selectAWSTemplate(value)
	}
}

func selectAWSTemplate(value string) string {
	switch {
	case strings.Contains(value, "contact"):
		return "contact-form"
	case strings.Contains(value, "discord") || strings.Contains(value, "bot"):
		return "discord-bot"
	case strings.Contains(value, "game") || strings.Contains(value, "minecraft") || strings.Contains(value, "server"):
		return "game-server"
	case strings.Contains(value, "fullstack") || strings.Contains(value, "full-stack") || strings.Contains(value, "full stack"):
		return "fullstack-app"
	case strings.Contains(value, "upload") || strings.Contains(value, "file"):
		return "file-upload"
	case strings.Contains(value, "short") || strings.Contains(value, "url"):
		return "url-shortener"
	case strings.Contains(value, "cron") || strings.Contains(value, "schedule") || strings.Contains(value, "scheduled"):
		return "cron-job"
	case strings.Contains(value, "api") || strings.Contains(value, "backend") || strings.Contains(value, "lambda"):
		return "lambda-api"
	case strings.Contains(value, "database") || strings.Contains(value, "db"):
		return "tiny-db"
	default:
		return "static-site"
	}
}

func selectAzureTemplate(value string) string {
	switch {
	case strings.Contains(value, "api") || strings.Contains(value, "backend") ||
		strings.Contains(value, "function") || strings.Contains(value, "serverless"):
		return "function-api"
	default:
		return "static-site"
	}
}

func selectGCPTemplate(value string) string {
	switch {
	case strings.Contains(value, "api") || strings.Contains(value, "backend") ||
		strings.Contains(value, "container") || strings.Contains(value, "cloud run"):
		return "cloud-run-api"
	default:
		return "static-site"
	}
}

func Slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}
	re := regexp.MustCompile(`[^a-z0-9]+`)
	value = re.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	if value == "" {
		return ""
	}
	return value
}
