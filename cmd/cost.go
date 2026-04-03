package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/simplestate"
	"github.com/spf13/cobra"
)

var costCmd = &cobra.Command{
	Use:   "cost",
	Short: "Show estimated monthly spend for tracked Simple Mode projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := simplestate.Load()
		if err != nil {
			return err
		}
		out := cmd.OutOrStdout()
		if len(state.Projects) == 0 {
			fmt.Fprintln(out, "No tracked projects yet. Run `aeroform launch` first.")
			return nil
		}

		total := 0.0
		fmt.Fprintln(out, "Estimated monthly cost by project")
		for _, project := range state.Projects {
			total += project.MonthlyEstimate
			fmt.Fprintf(out, "- %s: $%.2f/mo\n", project.Name, project.MonthlyEstimate)
		}
		fmt.Fprintf(out, "Total: $%.2f/mo\n", total)
		return nil
	},
}
