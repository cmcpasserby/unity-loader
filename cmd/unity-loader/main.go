package main

import (
	"github.com/spf13/cobra"
	"os"
)

var version = "3.0.0" // left as a var so it can be updated via ldflags

func main() {
	cmd := &cobra.Command{
		Use:     "unity-loader",
		Version: version,
		Short:   "Tool for loading unity projects with their respective unity versions and installing the proper version if required",
	}

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
