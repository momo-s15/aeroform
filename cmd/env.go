package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment workspaces",
}

var envAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Create a new environment workspace stub",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintf(cmd.OutOrStdout(), "Environment %q would be cloned from the current config in a future phase.\n", args[0])
		return nil
	},
}

func init() {
	envCmd.AddCommand(envAddCmd)
}
