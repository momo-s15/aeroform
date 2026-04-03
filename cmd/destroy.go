package cmd

import (
	"errors"
	"fmt"

	"github.com/momo-s15/aeroform/internal/simplestate"
	"github.com/momo-s15/aeroform/internal/terraform"
	"github.com/spf13/cobra"
)

var destroyConfirm bool

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy deployed Simple Mode infrastructure and clear tracked projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !destroyConfirm {
			return errors.New("destructive command: re-run with --confirm to destroy infrastructure and clear tracked projects")
		}

		out := cmd.OutOrStdout()
		st, err := simplestate.Load()
		if err != nil {
			return err
		}

		if len(st.Projects) == 0 {
			fmt.Fprintln(out, "No tracked projects to destroy.")
			return nil
		}

		var destroyErrors []string
		for _, project := range st.Projects {
			if project.WorkDir == "" {
				fmt.Fprintf(out, "Skipping %s — no working directory recorded (created before deployment wiring)\n", project.Name)
				continue
			}

			fmt.Fprintf(out, "-> Destroying %s (%s/%s) in %s\n", project.Name, project.Provider, project.Template, project.WorkDir)
			if err := terraform.Destroy(project.WorkDir, true); err != nil {
				destroyErrors = append(destroyErrors, fmt.Sprintf("%s: %v", project.Name, err))
				fmt.Fprintf(out, "   Failed to destroy %s: %v\n", project.Name, err)
				continue
			}
			fmt.Fprintf(out, "   Destroyed %s.\n", project.Name)
		}

		if err := simplestate.Clear(); err != nil {
			return err
		}
		fmt.Fprintln(out, "Tracked projects cleared.")

		if len(destroyErrors) > 0 {
			return fmt.Errorf("some projects failed to destroy: %v", destroyErrors)
		}

		return nil
	},
}

func init() {
	destroyCmd.Flags().BoolVar(&destroyConfirm, "confirm", false, "confirm infrastructure destruction")
}
