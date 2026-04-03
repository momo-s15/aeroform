package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/simplestate"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Simple Mode projects tracked locally",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := simplestate.Load()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if len(state.Projects) == 0 {
			fmt.Fprintln(out, "No projects found. Run `aeroform launch` first.")
			return nil
		}

		fmt.Fprintln(out, "Simple Mode projects")
		for _, project := range state.Projects {
			fmt.Fprintf(out, "- %s | %s | %s | $%.2f/mo\n", project.Name, project.Provider, project.Template, project.MonthlyEstimate)
		}
		return nil
	},
}
