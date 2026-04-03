package security

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestRunProSecurityScanUnavailable(t *testing.T) {
	originalLookPath := execLookPath
	originalCommand := execCommand
	t.Cleanup(func() {
		execLookPath = originalLookPath
		execCommand = originalCommand
	})

	execLookPath = func(file string) (string, error) {
		return "", fmt.Errorf("missing")
	}

	report := RunProSecurityScan(".")
	if report.Blocking {
		t.Fatalf("expected non-blocking report when scanners are unavailable")
	}
	if report.Checkov.Available || report.Tfsec.Available {
		t.Fatalf("expected scanners to be unavailable")
	}
}

func TestRunProSecurityScanBlocksOnFailure(t *testing.T) {
	originalLookPath := execLookPath
	originalCommand := execCommand
	t.Cleanup(func() {
		execLookPath = originalLookPath
		execCommand = originalCommand
	})

	execLookPath = func(file string) (string, error) {
		return file, nil
	}

	execCommand = func(name string, args ...string) *exec.Cmd {
		allArgs := []string{"-test.run=TestHelperProcess", "--", name}
		allArgs = append(allArgs, args...)
		cmd := exec.Command(os.Args[0], allArgs...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		return cmd
	}

	report := RunProSecurityScan(".")
	if !report.Blocking {
		t.Fatalf("expected blocking report on scanner failure")
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(1)
}
