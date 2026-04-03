package providers

import (
	"io"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/cost"
	"github.com/momo-s15/aeroform/internal/engine"
)

type CostEstimate = cost.ProEstimate

type BootstrapSection struct {
	Title    string
	Commands []string
}

type CloudProvider interface {
	Name() string
	Validate(cfg config.Config) error
	GenerateVars(cfg config.Config, params map[string]string) map[string]string
	GetTemplateDir(mode engine.Mode, tmpl string) string
	EstimateCost(templates []string) CostEstimate
	PostDeploy(cfg config.Config, w io.Writer) error
	BootstrapSections(cfg config.Config, repo string) []BootstrapSection
}

func ForCloud(cloud string) CloudProvider {
	switch cloud {
	case "azure":
		return AzureProvider{}
	case "gcp":
		return GCPProvider{}
	default:
		return AWSProvider{}
	}
}
