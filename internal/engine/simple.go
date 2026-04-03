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
	Prompt       string
	Provider     string
	ProjectName  string
	CustomDomain string
}

type SimpleLaunchPlan struct {
	Provider     string
	ProjectName  string
	Prompt       string
	Template     string
	TemplateDir  string
	CustomDomain string
	Cost         cost.Estimate
	Summary      []string
}

func BuildSimpleLaunchPlan(input SimpleLaunchInput) (SimpleLaunchPlan, error) {
	return BuildSimpleLaunchPlanWithClient(input, nil)
}

func BuildSimpleLaunchPlanWithClient(input SimpleLaunchInput, client llm.Client) (SimpleLaunchPlan, error) {
	provider := strings.ToLower(strings.TrimSpace(input.Provider))
	if provider == "" {
		provider = "aws"
	}
	if provider != "aws" {
		return SimpleLaunchPlan{}, fmt.Errorf("simple launch is only implemented for aws in this phase")
	}

	template := SelectSimpleTemplate(input.Prompt)
	if client != nil {
		selectionPrompt := prompt.BuildSimplePrompt(input.Prompt, provider, "us-east-1", cost.SupportedSimpleTemplates())
		selection, err := llm.GenerateTemplateSelection(client, selectionPrompt, cost.SupportedSimpleTemplates())
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
	costEstimate := cost.EstimateForSimpleTemplate(template, strings.TrimSpace(input.CustomDomain) != "", 20)

	return SimpleLaunchPlan{
		Provider:     provider,
		ProjectName:  projectName,
		Prompt:       strings.TrimSpace(input.Prompt),
		Template:     template,
		TemplateDir:  templateDir,
		CustomDomain: strings.TrimSpace(input.CustomDomain),
		Cost:         costEstimate,
		Summary: []string{
			fmt.Sprintf("provider: %s", provider),
			fmt.Sprintf("project: %s", projectName),
			fmt.Sprintf("template: %s", template),
			fmt.Sprintf("template_dir: %s", templateDir),
		},
	}, nil
}

func DefaultSimpleLLMClient() llm.Client {
	return llm.NewClient(llm.Config{Backend: "ollama", Host: "http://localhost:11434", Model: "llama3.2"})
}

func SelectSimpleTemplate(prompt string) string {
	value := strings.ToLower(prompt)
	switch {
	case strings.Contains(value, "contact"):
		return "contact-form"
	case strings.Contains(value, "api"), strings.Contains(value, "backend"), strings.Contains(value, "lambda"):
		return "lambda-api"
	case strings.Contains(value, "database"), strings.Contains(value, "db"):
		return "tiny-db"
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
