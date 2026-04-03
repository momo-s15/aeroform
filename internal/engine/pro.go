package engine

import (
	"fmt"
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/llm"
	"github.com/momo-s15/aeroform/internal/prompt"
	"github.com/momo-s15/aeroform/internal/providers"
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

	provider := providers.ForCloud(cfg.Cloud)
	templates := provider.SelectProTemplates(promptText)
	if client != nil {
		allowed := provider.SupportedProTemplates()
		promptBody := prompt.BuildProPrompt(promptText, cfg.Cloud, regionForCloud(cfg), "production", allowed)
		selection, err := llm.GenerateTemplateSelection(client, promptBody, allowed)
		if err == nil && len(selection.Templates) > 0 {
			templates = selection.Templates
		}
	}

	region := provider.Region(cfg)

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
	return providers.ForCloud(cfg.Cloud).Region(cfg)
}
