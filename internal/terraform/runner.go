package terraform

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/momo-s15/aeroform/internal/logger"
)

// Shims for testing — same pattern as internal/security/pro_scan.go.
var execLookPath = exec.LookPath
var execCommand = exec.Command

// PlanResult holds parsed output from terraform plan.
type PlanResult struct {
	AddCount     int
	ChangeCount  int
	DestroyCount int
	RawOutput    string
}

// EnsureBinary checks that the terraform binary is in PATH.
func EnsureBinary() error {
	if _, err := execLookPath("terraform"); err != nil {
		return fmt.Errorf("terraform is not installed or not in PATH: %w\nInstall it from https://developer.hashicorp.com/terraform/install", err)
	}
	return nil
}

// Init runs terraform init in the given working directory.
func Init(workDir string) error {
	if err := EnsureBinary(); err != nil {
		return err
	}
	logger.L().Debugw("terraform init", "workDir", workDir)
	cmd := execCommand("terraform", "init", "-input=false")
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform init failed:\n%s\n%w", string(output), err)
	}
	return nil
}

// Plan runs terraform plan and parses the resource change summary.
func Plan(workDir string) (*PlanResult, error) {
	if err := EnsureBinary(); err != nil {
		return nil, err
	}
	logger.L().Debugw("terraform plan", "workDir", workDir)
	cmd := execCommand("terraform", "plan", "-input=false", "-no-color")
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	raw := string(output)

	result := &PlanResult{RawOutput: raw}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
			// Exit code 2 means plan succeeded but there are changes.
			parsePlanCounts(raw, result)
			return result, nil
		}
		return nil, fmt.Errorf("terraform plan failed:\n%s\n%w", raw, err)
	}

	parsePlanCounts(raw, result)
	return result, nil
}

// Apply runs terraform apply in the given working directory.
func Apply(workDir string, autoApprove bool) error {
	if err := EnsureBinary(); err != nil {
		return err
	}
	args := []string{"apply", "-input=false"}
	if autoApprove {
		args = append(args, "-auto-approve")
	}
	cmd := execCommand("terraform", args...)
	cmd.Dir = workDir
	cmd.Stdout = nil
	cmd.Stderr = nil
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform apply failed:\n%s\n%w", string(output), err)
	}
	return nil
}

// Destroy runs terraform destroy in the given working directory.
func Destroy(workDir string, autoApprove bool) error {
	if err := EnsureBinary(); err != nil {
		return err
	}
	args := []string{"destroy", "-input=false"}
	if autoApprove {
		args = append(args, "-auto-approve")
	}
	cmd := execCommand("terraform", args...)
	cmd.Dir = workDir
	cmd.Stdout = nil
	cmd.Stderr = nil
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform destroy failed:\n%s\n%w", string(output), err)
	}
	return nil
}

// Output runs terraform output and returns the raw text.
func Output(workDir string) (string, error) {
	if err := EnsureBinary(); err != nil {
		return "", err
	}
	cmd := execCommand("terraform", "output", "-no-color")
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("terraform output failed:\n%s\n%w", string(output), err)
	}
	return strings.TrimSpace(string(output)), nil
}

// DriftResult holds the output of a detailed-exitcode plan used for drift detection.
type DriftResult struct {
	HasDrift bool
	Changes  []DriftChange
	Plan     PlanResult
}

// DriftChange represents a single resource that drifted.
type DriftChange struct {
	Address    string
	ChangeType string // "create", "update", "destroy", "replace"
}

// PlanDetailed runs terraform plan with -detailed-exitcode for drift detection.
// Exit 0 = no changes, exit 2 = changes detected, exit 1 = error.
func PlanDetailed(workDir string) (*DriftResult, error) {
	if err := EnsureBinary(); err != nil {
		return nil, err
	}
	cmd := execCommand("terraform", "plan", "-input=false", "-no-color", "-detailed-exitcode")
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	raw := string(output)

	result := &DriftResult{Plan: PlanResult{RawOutput: raw}}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
			result.HasDrift = true
			parsePlanCounts(raw, &result.Plan)
			result.Changes = parseDriftChanges(raw)
			return result, nil
		}
		return nil, fmt.Errorf("terraform plan failed:\n%s\n%w", raw, err)
	}

	parsePlanCounts(raw, &result.Plan)
	return result, nil
}

// planSummaryRe matches the Terraform plan summary line:
// "Plan: X to add, Y to change, Z to destroy."
var planSummaryRe = regexp.MustCompile(
	`Plan:\s+(\d+)\s+to add,\s+(\d+)\s+to change,\s+(\d+)\s+to destroy`,
)

func parsePlanCounts(output string, result *PlanResult) {
	matches := planSummaryRe.FindStringSubmatch(output)
	if len(matches) == 4 {
		result.AddCount, _ = strconv.Atoi(matches[1])
		result.ChangeCount, _ = strconv.Atoi(matches[2])
		result.DestroyCount, _ = strconv.Atoi(matches[3])
	}
}

var driftChangeRe = regexp.MustCompile(
	`#\s+(\S+)\s+(?:will be|must be)\s+(\S+)`,
)

func parseDriftChanges(output string) []DriftChange {
	var changes []DriftChange
	seen := map[string]bool{}
	for _, match := range driftChangeRe.FindAllStringSubmatch(output, -1) {
		addr := match[1]
		if seen[addr] {
			continue
		}
		seen[addr] = true
		action := normalizeChangeType(match[2])
		changes = append(changes, DriftChange{Address: addr, ChangeType: action})
	}
	return changes
}

func normalizeChangeType(raw string) string {
	switch raw {
	case "created":
		return "create"
	case "destroyed":
		return "destroy"
	case "updated", "updated-in-place", "changed":
		return "update"
	case "replaced":
		return "replace"
	default:
		return raw
	}
}
