package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cmcpasserby/scli"
)

func createBuildProfilesCmd() *scli.Command {
	const profilesPathRelative = "Assets/Settings/Build Profiles"

	return &scli.Command{
		Usage:         "buildProfiles",
		ShortHelp:     "Lists Builds Profiles found in project",
		LongHelp:      "Lists Builds Profiles found in project",
		ArgsValidator: scli.NoArgs(),
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

			profilesPath := filepath.Join(path, profilesPathRelative)
			profiles, err := filepath.Glob(filepath.Join(profilesPath, "*.asset"))
			if err != nil {
				return err
			}

			for _, profile := range profiles {
				profileName := filepath.Base(profile[:len(profile)-6])
				fmt.Println(profileName)
			}

			return nil
		},
	}
}
