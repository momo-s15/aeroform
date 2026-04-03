package engine

import (
	"strings"
	"testing"
	"time"

	"github.com/momo-s15/aeroform/internal/simplestate"
)

func TestBuildProConfigFromSimpleState(t *testing.T) {
	cfg, err := BuildProConfigFromSimpleState(simplestate.State{
		Projects: []simplestate.Project{{
			Name:      "demo",
			Provider:  "aws",
			Template:  "static-site",
			Prompt:    "a portfolio site",
			CreatedAt: time.Now().UTC(),
		}},
	})
	if err != nil {
		t.Fatalf("BuildProConfigFromSimpleState() error = %v", err)
	}
	if cfg.Mode != "pro" || cfg.Cloud != "aws" {
		t.Fatalf("unexpected config %#v", cfg)
	}
}

func TestRenderConfigYAML(t *testing.T) {
	cfg, err := BuildProConfigFromSimpleState(simplestate.State{
		Projects: []simplestate.Project{{Provider: "aws", CreatedAt: time.Now().UTC()}},
	})
	if err != nil {
		t.Fatalf("BuildProConfigFromSimpleState() error = %v", err)
	}

	rendered, err := RenderConfigYAML(cfg)
	if err != nil {
		t.Fatalf("RenderConfigYAML() error = %v", err)
	}
	if !strings.Contains(rendered, "cloud: aws") {
		t.Fatalf("expected cloud field in rendered yaml")
	}
}
