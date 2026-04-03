package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/spf13/cobra"
)

var driftCmd = &cobra.Command{
	Use:   "drift [prompt]",
	Short: "Detect infrastructure drift by comparing live state to desired config",
	Long: `Run a terraform plan against an existing Aeroform project directory to detect
resources that have drifted from the desired state. The prompt should match the
one used when the infrastructure was originally deployed.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadFromFile("config.yaml")
		if err != nil {
			return fmt.Errorf("pro mode requires config.yaml: %w", err)
		}

		out := cmd.OutOrStdout()
		fmt.Fprintln(out, "Aeroform drift detection")
		fmt.Fprintln(out, "")

		report, err := engine.BuildDriftReport(cfg, args[0])
		if err != nil {
			return err
		}

		for _, line := range report.Summary {
			fmt.Fprintln(out, "  "+line)
		}

		if !report.HasDrift {
			fmt.Fprintln(out, "")
			fmt.Fprintln(out, "All clear — your infrastructure matches the desired state.")
		}

		return nil
	},
}
