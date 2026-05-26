package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time or defaults to "dev".
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "mpt",
	Short: "my-power-toys - A tiny cross-platform project launcher for developers.",
	Long:  "my-power-toys (mpt) is a local development productivity tool.\nIt helps you manage projects and open them with your preferred tools.",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of mpt",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "my-power-toys "+Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
