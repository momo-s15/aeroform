package setup

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDetectPlatformMatchesRuntime(t *testing.T) {
	if got := DetectPlatform(); got != runtime.GOOS {
		t.Fatalf("DetectPlatform() = %q, want %q", got, runtime.GOOS)
	}
}

func TestInstallInstructions(t *testing.T) {
	tests := []struct {
		platform string
		want     string
	}{
		{platform: "darwin", want: "brew install ollama"},
		{platform: "linux", want: "install.sh"},
		{platform: "windows", want: "download/windows"},
	}

	for _, testCase := range tests {
		t.Run(testCase.platform, func(t *testing.T) {
			got := InstallInstructions(testCase.platform)
			if !strings.Contains(got, testCase.want) {
				t.Fatalf("InstallInstructions(%q) = %q, want substring %q", testCase.platform, got, testCase.want)
			}
		})
	}
}

func TestDefaultModelConstant(t *testing.T) {
	if DefaultModel == "" {
		t.Fatal("DefaultModel should not be empty")
	}
	if DefaultModel != "llama3.2" {
		t.Fatalf("DefaultModel = %q, want %q", DefaultModel, "llama3.2")
	}
}

func TestPullModelReturnsErrorWhenOllamaMissing(t *testing.T) {
	originalPath := os.Getenv("PATH")
	t.Setenv("PATH", t.TempDir())
	t.Cleanup(func() {
		_ = os.Setenv("PATH", originalPath)
	})

	var buf strings.Builder
	err := PullModel("testmodel", &buf)
	if err == nil {
		t.Fatal("expected error when ollama is not in PATH")
	}
}

func TestOllamaAvailableWithFakeBinary(t *testing.T) {
	tempDir := t.TempDir()
	binaryName := "ollama"
	if runtime.GOOS == "windows" {
		binaryName = "ollama.exe"
	}

	binaryPath := filepath.Join(tempDir, binaryName)
	if err := os.WriteFile(binaryPath, []byte("fake ollama"), 0o600); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}
	// Unix exec.LookPath ignores non-executable files; Windows does not require +x.
	if runtime.GOOS != "windows" {
		if err := os.Chmod(binaryPath, 0o755); err != nil {
			t.Fatalf("chmod fake binary: %v", err)
		}
	}

	originalPath := os.Getenv("PATH")
	if err := os.Setenv("PATH", tempDir+string(os.PathListSeparator)+originalPath); err != nil {
		t.Fatalf("set PATH: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Setenv("PATH", originalPath)
	})

	if !OllamaAvailable() {
		t.Fatalf("OllamaAvailable() = false, want true")
	}
}
