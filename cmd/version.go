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
		info, ok := debug.ReadBuildInfo()

		display := Version
		if display == "dev" && ok && info != nil {
			if mv := info.Main.Version; mv != "" && mv != "(devel)" {
				display = mv
			}
		}
		fmt.Fprintf(out, "aeroform %s\n", display)

		if Version != "dev" {
			fmt.Fprintf(out, "  commit: %s\n", Commit)
			fmt.Fprintf(out, "  built:  %s\n", Date)
			return nil
		}

		if !ok || info == nil {
			return nil
		}

		if info.Main.Path != "" {
			fmt.Fprintf(out, "  module: %s\n", info.Main.Path)
		}
		if mv := info.Main.Version; mv != "" {
			fmt.Fprintf(out, "  mod_version: %s\n", mv)
		}
		if info.GoVersion != "" {
			fmt.Fprintf(out, "  go: %s\n", info.GoVersion)
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if len(s.Value) >= 7 {
					short := s.Value
					if len(short) > 12 {
						short = short[:12]
					}
					fmt.Fprintf(out, "  commit: %s\n", short)
				}
			case "vcs.time":
				fmt.Fprintf(out, "  built:  %s\n", s.Value)
			case "vcs.modified":
				if s.Value == "true" {
					fmt.Fprintf(out, "  vcs_dirty: true\n")
				}
			}
		}
		return nil
	},
}
