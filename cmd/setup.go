package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/setup"
	"github.com/spf13/cobra"
)

var setupInstall bool

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install and verify the local AI backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		platform := setup.DetectPlatform()

		fmt.Fprintln(out, "Aeroform setup")

		if setup.OllamaAvailable() {
			fmt.Fprintln(out, "✓ Ollama found in PATH")
			if err := setup.VerifyOllamaBinary(); err != nil {
				return fmt.Errorf("ollama is installed but could not be verified: %w", err)
			}
			fmt.Fprintln(out, "✓ Ollama is ready")
			fmt.Fprintln(out, "Next: run `aeroform launch \"...\"` or `aeroform generate \"...\"` once the rest of the CLI is ready.")
			return nil
		}

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

		fmt.Fprintln(out, "✓ Ollama installation step completed")
		fmt.Fprintln(out, "Run `aeroform setup` again to verify the binary and finish startup.")
		return nil
	},
}

func init() {
	setupCmd.Flags().BoolVar(&setupInstall, "install", false, "attempt automatic Ollama installation when missing")
}
