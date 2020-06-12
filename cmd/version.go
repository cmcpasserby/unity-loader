package cmd

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/spf13/cobra"
	"os"
)

func createVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version [projectDirectory]",
		Short: "check what version of unity a project is using",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var path string
			if len(args) > 0 {
				path = args[0]
			} else {
				wd, err := os.Getwd()
				if err != nil {
					return err
				}
				path = wd
			}

			version, err := unity.GetVersionFromProject(path)
			if err != nil {
				return err
			}

			_, err = unity.GetInstallFromVersion(version)
			fmt.Printf("version: %q installed: %t\n", version, err == nil)

			return nil
		},
	}
	return cmd
}
