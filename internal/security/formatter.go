package security

import (
	"fmt"
	"io"
)

func PrintSimpleReport(out io.Writer, report Report) {
	if !report.HasFindings {
		fmt.Fprintln(out, "Security check: no issues found")
		return
	}

	fmt.Fprintln(out, "Security check")
	for _, finding := range report.Findings {
		fmt.Fprintf(out, "- [%s] %s\n", finding.Severity, finding.Title)
		fmt.Fprintf(out, "  %s\n", finding.Description)
		if finding.Fix != "" {
			fmt.Fprintf(out, "  Fix: %s\n", finding.Fix)
		}
	}
	if len(report.AutoCorrected) > 0 {
		fmt.Fprintln(out, "Auto-corrected")
		for _, item := range report.AutoCorrected {
			fmt.Fprintf(out, "- %s\n", item)
		}
	}
}
