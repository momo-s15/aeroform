package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold a config.yaml from the example",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("config.yaml"); err == nil {
			fmt.Fprintln(cmd.OutOrStdout(), "config.yaml already exists")
			return nil
		}
		data, err := os.ReadFile("config.yaml.example")
		if err != nil {
			return err
		}
		if err := os.WriteFile("config.yaml", data, 0o644); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "config.yaml created from config.yaml.example")
		return nil
	},
}
