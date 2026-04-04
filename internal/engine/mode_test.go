package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectMode(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "config.yaml")

	t.Run("missing file", func(t *testing.T) {
		if got := DetectMode(filepath.Join(dir, "nope.yaml"), Flags{}); got != SimpleMode {
			t.Fatalf("got %v, want SimpleMode", got)
		}
	})

	t.Run("unrelated yaml without cloud", func(t *testing.T) {
		if err := os.WriteFile(cfg, []byte("foo: bar\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := DetectMode(cfg, Flags{}); got != SimpleMode {
			t.Fatalf("unrelated config.yaml must not force ProMode, got %v", got)
		}
	})

	t.Run("aeroform pro by mode field", func(t *testing.T) {
		if err := os.WriteFile(cfg, []byte("mode: pro\ncloud: aws\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := DetectMode(cfg, Flags{}); got != ProMode {
			t.Fatalf("got %v, want ProMode", got)
		}
	})

	t.Run("aeroform simple by mode field", func(t *testing.T) {
		if err := os.WriteFile(cfg, []byte("mode: simple\ncloud: aws\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := DetectMode(cfg, Flags{}); got != SimpleMode {
			t.Fatalf("got %v, want SimpleMode", got)
		}
	})

	t.Run("aeroform pro when cloud set mode omitted", func(t *testing.T) {
		if err := os.WriteFile(cfg, []byte("cloud: azure\nazure:\n  location: eastus\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := DetectMode(cfg, Flags{}); got != ProMode {
			t.Fatalf("got %v, want ProMode", got)
		}
	})

	t.Run("flags simple wins", func(t *testing.T) {
		if err := os.WriteFile(cfg, []byte("cloud: aws\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := DetectMode(cfg, Flags{Simple: true}); got != SimpleMode {
			t.Fatalf("got %v, want SimpleMode", got)
		}
	})
}
