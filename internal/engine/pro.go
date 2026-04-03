package engine

import (
	"fmt"
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/llm"
	"github.com/momo-s15/aeroform/internal/prompt"
)

type ProPlan struct {
	Prompt    string
	Templates []string
	Summary   []string
}

func BuildProPlan(cfg config.Config, prompt string) (ProPlan, error) {
	return BuildProPlanWithClient(cfg, prompt, nil)
}

func BuildProPlanWithClient(cfg config.Config, promptText string, client llm.Client) (ProPlan, error) {
	promptText = strings.TrimSpace(promptText)
	if promptText == "" {
		return ProPlan{}, fmt.Errorf("prompt is required")
	}

	templates := selectProTemplates(cfg.Cloud, promptText)
	if client != nil {
		allowed := supportedProTemplates(cfg.Cloud)
		promptBody := prompt.BuildProPrompt(promptText, cfg.Cloud, regionForCloud(cfg), "production", allowed)
		selection, err := llm.GenerateTemplateSelection(client, promptBody, allowed)
		if err == nil && len(selection.Templates) > 0 {
			templates = selection.Templates
		}
	}

	region := regionForCloud(cfg)

	return ProPlan{
		Prompt:    promptText,
		Templates: templates,
		Summary: []string{
			fmt.Sprintf("cloud: %s", cfg.Cloud),
			fmt.Sprintf("region: %s", region),
			fmt.Sprintf("templates: %s", strings.Join(templates, ", ")),
			"pipeline: render -> security scan -> terraform plan/apply",
		},
	}, nil
}

func DefaultProLLMClient() llm.Client {
	return llm.NewClient(llm.Config{Backend: "ollama", Host: "http://localhost:11434", Model: "llama3.2"})
}

func ProLLMClientFromConfig(cfg config.Config) llm.Client {
	backend := cfg.LLM.Backend
	if backend == "" {
		backend = "ollama"
	}
	return llm.NewClient(llm.Config{
		Backend: backend,
		Host:    cfg.LLM.Host,
		Model:   cfg.LLM.Model,
		Region:  cfg.LLM.Region,
		APIKey:  cfg.LLM.APIKey,
		RoleARN: cfg.LLM.RoleARN,
	})
}

func regionForCloud(cfg config.Config) string {
	switch strings.ToLower(strings.TrimSpace(cfg.Cloud)) {
	case "azure":
		return cfg.Azure.Location
	case "gcp":
		return cfg.GCP.Region
	default:
		return cfg.AWS.Region
	}
}

func supportedProTemplates(cloud string) []string {
	switch strings.ToLower(strings.TrimSpace(cloud)) {
	case "azure":
		return []string{"vnet", "aks", "cosmos-db", "app-service", "storage", "key-vault"}
	case "gcp":
		return []string{"vpc", "gke", "cloudsql", "gcs", "cloud-run"}
	default:
		return []string{"vpc", "eks", "rds-private", "s3-private", "alb", "lambda-api", "ecs-fargate", "cloudfront-api"}
	}
}

func selectProTemplates(cloud, prompt string) []string {
	value := strings.ToLower(prompt)
	switch strings.ToLower(strings.TrimSpace(cloud)) {
	case "azure":
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
	case "gcp":
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
	default:
		templates := []string{"vpc"}
		if strings.Contains(value, "eks") || strings.Contains(value, "kubernetes") {
			templates = append(templates, "eks")
		}
		if strings.Contains(value, "rds") || strings.Contains(value, "database") {
			templates = append(templates, "rds-private")
		}
		if strings.Contains(value, "s3") || strings.Contains(value, "bucket") || strings.Contains(value, "storage") {
			templates = append(templates, "s3-private")
		}
		if strings.Contains(value, "alb") || strings.Contains(value, "load balancer") {
			templates = append(templates, "alb")
		}
		if strings.Contains(value, "ecs") || strings.Contains(value, "fargate") || strings.Contains(value, "container") {
			templates = append(templates, "ecs-fargate")
		}
		if strings.Contains(value, "cdn") || strings.Contains(value, "cloudfront") {
			templates = append(templates, "cloudfront-api")
		}
		if len(templates) == 1 {
			templates = append(templates, "lambda-api")
		}
		return templates
	}
}
