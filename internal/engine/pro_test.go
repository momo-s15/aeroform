package engine

import (
	"testing"

	"github.com/momo-s15/aeroform/internal/config"
)

func TestBuildProPlan(t *testing.T) {
	plan, err := BuildProPlan(config.Config{
		Cloud: "aws",
		AWS:   config.AWSConfig{Region: "us-east-1"},
	}, "resilient eks cluster with private rds and alb")
	if err != nil {
		t.Fatalf("BuildProPlan() error = %v", err)
	}
	if len(plan.Templates) < 3 {
		t.Fatalf("expected multi-template pro plan, got %#v", plan.Templates)
	}
}
