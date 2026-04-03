package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/bootstrap"
	"github.com/momo-s15/aeroform/internal/config"
	"github.com/spf13/cobra"
)

var bootstrapRepo string

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Generate OIDC trust + remote state setup commands for your cloud",
	Long: `Bootstrap generates the CLI commands needed to set up:
  - OIDC trust between GitHub Actions and your cloud provider
  - A remote state backend (S3/Azure Blob/GCS) for Terraform

Review the output and run the commands in your terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadFromFile("config.yaml")
		if err != nil {
			return fmt.Errorf("pro mode requires config.yaml: %w", err)
		}

		out := cmd.OutOrStdout()
		report := bootstrap.Run(cfg, bootstrap.Params{Repo: bootstrapRepo})
		fmt.Fprint(out, bootstrap.Render(report))
		return nil
	},
}

func init() {
	bootstrapCmd.Flags().StringVar(&bootstrapRepo, "repo", "", "GitHub repository (owner/repo) for OIDC trust policy")
}
