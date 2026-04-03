package cmd

import (
	"errors"
	"fmt"

	"github.com/momo-s15/aeroform/internal/simplestate"
	"github.com/spf13/cobra"
)

var destroyConfirm bool

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Remove tracked Simple Mode projects (state only for now)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !destroyConfirm {
			return errors.New("destructive command: re-run with --confirm to clear tracked projects")
		}
		if err := simplestate.Clear(); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Tracked projects cleared.")
		return nil
	},
}

func init() {
	destroyCmd.Flags().BoolVar(&destroyConfirm, "confirm", false, "confirm project state deletion")
}
