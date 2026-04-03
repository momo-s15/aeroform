package simplestate

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStateRoundTrip(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	project := Project{
		Name:            "demo",
		Provider:        "aws",
		Template:        "static-site",
		Prompt:          "a personal website",
		MonthlyEstimate: 0.5,
		CreatedAt:       time.Now().UTC(),
	}

	if err := AddProject(project); err != nil {
		t.Fatalf("AddProject: %v", err)
	}

	state, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(state.Projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(state.Projects))
	}

	if err := Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	state, err = Load()
	if err != nil {
		t.Fatalf("Load after clear: %v", err)
	}
	if len(state.Projects) != 0 {
		t.Fatalf("expected 0 projects after clear")
	}

	if _, err := os.Stat(filepath.Join(tempDir, statePath)); err != nil {
		t.Fatalf("expected state file: %v", err)
	}
}
