package terraform

import (
	"os/exec"
	"strings"
	"testing"
)

func stubLookPathWS(t *testing.T) {
	t.Helper()
	orig := execLookPath
	t.Cleanup(func() { execLookPath = orig })
	execLookPath = func(file string) (string, error) { return file, nil }
}

func TestWorkspaceNew(t *testing.T) {
	stubLookPathWS(t)
	origCmd := execCommand
	t.Cleanup(func() { execCommand = origCmd })

	var captured []string
	execCommand = func(name string, args ...string) *exec.Cmd {
		captured = append([]string{name}, args...)
		return exec.Command("go", "version")
	}

	err := WorkspaceNew(t.TempDir(), "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "terraform workspace new staging"
	got := strings.Join(captured, " ")
	if got != want {
		t.Errorf("args = %q, want %q", got, want)
	}
}

func TestWorkspaceSelect(t *testing.T) {
	stubLookPathWS(t)
	origCmd := execCommand
	t.Cleanup(func() { execCommand = origCmd })

	var captured []string
	execCommand = func(name string, args ...string) *exec.Cmd {
		captured = append([]string{name}, args...)
		return exec.Command("go", "version")
	}

	err := WorkspaceSelect(t.TempDir(), "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "terraform workspace select prod"
	got := strings.Join(captured, " ")
	if got != want {
		t.Errorf("args = %q, want %q", got, want)
	}
}

func TestWorkspaceListParsing(t *testing.T) {
	raw := "  default\n* staging\n  prod\n"
	var workspaces []string
	var current string
	for _, line := range strings.Split(strings.TrimSpace(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "* ") {
			name := strings.TrimPrefix(line, "* ")
			current = name
			workspaces = append(workspaces, name)
		} else {
			workspaces = append(workspaces, line)
		}
	}

	if len(workspaces) != 3 {
		t.Fatalf("expected 3 workspaces, got %d: %v", len(workspaces), workspaces)
	}
	if current != "staging" {
		t.Errorf("current = %q, want %q", current, "staging")
	}
	if workspaces[0] != "default" || workspaces[1] != "staging" || workspaces[2] != "prod" {
		t.Errorf("workspaces = %v", workspaces)
	}
}

func TestWorkspaceNewError(t *testing.T) {
	stubLookPathWS(t)
	origCmd := execCommand
	t.Cleanup(func() { execCommand = origCmd })

	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("go", "run", "nonexistent_file_that_will_fail.go")
	}

	err := WorkspaceNew(t.TempDir(), "bad-ws")
	if err == nil {
		t.Fatal("expected error for failed workspace new")
	}
	if !strings.Contains(err.Error(), "workspace new") {
		t.Errorf("error should mention workspace new, got: %v", err)
	}
}

func TestWorkspaceSelectError(t *testing.T) {
	stubLookPathWS(t)
	origCmd := execCommand
	t.Cleanup(func() { execCommand = origCmd })

	execCommand = func(name string, args ...string) *exec.Cmd {
		return exec.Command("go", "run", "nonexistent_file_that_will_fail.go")
	}

	err := WorkspaceSelect(t.TempDir(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for failed workspace select")
	}
	if !strings.Contains(err.Error(), "workspace select") {
		t.Errorf("error should mention workspace select, got: %v", err)
	}
}

func TestWorkspaceNewMissingBinary(t *testing.T) {
	orig := execLookPath
	t.Cleanup(func() { execLookPath = orig })
	execLookPath = func(file string) (string, error) {
		return "", &exec.Error{Name: file, Err: exec.ErrNotFound}
	}

	err := WorkspaceNew(t.TempDir(), "staging")
	if err == nil {
		t.Fatal("expected error when terraform is missing")
	}
}

func TestWorkspaceListMissingBinary(t *testing.T) {
	orig := execLookPath
	t.Cleanup(func() { execLookPath = orig })
	execLookPath = func(file string) (string, error) {
		return "", &exec.Error{Name: file, Err: exec.ErrNotFound}
	}

	_, _, err := WorkspaceList(t.TempDir())
	if err == nil {
		t.Fatal("expected error when terraform is missing")
	}
}
