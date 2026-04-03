package cost

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestEstimateForSimpleTemplate(t *testing.T) {
	estimate := EstimateForSimpleTemplate("aws", "static-site", true, 20)
	if !estimate.Monthly.GreaterThan(decimal.Zero) {
		t.Fatalf("expected a positive estimate")
	}
	if !estimate.HasDomainFee {
		t.Fatalf("expected domain fee to be included")
	}
	if len(estimate.Lines) == 0 {
		t.Fatalf("expected breakdown lines")
	}
}

func TestEstimateForSimpleTemplateExactArithmetic(t *testing.T) {
	est := EstimateForSimpleTemplate("aws", "static-site", true, 20)
	want := decimal.NewFromFloat(1.00)
	if !est.Monthly.Equal(want) {
		t.Fatalf("expected $1.00 exactly, got %s", est.Monthly.StringFixed(2))
	}
}

func TestEstimateForSimpleTemplateAzure(t *testing.T) {
	est := EstimateForSimpleTemplate("azure", "function-api", false, 20)
	if !est.Monthly.Equal(decimal.Zero) {
		t.Fatalf("expected $0.00 for azure function-api, got %s", est.Monthly.StringFixed(2))
	}
}

func TestEstimateForSimpleTemplateGCP(t *testing.T) {
	est := EstimateForSimpleTemplate("gcp", "cloud-run-api", false, 20)
	if !est.Monthly.Equal(decimal.Zero) {
		t.Fatalf("expected $0.00 for gcp cloud-run-api, got %s", est.Monthly.StringFixed(2))
	}
}

func TestSupportedSimpleTemplatesForProvider(t *testing.T) {
	azure := SupportedSimpleTemplatesForProvider("azure")
	if len(azure) != 2 {
		t.Fatalf("expected 2 azure templates, got %d", len(azure))
	}
	gcp := SupportedSimpleTemplatesForProvider("gcp")
	if len(gcp) != 2 {
		t.Fatalf("expected 2 gcp templates, got %d", len(gcp))
	}
}

func TestSupportedSimpleTemplates(t *testing.T) {
	templates := SupportedSimpleTemplates()
	if len(templates) == 0 {
		t.Fatalf("expected templates")
	}
}

func TestEstimateForProTemplatesAWS(t *testing.T) {
	est := EstimateForProTemplates("aws", []string{"vpc", "eks", "rds-private"})
	if est.Cloud != "aws" {
		t.Fatalf("expected cloud=aws, got %q", est.Cloud)
	}
	if !est.Monthly.GreaterThan(decimal.Zero) {
		t.Fatal("expected positive total for vpc+eks+rds-private")
	}
	if len(est.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(est.Items))
	}
	if est.FreeTier {
		t.Fatal("vpc+eks+rds should not be free tier")
	}
}

func TestEstimateForProTemplatesAzure(t *testing.T) {
	est := EstimateForProTemplates("azure", []string{"vnet", "storage"})
	if est.Cloud != "azure" {
		t.Fatalf("expected cloud=azure, got %q", est.Cloud)
	}
	if len(est.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(est.Items))
	}
	if est.Monthly.LessThan(decimal.Zero) {
		t.Fatal("expected non-negative cost")
	}
}

func TestEstimateForProTemplatesGCP(t *testing.T) {
	est := EstimateForProTemplates("gcp", []string{"gcs"})
	if est.Cloud != "gcp" {
		t.Fatalf("expected cloud=gcp, got %q", est.Cloud)
	}
	if len(est.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(est.Items))
	}
	if !est.FreeTier {
		t.Fatal("gcs at $0.50 should be under the $1.00 free-tier threshold")
	}
}

func TestEstimateForProTemplatesFreeTier(t *testing.T) {
	est := EstimateForProTemplates("aws", []string{"lambda-api"})
	if !est.FreeTier {
		t.Fatal("lambda-api alone should be free tier")
	}
}

func TestEstimateForProTemplatesUnknown(t *testing.T) {
	est := EstimateForProTemplates("aws", []string{"mystery-template"})
	if len(est.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(est.Items))
	}
	if !strings.Contains(est.Items[0].Note, "not yet catalogued") {
		t.Fatalf("expected 'not yet catalogued' for unknown, got %q", est.Items[0].Note)
	}
}

func TestProEstimateLines(t *testing.T) {
	est := EstimateForProTemplates("aws", []string{"vpc", "s3-private"})
	lines := est.Lines()
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines (2 items + total), got %d", len(lines))
	}
	found := false
	for _, l := range lines {
		if strings.Contains(l, "TOTAL") {
			found = true
		}
	}
	if !found {
		t.Fatal("expected TOTAL line in output")
	}
}

func TestProEstimateExactArithmetic(t *testing.T) {
	est := EstimateForProTemplates("aws", []string{"s3-private", "s3-private"})
	want := decimal.NewFromFloat(1.00)
	if !est.Monthly.Equal(want) {
		t.Fatalf("expected $1.00 exactly for 2x s3-private, got %s", est.Monthly.StringFixed(2))
	}
}
