package main

import (
	"context"
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/peterbourgon/ff/v3/ffcli"
	"os"
)

func createVersionCmd() *ffcli.Command {
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "version [projectDirectory]",
		ShortHelp:  "Check what version of unity a project is using",
		LongHelp:   "Check what version of unity a project is using",
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

			version, err := unity.GetVersionFromProject(path)
			if err != nil {
				return err
			}

			config, err := getConfig()
			if err != nil {
				return err
			}

			_, err = unity.GetInstallFromVersion(version, config.SearchPaths...)
			fmt.Printf("version: %q installed: %t\n", version, err == nil)

			return nil
		},
	}
}
