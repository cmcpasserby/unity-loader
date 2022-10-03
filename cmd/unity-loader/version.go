package main

import (
	"context"
	"fmt"
	"github.com/cmcpasserby/scli"
	"github.com/cmcpasserby/unity-loader/unity"
	"os"
)

func createVersionCmd() *scli.Command {
	return &scli.Command{
		Usage:         "version [projectDirectory]",
		ShortHelp:     "Check what version of unity a project is using",
		LongHelp:      "Check what version of unity a project is using",
		ArgsValidator: scli.MaxArgs(1),
		Exec: func(ctx context.Context, args []string) error {
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

			projectVersion, err := unity.GetVersionFromProject(path)
			if err != nil {
				return err
			}

			config, err := getConfig()
			if err != nil {
				return err
			}

			_, err = unity.GetInstallFromVersion(projectVersion, config.SearchPaths...)
			fmt.Printf("projectVersion: %q installed: %t\n", projectVersion, err == nil)

			return nil
		},
	}
}
