package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/peterbourgon/ff/v3/ffcli"
	"os"
	"path/filepath"
)

func createRunCmd() *ffcli.Command {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	forceVersionFlag := fs.String("force", "", "force project to be opened with a specific Unity version")
	buildTargetFlag := fs.String("buildTarget", "", "opens project with a specific build target set")

	return &ffcli.Command{
		Name:       "run",
		ShortUsage: "unity-loader run [projectDirectory]",
		ShortHelp:  "Launches unity and opens the selected project",
		LongHelp:   "Launches unity and opens the selected project",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("accepts at most %d args(s), received %d", 1, len(args))
			}

			config, err := getConfig()
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

			var version unity.VersionData

			if *forceVersionFlag != "" {
				version = unity.VersionFromString(*forceVersionFlag)
			} else {
				version, err = unity.GetVersionFromProject(path)
				if err != nil {
					return err
				}
			}

			appInstall, err := unity.GetInstallFromVersion(version, config.SearchPaths...)
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
