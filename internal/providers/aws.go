package providers

import (
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
)

type AWSProvider struct{}

func (AWSProvider) Name() string { return "aws" }

func (AWSProvider) Region(cfg config.Config) string {
	return cfg.AWS.Region
}

func (AWSProvider) BootstrapSteps(cfg config.Config) []string {
	steps := []string{
		"Create or validate OIDC role trust policy",
		"Validate Terraform state backend permissions",
	}
	if cfg.State.Backend == "s3" {
		steps = append(steps, "Ensure S3 bucket and lock table exist for Terraform state")
	}
	return steps
}

func (AWSProvider) SupportedProTemplates() []string {
	return []string{"vpc", "eks", "rds-private", "alb", "lambda-api", "ecs-fargate", "cloudfront-api"}
}

func (AWSProvider) SelectProTemplates(prompt string) []string {
	value := strings.ToLower(prompt)
	templates := []string{"vpc"}

	if strings.Contains(value, "eks") || strings.Contains(value, "kubernetes") {
		templates = append(templates, "eks")
	}
	if strings.Contains(value, "rds") || strings.Contains(value, "database") {
		templates = append(templates, "rds-private")
	}
	if strings.Contains(value, "alb") || strings.Contains(value, "load balancer") {
		templates = append(templates, "alb")
	}
	if len(templates) == 1 {
		templates = append(templates, "lambda-api")
	}

	return templates
}
