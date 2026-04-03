package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/spf13/cobra"
)

var driftCmd = &cobra.Command{
	Use:   "drift [prompt]",
	Short: "Generate a Pro Mode drift summary",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadFromFile("config.yaml")
		if err != nil {
			return fmt.Errorf("pro mode requires config.yaml: %w", err)
		}

		report, err := engine.BuildDriftReport(cfg, args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintln(out, "Aeroform Pro Mode drift")
		for _, line := range report.Summary {
			fmt.Fprintln(out, line)
		}
		return nil
	},
}
