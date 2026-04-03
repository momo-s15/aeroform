package setup

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

const DefaultModel = "llama3.2"

func InstallOllama() error {
	return InstallOllamaForPlatform(runtime.GOOS)
}

func InstallOllamaForPlatform(platform string) error {
	switch platform {
	case "darwin":
		if _, err := exec.LookPath("brew"); err == nil {
			cmd := exec.Command("brew", "install", "ollama")
			cmd.Stdout = nil
			cmd.Stderr = nil
			return cmd.Run()
		}
		return fmt.Errorf("Homebrew is required for automatic Ollama installation on macOS")
	case "linux":
		cmd := exec.Command("sh", "-c", "curl -fsSL https://ollama.com/install.sh | sh")
		cmd.Stdout = nil
		cmd.Stderr = nil
		return cmd.Run()
	case "windows":
		return fmt.Errorf("automatic Ollama installation is not available on Windows yet; open https://ollama.com/download/windows")
	default:
		return fmt.Errorf("unsupported platform %q", platform)
	}
}

// PullModel runs `ollama pull <model>` and streams progress to the given writer.
// If w is nil, output goes to os.Stdout.
func PullModel(model string, w io.Writer) error {
	if model == "" {
		model = DefaultModel
	}
	if w == nil {
		w = os.Stdout
	}

	cmd := exec.Command("ollama", "pull", model)
	cmd.Stdout = w
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ollama pull %s: %w", model, err)
	}
	return nil
}

func InstallInstructions(platform string) string {
	switch platform {
	case "darwin":
		return "Install Ollama on macOS with: brew install ollama\nIf you do not have Homebrew, download it from https://brew.sh/"
	case "linux":
		return "Install Ollama on Linux with: curl -fsSL https://ollama.com/install.sh | sh"
	case "windows":
		return "Install Ollama on Windows from: https://ollama.com/download/windows"
	default:
		return "Install Ollama from: https://ollama.com/download"
	}
}
