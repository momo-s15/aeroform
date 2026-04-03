package security

import (
	"bytes"
	"strings"
	"testing"

	"github.com/momo-s15/aeroform/internal/engine"
)

func TestEvaluateSimplePlanAllowsStaticSite(t *testing.T) {
	report := EvaluateSimplePlan(engine.SimpleLaunchPlan{
		Provider:    "aws",
		ProjectName: "my-site",
		Template:    "static-site",
	})
	if report.BlockingCount != 0 {
		t.Fatalf("expected no blocking findings")
	}
	if report.HasFindings {
		t.Fatalf("expected no findings for a clean plan")
	}
}

func TestEvaluateSimplePlanBlocksTinyDB(t *testing.T) {
	report := EvaluateSimplePlan(engine.SimpleLaunchPlan{
		Provider:    "aws",
		ProjectName: "my-database",
		Template:    "tiny-db",
	})
	if report.BlockingCount == 0 {
		t.Fatalf("expected blocking findings")
	}
}

func TestPrintSimpleReport(t *testing.T) {
	report := EvaluateSimplePlan(engine.SimpleLaunchPlan{
		Provider:     "aws",
		ProjectName:  "My Site",
		Template:     "static-site",
		CustomDomain: "example.com",
	})
	var buffer bytes.Buffer
	PrintSimpleReport(&buffer, report)
	output := buffer.String()
	if !strings.Contains(output, "Security check") {
		t.Fatalf("expected security output")
	}
	if !strings.Contains(output, "Auto-corrected") {
		t.Fatalf("expected auto-corrected output")
	}
}
