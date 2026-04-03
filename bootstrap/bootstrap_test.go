package bootstrap

import (
	"strings"
	"testing"

	"github.com/momo-s15/aeroform/internal/config"
)

func TestRunAWS(t *testing.T) {
	report := Run(config.Config{Cloud: "aws", AWS: config.AWSConfig{Region: "us-east-1"}})
	if report.Cloud != "aws" {
		t.Fatalf("expected aws cloud")
	}
	if len(report.Steps) == 0 {
		t.Fatalf("expected bootstrap steps")
	}
}

func TestRender(t *testing.T) {
	text := Render(Report{Cloud: "aws", Steps: []string{"one", "two"}})
	if !strings.Contains(text, "one") || !strings.Contains(text, "two") {
		t.Fatalf("expected rendered steps")
	}
}
