// Package main is the entry point to the Plugin Marketplace server.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Marketplace is a repository of Mattermost plugins.",
	// SilenceErrors allows us to explicitly log the error returned from rootCmd below.
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.WithError(err).Error("command failed")
		os.Exit(1)
	}
}
