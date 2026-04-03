package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/bootstrap"
	"github.com/momo-s15/aeroform/internal/config"
	"github.com/spf13/cobra"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Prepare Pro Mode cloud trust and state prerequisites",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadFromFile("config.yaml")
		if err != nil {
			return fmt.Errorf("pro mode requires config.yaml: %w", err)
		}

		out := cmd.OutOrStdout()
		report := bootstrap.Run(cfg)
		fmt.Fprint(out, bootstrap.Render(report))
		fmt.Fprintf(out, "state backend: %s\n", cfg.State.Backend)
		return nil
	},
}
