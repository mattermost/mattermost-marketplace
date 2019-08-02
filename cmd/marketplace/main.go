// Package main is the entry point to the Mattermost Integrations Marketplace server.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Marketplace is a repository of Mattermost plugins.",
	Run: func(cmd *cobra.Command, args []string) {
		serverCmd.RunE(cmd, args)
	},
	// SilenceErrors allows us to explicitly log the error returned from rootCmd below.
	SilenceErrors: true,
}

func init() {
	rootCmd.MarkFlagRequired("database")

	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.WithError(err).Error("command failed")
		os.Exit(1)
	}
}
