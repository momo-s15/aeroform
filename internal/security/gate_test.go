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

func TestEvaluateSimplePlanAllowsTinyDB(t *testing.T) {
	report := EvaluateSimplePlan(engine.SimpleLaunchPlan{
		Provider:    "aws",
		ProjectName: "my-database",
		Template:    "tiny-db",
	})
	if report.BlockingCount != 0 {
		t.Fatalf("tiny-db should not be blocked (autocorrect handles file-level safety)")
	}
	if !report.HasFindings {
		t.Fatal("expected informational finding about database security enforcement")
	}
}

func TestEvaluateSimplePlanAllowsAzure(t *testing.T) {
	report := EvaluateSimplePlan(engine.SimpleLaunchPlan{
		Provider:    "azure",
		ProjectName: "my-site",
		Template:    "static-site",
	})
	if report.BlockingCount != 0 {
		t.Fatalf("expected no blocking findings for azure, got %d", report.BlockingCount)
	}
}

func TestEvaluateSimplePlanAllowsGCP(t *testing.T) {
	report := EvaluateSimplePlan(engine.SimpleLaunchPlan{
		Provider:    "gcp",
		ProjectName: "my-site",
		Template:    "static-site",
	})
	if report.BlockingCount != 0 {
		t.Fatalf("expected no blocking findings for gcp, got %d", report.BlockingCount)
	}
}

func TestEvaluateSimplePlanBlocksUnsupported(t *testing.T) {
	report := EvaluateSimplePlan(engine.SimpleLaunchPlan{
		Provider:    "digitalocean",
		ProjectName: "my-site",
		Template:    "static-site",
	})
	if report.BlockingCount == 0 {
		t.Fatal("expected blocking finding for unsupported provider")
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
