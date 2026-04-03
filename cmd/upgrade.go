package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/simplestate"
	"github.com/spf13/cobra"
)

var upgradeConfirm bool
var upgradeOutput string
var upgradeForce bool

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Generate a Pro Mode config from Simple Mode projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := simplestate.Load()
		if err != nil {
			return err
		}
		out := cmd.OutOrStdout()

		if len(state.Projects) == 0 {
			fmt.Fprintln(out, "No projects tracked yet. Build your first project with aeroform launch.")
			return nil
		}

		cfg, err := engine.BuildProConfigFromSimpleState(state)
		if err != nil {
			return err
		}
		rendered, err := engine.RenderConfigYAML(cfg)
		if err != nil {
			return err
		}

		total := 0.0
		for _, project := range state.Projects {
			total += project.MonthlyEstimate
		}

		fmt.Fprintln(out, "Simple to Pro upgrade preview")
		fmt.Fprintf(out, "Tracked projects: %d\n", len(state.Projects))
		fmt.Fprintf(out, "Estimated monthly spend: $%.2f\n", total)
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, "Generated config preview:")
		fmt.Fprintln(out, rendered)

		if !upgradeConfirm {
			fmt.Fprintln(out, "Run with --confirm to write config.yaml and switch to Pro Mode.")
			return nil
		}

		if !upgradeForce {
			if _, err := os.Stat(upgradeOutput); err == nil {
				return errors.New("output file already exists; use --force to overwrite")
			}
		}

		if err := os.WriteFile(upgradeOutput, []byte(rendered), 0o644); err != nil {
			return err
		}

		fmt.Fprintf(out, "Wrote %s\n", upgradeOutput)
		fmt.Fprintln(out, "Pro Mode config generated. You can now run aeroform plan or aeroform generate.")
		return nil
	},
}

func init() {
	upgradeCmd.Flags().BoolVar(&upgradeConfirm, "confirm", false, "write the generated Pro Mode config")
	upgradeCmd.Flags().StringVar(&upgradeOutput, "output", "config.yaml", "output file for generated Pro Mode config")
	upgradeCmd.Flags().BoolVar(&upgradeForce, "force", false, "overwrite output file if it already exists")
}
