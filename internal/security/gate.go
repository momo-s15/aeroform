package security

import (
	"fmt"
	"strings"

	"github.com/momo-s15/aeroform/internal/engine"
)

type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

type Finding struct {
	Severity    Severity
	Title       string
	Description string
	Fix         string
	Blocking    bool
}

type Report struct {
	Findings       []Finding
	AutoCorrected  []string
	CorrectedPlan  engine.SimpleLaunchPlan
	BlockingCount  int
	WasAutoFixed   bool
	HasFindings    bool
	PlainTextLines []string
}

func EvaluateSimplePlan(plan engine.SimpleLaunchPlan) Report {
	corrected := plan
	report := Report{CorrectedPlan: corrected}

	if normalized := engine.Slugify(plan.ProjectName); normalized != "" && normalized != plan.ProjectName {
		report.Findings = append(report.Findings, Finding{
			Severity:    SeverityLow,
			Title:       "Project name normalized",
			Description: fmt.Sprintf("Simple Mode keeps project names lowercase and hyphenated. %q will be stored as %q.", plan.ProjectName, normalized),
			Fix:         "Use the normalized project name for consistent file and resource names.",
		})
		report.AutoCorrected = append(report.AutoCorrected, fmt.Sprintf("project name -> %s", normalized))
		corrected.ProjectName = normalized
		report.WasAutoFixed = true
	}

	if plan.Template == "tiny-db" || plan.Template == "fullstack-app" {
		report.Findings = append(report.Findings, Finding{
			Severity:    SeverityLow,
			Title:       "Database security enforced",
			Description: "This template includes a database. Aeroform will auto-correct any insecure settings (public access, unencrypted storage) after rendering.",
			Fix:         "No action needed — autocorrect runs before deployment.",
		})
	}

	if plan.CustomDomain != "" {
		report.Findings = append(report.Findings, Finding{
			Severity:    SeverityLow,
			Title:       "HTTPS will be enforced",
			Description: fmt.Sprintf("Custom domain %q will be routed through HTTPS only.", plan.CustomDomain),
			Fix:         "Aeroform will provision the certificate and redirect HTTP to HTTPS.",
		})
	}

	supportedProviders := map[string]bool{"aws": true, "azure": true, "gcp": true}
	if !supportedProviders[strings.TrimSpace(plan.Provider)] {
		report.Findings = append(report.Findings, Finding{
			Severity:    SeverityHigh,
			Title:       "Unsupported provider for Simple Mode",
			Description: fmt.Sprintf("Simple Mode supports aws, azure, and gcp, but %q was requested.", plan.Provider),
			Fix:         "Choose aws, azure, or gcp.",
			Blocking:    true,
		})
		report.BlockingCount++
	}

	report.CorrectedPlan = corrected
	report.HasFindings = len(report.Findings) > 0
	return report
}
