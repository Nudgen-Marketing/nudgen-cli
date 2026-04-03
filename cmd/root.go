package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var UseJSON bool

var rootCmd = &cobra.Command{
	Use:   "nudgen",
	Short: "Nudgen CLI for Agent-Driven Lead Ops",
	Long:  `Nudgen is an interactive CLI for managing campaigns, contacts, and telemetry for Nudgen Agent Ops.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&UseJSON, "json", false, "Output results in JSON format for machine readability")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
