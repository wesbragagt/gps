package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wesbragagt/gps/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Display the version, commit, build date, and platform information for gps.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Get().String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
