package cmd

import (
	"fmt"

	"github.com/momo-s15/aeroform/internal/security"
	"github.com/momo-s15/aeroform/internal/ui"
	"github.com/spf13/cobra"
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Manage custom security policies",
	Long: `Manage custom Checkov policies that are included in security scans.

Policies can be Python (.py) or YAML (.yaml/.yml) files following the
Checkov custom policy format. Added policies are stored in .aeroform/policies/
and automatically included in every 'aeroform scan' and Pro Mode security gate.`,
}

var policyAddCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add a custom Checkov policy file",
	Long: `Copy a Checkov custom policy into the project's policy directory.

Supported formats:
  - Python (.py)  — Checkov Python-based custom checks
  - YAML (.yaml)  — Checkov YAML-based custom checks

Example:
  aeroform policy add checks/my_s3_policy.py
  aeroform policy add rules/naming_convention.yaml`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		dst, err := security.AddPolicy(args[0])
		if err != nil {
			return err
		}
		ui.Successln(out, "✓ Policy added: "+dst)
		fmt.Fprintln(out, "  It will be included in the next 'aeroform scan' run.")
		return nil
	},
}

var policyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered custom policies",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		policies, err := security.ListPolicies()
		if err != nil {
			return err
		}
		if len(policies) == 0 {
			fmt.Fprintln(out, "No custom policies registered. Run 'aeroform policy add <file>' to add one.")
			return nil
		}
		ui.Boldln(out, "Custom policies")
		for _, p := range policies {
			fmt.Fprintf(out, "  %s\n", p)
		}
		fmt.Fprintf(out, "\nStored in: %s\n", security.PolicyDir)
		return nil
	},
}

var policyRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a custom policy by filename",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		if err := security.RemovePolicy(args[0]); err != nil {
			return err
		}
		ui.Successln(out, "✓ Policy removed: "+args[0])
		return nil
	},
}

func init() {
	policyCmd.AddCommand(policyAddCmd)
	policyCmd.AddCommand(policyListCmd)
	policyCmd.AddCommand(policyRemoveCmd)
}
