package config

import "fmt"

func Validate(cfg Config) error {
	if cfg.Cloud == "" {
		return fmt.Errorf("cloud is required")
	}
	if cfg.Mode != "" && cfg.Mode != "pro" && cfg.Mode != "simple" {
		return fmt.Errorf("mode must be one of: simple, pro")
	}

	switch cfg.Cloud {
	case "aws":
		if cfg.AWS.Region == "" {
			return fmt.Errorf("aws.region is required when cloud=aws")
		}
	case "azure":
		if cfg.Azure.Location == "" {
			return fmt.Errorf("azure.location is required when cloud=azure")
		}
	case "gcp":
		if cfg.GCP.Region == "" {
			return fmt.Errorf("gcp.region is required when cloud=gcp")
		}
	default:
		return fmt.Errorf("unsupported cloud %q", cfg.Cloud)
	}
	return nil
}
