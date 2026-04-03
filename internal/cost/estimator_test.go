package cost

import "testing"

func TestEstimateForSimpleTemplate(t *testing.T) {
	estimate := EstimateForSimpleTemplate("static-site", true, 20)
	if estimate.Monthly <= 0 {
		t.Fatalf("expected a positive estimate")
	}
	if !estimate.HasDomainFee {
		t.Fatalf("expected domain fee to be included")
	}
	if len(estimate.Lines) == 0 {
		t.Fatalf("expected breakdown lines")
	}
}

func TestSupportedSimpleTemplates(t *testing.T) {
	templates := SupportedSimpleTemplates()
	if len(templates) == 0 {
		t.Fatalf("expected templates")
	}
}
