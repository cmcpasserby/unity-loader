package cmd

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

func Execute() error {
	// get config path stuff
	var configPath string
	rootCmd.PersistentFlags().StringVarP(&gFlagConfig, "config", "c", configPath, "sets teh active config for unity-loader")
	rootCmd.PersistentFlags().BoolVarP(&gFlagNoConfirm, "no-Confirm", "y", false, "removes confirmation prompts")

	rootCmd.AddCommand(
		createRunCmd(),
		createVersionCmd(),
	)

	return rootCmd.Execute()
}
