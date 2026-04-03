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

	if plan.Template == "tiny-db" {
		report.Findings = append(report.Findings, Finding{
			Severity:    SeverityHigh,
			Title:       "Public database risk",
			Description: "Simple Mode never deploys a public database. A database must remain private, encrypted, and behind a controlled access path.",
			Fix:         "Use a private database pattern or switch to a serverless data store for Simple Mode.",
			Blocking:    true,
		})
		report.BlockingCount++
	}

	if plan.CustomDomain != "" {
		report.Findings = append(report.Findings, Finding{
			Severity:    SeverityLow,
			Title:       "HTTPS will be enforced",
			Description: fmt.Sprintf("Custom domain %q will be routed through HTTPS only.", plan.CustomDomain),
			Fix:         "Aeroform will provision the certificate and redirect HTTP to HTTPS.",
		})
	}

	if strings.TrimSpace(plan.Provider) != "aws" {
		report.Findings = append(report.Findings, Finding{
			Severity:    SeverityHigh,
			Title:       "Unsupported provider for Simple Mode",
			Description: fmt.Sprintf("Simple Mode currently supports AWS only, but %q was requested.", plan.Provider),
			Fix:         "Choose AWS for the first Simple Mode release.",
			Blocking:    true,
		})
		report.BlockingCount++
	}

	report.CorrectedPlan = corrected
	report.HasFindings = len(report.Findings) > 0
	return report
}
