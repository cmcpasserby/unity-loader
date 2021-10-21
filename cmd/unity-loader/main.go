package main

import (
	"github.com/spf13/cobra"
	"os"
)

var ( // Global Flags
	gFlagNoConfirm bool
)

func main() {
	cmd := &cobra.Command{
		Use:     "unity-loader",
		Version: "3.0.0",
		Short:   "Tool for loading unity projects with their respective unity versions and installing the proper version if required",
	}

	cmd.PersistentFlags().BoolVarP(&gFlagNoConfirm, "no-Confirm", "y", false, "removes confirmation prompts")

	cmd.AddCommand(
		createRunCmd(),
		createVersionCmd(),
		createListCmd(),
	)

	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}
