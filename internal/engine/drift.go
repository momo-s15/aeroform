package engine

import (
	"fmt"
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
)

type DriftReport struct {
	Summary []string
}

func BuildDriftReport(cfg config.Config, prompt string) (DriftReport, error) {
	plan, err := BuildProPlan(cfg, prompt)
	if err != nil {
		return DriftReport{}, err
	}

	return DriftReport{Summary: []string{
		fmt.Sprintf("cloud: %s", cfg.Cloud),
		fmt.Sprintf("desired templates: %s", strings.Join(plan.Templates, ", ")),
		"drift status: baseline only in this phase (full live-state comparison arrives in a later phase)",
	}}, nil
}
