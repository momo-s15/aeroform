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
