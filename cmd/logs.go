package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/simplestate"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show local launch history until cloud log streaming is wired",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := simplestate.Load()
		if err != nil {
			return err
		}
		out := cmd.OutOrStdout()
		if len(state.Projects) == 0 {
			fmt.Fprintln(out, "No launch history found yet.")
			return nil
		}

		fmt.Fprintln(out, "Launch history")
		for _, project := range state.Projects {
			fmt.Fprintf(out, "- %s | template=%s | created=%s\n", project.Name, project.Template, project.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		return nil
	},
}
