package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the Aeroform version",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "aeroform %s\n", Version)
		if Version == "dev" {
			if info, ok := debug.ReadBuildInfo(); ok {
				for _, s := range info.Settings {
					if s.Key == "vcs.revision" && len(s.Value) >= 7 {
						fmt.Fprintf(out, "  commit: %s\n", s.Value[:7])
					}
					if s.Key == "vcs.time" {
						fmt.Fprintf(out, "  built:  %s\n", s.Value)
					}
				}
			}
		} else {
			fmt.Fprintf(out, "  commit: %s\n", Commit)
			fmt.Fprintf(out, "  built:  %s\n", Date)
		}
		return nil
	},
}
