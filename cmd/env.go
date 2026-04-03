package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/terraform"
	"github.com/momo-s15/aeroform/internal/ui"
	"github.com/spf13/cobra"
)

var envDir string

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment workspaces (staging, prod, etc.)",
	Long: `Manage Terraform workspaces that isolate state for different environments.

Each workspace gets its own state file, so you can deploy staging and prod
side by side from the same templates.

Use --dir to point at the project directory that was created by 'aeroform generate'.`,
}

var envAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Create a new Terraform workspace for an environment",
	Long: `Create a new Terraform workspace to isolate state for a named environment.

Example:
  aeroform env add staging --dir .aeroform/pro/my-app
  aeroform env add prod`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envName := engine.Slugify(args[0])
		if envName == "" {
			return fmt.Errorf("invalid environment name %q", args[0])
		}

		workDir, err := resolveEnvDir()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Creating environment %q in %s\n", envName, workDir)

		if err := terraform.WorkspaceNew(workDir, envName); err != nil {
			return fmt.Errorf("create workspace: %w", err)
		}

		ui.Successln(out, "✓ Environment "+envName+" created")
		fmt.Fprintf(out, "  Terraform workspace %q is now active.\n", envName)
		fmt.Fprintf(out, "  Run 'terraform apply' in %s to deploy this environment.\n", workDir)
		return nil
	},
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Terraform workspaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		workDir, err := resolveEnvDir()
		if err != nil {
			return err
		}

		workspaces, current, err := terraform.WorkspaceList(workDir)
		if err != nil {
			fmt.Fprintln(out, "No environments found. Run 'aeroform env add <name>' first.")
			return nil
		}

		ui.Boldln(out, "Environments ("+workDir+")")
		for _, ws := range workspaces {
			if ws == current {
				ui.Success(out, "  * %s (active)\n", ws)
			} else {
				fmt.Fprintf(out, "    %s\n", ws)
			}
		}
		return nil
	},
}

var envSelectCmd = &cobra.Command{
	Use:   "select [name]",
	Short: "Switch to a Terraform workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envName := engine.Slugify(args[0])
		if envName == "" {
			return fmt.Errorf("invalid environment name %q", args[0])
		}

		workDir, err := resolveEnvDir()
		if err != nil {
			return err
		}

		if err := terraform.WorkspaceSelect(workDir, envName); err != nil {
			return fmt.Errorf("select workspace: %w", err)
		}

		out := cmd.OutOrStdout()
		ui.Successln(out, "✓ Switched to environment "+envName)
		return nil
	},
}

func resolveEnvDir() (string, error) {
	if envDir != "" {
		return envDir, nil
	}
	entries, err := os.ReadDir(filepath.Join(".aeroform", "pro"))
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				return filepath.Join(".aeroform", "pro", e.Name()), nil
			}
		}
	}
	return "", fmt.Errorf("no project directory found; use --dir to specify one (e.g. --dir .aeroform/pro/my-app)")
}

func init() {
	envCmd.PersistentFlags().StringVar(&envDir, "dir", "", "path to the Terraform project directory")
	envCmd.AddCommand(envAddCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envSelectCmd)
}
