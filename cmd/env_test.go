package cmd

import (
	"strings"
	"testing"

	"github.com/momo-s15/aeroform/internal/engine"
)

func TestResolveEnvDirExplicit(t *testing.T) {
	origDir := envDir
	defer func() { envDir = origDir }()
	envDir = "/some/explicit/path"

	dir, err := resolveEnvDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != "/some/explicit/path" {
		t.Errorf("dir = %q, want %q", dir, "/some/explicit/path")
	}
}

func TestResolveEnvDirNone(t *testing.T) {
	origDir := envDir
	defer func() { envDir = origDir }()
	envDir = ""

	_, err := resolveEnvDir()
	if err == nil {
		t.Fatal("expected error when no dir set and no .aeroform/pro exists")
	}
	if !strings.Contains(err.Error(), "--dir") {
		t.Errorf("error should mention --dir, got: %v", err)
	}
}

func TestSlugifyRejectsInvalidEnvNames(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"staging", "staging"},
		{"prod", "prod"},
		{"My Env", "my-env"},
		{"!!!", ""},
		{"", ""},
	}
	for _, tc := range cases {
		got := engine.Slugify(tc.input)
		if got != tc.want {
			t.Errorf("Slugify(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestEnvCommandsRegistered(t *testing.T) {
	subs := envCmd.Commands()
	names := make(map[string]bool)
	for _, c := range subs {
		names[c.Name()] = true
	}
	for _, want := range []string{"add", "list", "select"} {
		if !names[want] {
			t.Errorf("env subcommand %q not registered", want)
		}
	}
}
