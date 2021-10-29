package main

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/spf13/cobra"
)

func createListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all installed unity versions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := GetConfig()
			if err != nil {
				return err
			}

			installs, err := unity.GetInstalls(config.SearchPaths...)
			if err != nil {
				return err
			}

			for _, data := range installs {
				fmt.Println(data.String())
			}

			return nil
		},
	}
	return cmd
}
