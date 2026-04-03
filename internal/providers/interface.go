package providers

import "github.com/momo-s15/aeroform/internal/config"

type CloudProvider interface {
	Name() string
	Region(cfg config.Config) string
	BootstrapSteps(cfg config.Config) []string
	SupportedProTemplates() []string
	SelectProTemplates(prompt string) []string
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
