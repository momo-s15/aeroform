package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Manage custom security policies",
}

var policyAddCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add a custom Checkov policy stub",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintf(cmd.OutOrStdout(), "Policy file %q would be registered in a future phase.\n", args[0])
		return nil
	},
}

func init() {
	policyCmd.AddCommand(policyAddCmd)
}
