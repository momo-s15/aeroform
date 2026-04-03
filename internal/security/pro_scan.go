package security

import (
	"fmt"
	"os/exec"
	"strings"
)

var execLookPath = exec.LookPath
var execCommand = exec.Command

type ToolScanResult struct {
	Tool      string
	Available bool
	Passed    bool
	ExitCode  int
	Output    string
	Err       error
}

type ProScanReport struct {
	Checkov  ToolScanResult
	Tfsec    ToolScanResult
	Blocking bool
}

func RunProSecurityScan(terraformDir string) ProScanReport {
	checkov := runTool("checkov", []string{"-d", terraformDir})
	tfsec := runTool("tfsec", []string{terraformDir})

	report := ProScanReport{
		Checkov: checkov,
		Tfsec:   tfsec,
	}

	if checkov.Available && !checkov.Passed {
		report.Blocking = true
	}
	if tfsec.Available && !tfsec.Passed {
		report.Blocking = true
	}

	return report
}

func (r ProScanReport) SummaryLines() []string {
	lines := []string{
		formatToolLine(r.Checkov),
		formatToolLine(r.Tfsec),
	}
	if !r.Checkov.Available && !r.Tfsec.Available {
		lines = append(lines, "No scanner binaries found in PATH. Install checkov and tfsec for full Pro Mode gating.")
	}
	if r.Blocking {
		lines = append(lines, "Security gate: blocked")
	} else {
		lines = append(lines, "Security gate: pass")
	}
	return lines
}

func formatToolLine(result ToolScanResult) string {
	if !result.Available {
		return fmt.Sprintf("%s: unavailable", result.Tool)
	}
	if result.Passed {
		return fmt.Sprintf("%s: pass", result.Tool)
	}
	return fmt.Sprintf("%s: fail (exit %d)", result.Tool, result.ExitCode)
}

func runTool(name string, args []string) ToolScanResult {
	result := ToolScanResult{Tool: name}

	if _, err := execLookPath(name); err != nil {
		result.Available = false
		result.Passed = true
		result.Output = "binary not found"
		return result
	}

	result.Available = true
	cmd := execCommand(name, args...)
	output, err := cmd.CombinedOutput()
	result.Output = strings.TrimSpace(string(output))
	if err == nil {
		result.Passed = true
		return result
	}

	result.Passed = false
	result.Err = err
	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
	}
	return result
}
