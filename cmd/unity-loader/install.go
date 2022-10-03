package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cmcpasserby/scli"
	"github.com/cmcpasserby/unity-loader/unity"
	"strings"
)

func createInstallCmd() *scli.Command {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	modulesFlag := fs.String("modules", "", "module overrides")

	return &scli.Command{
		Usage:         "unity-loader install [version] [versionHash]",
		ShortHelp:     "Installs a Unity version.",
		LongHelp:      "Installs a Unity version.",
		FlagSet:       fs,
		ArgsValidator: scli.RangeArgs(1, 2),
		Exec: func(ctx context.Context, args []string) error {
			cfg, err := getConfig()
			if err != nil {
				return err
			}

			var ver unity.VersionData

			if len(args) == 2 {
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

			modules := cfg.DefaultModules
			if *modulesFlag != "" {
				modules = strings.Split(*modulesFlag, ",")
			}

			fmt.Printf("installing: %s (%s)\n", ver.String(), ver.RevisionHash)
			installInfo, err := unity.InstallFromArchive(ver, cfg.UnityHubPath, modules, cfg.SearchPaths)
			if err != nil {
				return err
			}

			fmt.Printf("installed %s to %s", installInfo.Version, installInfo.Path)
			return nil
		},
	}
}
