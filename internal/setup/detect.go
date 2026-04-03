package setup

import "os/exec"

func OllamaAvailable() bool {
	_, err := exec.LookPath("ollama")
	return err == nil
}

func VerifyOllamaBinary() error {
	_, err := exec.Command("ollama", "--version").CombinedOutput()
	return err
}
