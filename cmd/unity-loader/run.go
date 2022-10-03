package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cmcpasserby/scli"
	"github.com/cmcpasserby/unity-loader/unity"
	"os"
	"path/filepath"
)

func createRunCmd() *scli.Command {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	forceVersionFlag := fs.String("force", "", "force project to be opened with a specific Unity version")
	buildTargetFlag := fs.String("buildTarget", "", "opens project with a specific build target set")

	return &scli.Command{
		Usage:         "run [projectDirectory]",
		ShortHelp:     "Launches unity and opens the selected project",
		LongHelp:      "Launches unity and opens the selected project",
		FlagSet:       fs,
		ArgsValidator: scli.MaxArgs(1),
		Exec: func(ctx context.Context, args []string) error {
			cfg, err := getConfig()
			if err != nil {
				return err
			}

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

			expandedPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			var ver unity.VersionData

			if *forceVersionFlag != "" {
				ver, err = unity.VersionFromString(*forceVersionFlag)
				if err != nil {
					return err
				}
			} else {
				ver, err = unity.GetVersionFromProject(path)
				if err != nil {
					return err
				}
			}

			appInstall, err := unity.GetInstallFromVersion(ver, cfg.SearchPaths...)
			if err != nil {
				return err
			}

			return runInstalledVersion(appInstall, expandedPath, *buildTargetFlag)
		},
	}
}

func runInstalledVersion(installInfo unity.InstallInfo, projectPath, target string) error {
	fmt.Printf("Opening project \"%s\" in version %s\n", projectPath, installInfo.Version.String())
	return installInfo.RunWithTarget(projectPath, target)
}
