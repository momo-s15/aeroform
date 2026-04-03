package security

import (
	"fmt"
	"io"

	"github.com/momo-s15/aeroform/internal/ui"
)

func PrintSimpleReport(out io.Writer, report Report) {
	if !report.HasFindings {
		ui.Successln(out, "✓ Security check: no issues found")
		return
	}

	ui.Boldln(out, "Security check")
	for _, finding := range report.Findings {
		switch finding.Severity {
		case SeverityHigh:
			ui.Error(out, "  ✗ [%s] %s\n", finding.Severity, finding.Title)
		case SeverityMedium:
			ui.Warn(out, "  ⚠ [%s] %s\n", finding.Severity, finding.Title)
		default:
			ui.Info(out, "  ℹ [%s] %s\n", finding.Severity, finding.Title)
		}
		fmt.Fprintf(out, "    %s\n", finding.Description)
		if finding.Fix != "" {
			fmt.Fprintf(out, "    Fix: %s\n", finding.Fix)
		}
	}
	if len(report.AutoCorrected) > 0 {
		ui.Successln(out, "Auto-corrected")
		for _, item := range report.AutoCorrected {
			ui.Success(out, "  ✓ %s\n", item)
		}
	}
}
