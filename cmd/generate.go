package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/llm"
	"github.com/momo-s15/aeroform/internal/security"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [prompt]",
	Short: "Pro Mode pipeline: config + plan + apply preview",
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
		fmt.Fprintln(out, "Aeroform Pro Mode generate")
		for _, line := range plan.Summary {
			fmt.Fprintln(out, line)
		}

		scanReport := security.RunProSecurityScan(".")
		for _, line := range scanReport.SummaryLines() {
			fmt.Fprintln(out, line)
		}
		if scanReport.Blocking {
			return fmt.Errorf("security gate blocked generate")
		}

		fmt.Fprintln(out, "apply: deferred in this phase; foundation wiring only")
		return nil
	},
}

func proLLMClientFromConfig(cfg config.Config) llm.Client {
	client := engine.ProLLMClientFromConfig(cfg)
	if client != nil && client.IsAvailable() {
		return client
	}
	return nil
}
