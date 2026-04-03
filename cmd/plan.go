package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan [prompt]",
	Short: "Pro Mode dry-run planning",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadFromFile("config.yaml")
		if err != nil {
			return fmt.Errorf("pro mode requires config.yaml: %w", err)
		}

		plan, err := engine.BuildProPlanWithClient(cfg, args[0], proLLMClientFromConfig(cfg))
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintln(out, "Aeroform Pro Mode plan")
		for _, line := range plan.Summary {
			fmt.Fprintln(out, line)
		}
		fmt.Fprintln(out, "plan-only mode: no infrastructure changes made")
		return nil
	},
}
