package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/setup"
	"github.com/spf13/cobra"
)

var (
	setupInstall bool
	setupModel   string
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install and verify the local AI backend, then pull the default model",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		platform := setup.DetectPlatform()

		fmt.Fprintln(out, "Aeroform setup")

		if setup.OllamaAvailable() {
			fmt.Fprintln(out, "  Ollama found in PATH")
			if err := setup.VerifyOllamaBinary(); err != nil {
				return fmt.Errorf("ollama is installed but could not be verified: %w", err)
			}
			fmt.Fprintln(out, "  Ollama is ready")
		} else {
			fmt.Fprintf(out, "Ollama is not installed for %s.\n", platform)
			fmt.Fprintln(out, setup.InstallInstructions(platform))

			if !setupInstall {
				fmt.Fprintln(out, "Run again with --install to attempt automatic installation where supported.")
				return nil
			}

			fmt.Fprintln(out, "Attempting automatic installation...")
			if err := setup.InstallOllamaForPlatform(platform); err != nil {
				return err
			}
			fmt.Fprintln(out, "  Ollama installation step completed")
		}

		model := setupModel
		if model == "" {
			model = setup.DefaultModel
		}
		fmt.Fprintf(out, "-> Pulling model %s (this may take a few minutes on first run)...\n", model)
		if err := setup.PullModel(model, out); err != nil {
			return err
		}
		fmt.Fprintf(out, "  Model %s is ready\n", model)

		fmt.Fprintln(out, "")
		fmt.Fprintln(out, "Setup complete! Run `aeroform launch` or `aeroform generate` to get started.")
		return nil
	},
}

func init() {
	setupCmd.Flags().BoolVar(&setupInstall, "install", false, "attempt automatic Ollama installation when missing")
	setupCmd.Flags().StringVar(&setupModel, "model", "", "model to pull (default: "+setup.DefaultModel+")")
}
