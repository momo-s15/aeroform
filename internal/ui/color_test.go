package ui

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestSuccessContainsText(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	Success(&buf, "deploy %s", "done")
	if !strings.Contains(buf.String(), "deploy done") {
		t.Fatalf("expected 'deploy done', got %q", buf.String())
	}
}

func TestErrorContainsText(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	Errorln(&buf, "blocked")
	if !strings.Contains(buf.String(), "blocked") {
		t.Fatalf("expected 'blocked', got %q", buf.String())
	}
}

func TestWarnContainsText(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	Warn(&buf, "over budget: $%d", 50)
	if !strings.Contains(buf.String(), "over budget: $50") {
		t.Fatalf("expected 'over budget: $50', got %q", buf.String())
	}
}

func TestInfoContainsText(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	Infoln(&buf, "scanning")
	if !strings.Contains(buf.String(), "scanning") {
		t.Fatalf("expected 'scanning', got %q", buf.String())
	}
}

func TestBoldContainsText(t *testing.T) {
	color.NoColor = true
	defer func() { color.NoColor = false }()

	var buf bytes.Buffer
	Boldln(&buf, "heading")
	if !strings.Contains(buf.String(), "heading") {
		t.Fatalf("expected 'heading', got %q", buf.String())
	}
}

func TestNoColorEnvDisablesColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	orig := color.NoColor
	defer func() { color.NoColor = orig }()

	if _, ok := os.LookupEnv("NO_COLOR"); !ok {
		t.Fatal("expected NO_COLOR to be set")
	}
}
