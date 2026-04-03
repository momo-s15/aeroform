package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/security"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run Pro Mode security scanners against Terraform files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := config.LoadFromFile("config.yaml"); err != nil {
			return fmt.Errorf("pro mode requires config.yaml: %w", err)
		}

		report := security.RunProSecurityScan(".")
		out := cmd.OutOrStdout()
		fmt.Fprintln(out, "Aeroform Pro Mode scan")
		for _, line := range report.SummaryLines() {
			fmt.Fprintln(out, line)
		}
		if report.Blocking {
			return fmt.Errorf("security gate blocked")
		}
		return nil
	},
}
