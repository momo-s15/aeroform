package engine

import (
	"fmt"
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/simplestate"
	"gopkg.in/yaml.v3"
)

func BuildProConfigFromSimpleState(st simplestate.State) (config.Config, error) {
	if len(st.Projects) == 0 {
		return config.Config{}, fmt.Errorf("no simple mode projects found")
	}

	latest := st.Projects[len(st.Projects)-1]
	cloud := strings.TrimSpace(strings.ToLower(latest.Provider))
	if cloud == "" {
		cloud = "aws"
	}

	cfg := config.Config{
		Cloud: cloud,
		Mode:  "pro",
		State: config.StateConfig{Backend: "s3"},
	}

	switch cloud {
	case "aws":
		cfg.AWS = config.AWSConfig{Region: "us-east-1", AccountID: "000000000000"}
	case "azure":
		cfg.Azure = config.AzureConfig{Location: "eastus"}
	case "gcp":
		cfg.GCP = config.GCPConfig{Region: "us-central1"}
	default:
		cfg.Cloud = "aws"
		cfg.AWS = config.AWSConfig{Region: "us-east-1", AccountID: "000000000000"}
	}

	return cfg, nil
}

func RenderConfigYAML(cfg config.Config) (string, error) {
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
