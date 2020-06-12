package cmd

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/spf13/cobra"
)

func createListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all installed unity versions",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			installs, err := unity.GetInstalls()
			if err != nil {
				return err
			}

			for _, data := range installs {
				fmt.Printf("Version: %q Path: %q\n", data.Version.String(), data.Path)
			}

			return nil
		},
	}
	return cmd
}
