package main

import (
	"github.com/spf13/cobra"
)

var ( // Global Flags
	gFlagConfig    string
	gFlagNoConfirm bool
)

var rootCmd = &cobra.Command{
	Use:     "unity-loader",
	Version: "3.0.0",
	Short:   "Tool for loading unity projects with their respective unity versions and installing the proper version if required",
}

func execute() error {
	// get config path stuff
	var configPath string
	rootCmd.PersistentFlags().StringVarP(&gFlagConfig, "config", "c", configPath, "sets the active config for unity-loader")
	rootCmd.PersistentFlags().BoolVarP(&gFlagNoConfirm, "no-Confirm", "y", false, "removes confirmation prompts")

	rootCmd.AddCommand(
		createRunCmd(),
		createVersionCmd(),
		createListCmd(),
	)

	return rootCmd.Execute()
}
