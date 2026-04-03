package providers

import (
	"testing"

	"github.com/momo-s15/aeroform/internal/config"
)

func TestForCloud(t *testing.T) {
	tests := []struct {
		cloud string
		want  string
	}{
		{cloud: "aws", want: "aws"},
		{cloud: "azure", want: "azure"},
		{cloud: "gcp", want: "gcp"},
		{cloud: "unknown", want: "aws"},
	}

	for _, testCase := range tests {
		t.Run(testCase.cloud, func(t *testing.T) {
			provider := ForCloud(testCase.cloud)
			if provider.Name() != testCase.want {
				t.Fatalf("ForCloud(%q) = %q, want %q", testCase.cloud, provider.Name(), testCase.want)
			}
		})
	}
}

func TestProviderTemplates(t *testing.T) {
	cfg := config.Config{Cloud: "azure", Azure: config.AzureConfig{Location: "eastus"}}
	provider := ForCloud(cfg.Cloud)

	templates := provider.SelectProTemplates("kubernetes database app")
	if len(templates) == 0 {
		t.Fatalf("expected templates")
	}
	if provider.Region(cfg) == "" {
		t.Fatalf("expected region")
	}
}
