package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aeroform",
	Short: "Cloud infrastructure for everyone",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(launchCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(bootstrapCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(driftCmd)
	rootCmd.AddCommand(templateCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(policyCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(costCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(domainCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(upgradeCmd)
}
