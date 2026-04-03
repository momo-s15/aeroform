package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/simplestate"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Print the latest tracked project URL",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := simplestate.Load()
		if err != nil {
			return err
		}
		if len(state.Projects) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No tracked project found.")
			return nil
		}
		project := state.Projects[len(state.Projects)-1]
		fmt.Fprintf(cmd.OutOrStdout(), "Latest tracked project: %s\n", project.Name)
		if project.CustomDomain != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Custom domain: https://%s\n", project.CustomDomain)
		}
		return nil
	},
}
