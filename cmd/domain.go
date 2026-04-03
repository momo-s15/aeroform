package cmd

import (
	"errors"
	"fmt"

	"github.com/momo-s15/aeroform/internal/simplestate"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage Simple Mode custom domains",
}

var domainAddCmd = &cobra.Command{
	Use:   "add [domain]",
	Short: "Attach a custom domain to the newest tracked project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		state, err := simplestate.Load()
		if err != nil {
			return err
		}
		if len(state.Projects) == 0 {
			return errors.New("no tracked project found; run `aeroform launch` first")
		}

		state.Projects[len(state.Projects)-1].CustomDomain = domain
		if err := simplestate.Save(state); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Domain %s saved for project %s.\n", domain, state.Projects[len(state.Projects)-1].Name)
		return nil
	},
}

func init() {
	domainCmd.AddCommand(domainAddCmd)
}
