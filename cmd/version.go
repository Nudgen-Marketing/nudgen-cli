package cmd

import (
	"fmt"

	"github.com/nudgen/nudgen-cli/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Nudgen CLI",
	Long:  `All software has versions. This is Nudgen CLI's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.FullVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
