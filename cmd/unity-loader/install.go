package main

import (
	"context"
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func createInstallCmd() *ffcli.Command {
	return &ffcli.Command{
		Name:       "install",
		ShortUsage: "unity-loader install [version] [versionHash]",
		ShortHelp:  "Installs a Unity version.",
		LongHelp:   "Installs a Unity version.",
		Exec: func(ctx context.Context, args []string) error {
			argc := len(args)
			if argc == 0 || argc > 2 {
				return fmt.Errorf("install expected 1 or 2 arguments got %d", argc)
			}

			cfg, err := getConfig()
			if err != nil {
				return err
			}

			var ver unity.VersionData

			if argc == 2 {
				ver, err = unity.VersionFromString(args[0])
				if err != nil {
					return err
				}
				ver.RevisionHash = args[1]
			} else {
				archives, err := unity.GetAllVersions()
				if err != nil {
					return err
				}

				for _, item := range archives {
					if item.String() == args[0] {
						ver = item
						break
					}
				}

				if (ver == unity.VersionData{}) {
					return fmt.Errorf("version %s is not in the archive", args[0])
				}
			}

			fmt.Printf("installing: %s (%s)\n", ver.String(), ver.RevisionHash)

			installInfo, err := unity.InstallFromArchive(ver, cfg.UnityHubPath, cfg.DefaultModules, cfg.SearchPaths)
			if err != nil {
				return err
			}

			fmt.Printf("installed %s to %s", installInfo.Version, installInfo.Path)
			return nil
		},
	}
}
