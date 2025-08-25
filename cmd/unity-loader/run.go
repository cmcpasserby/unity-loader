package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cmcpasserby/scli"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/joho/godotenv"
)

func createRunCmd() *scli.Command {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	forceVersionFlag := fs.String("force", "", "force project to be opened with a specific Unity version")
	buildTargetFlag := fs.String("buildTarget", "", "opens project with a specific build target set")
	overloadEnvFlag := fs.Bool("overloadEnv", false, "should pre-existing env vars be overwritten but dotenv file")
	noEnvFlag := fs.Bool("noEnv", false, "prevents loading or overloading of dotenv file and applying it to the environment")

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

			if !*noEnvFlag {
				if err = loadEnv(*overloadEnvFlag); err != nil {
					return err
				}
			}

			return runInstalledVersion(appInstall, expandedPath, *buildTargetFlag)
		},
	}
}

func loadEnv(overload bool) error {
	f := godotenv.Load
	if overload {
		f = godotenv.Overload
	}

	if err := f(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func runInstalledVersion(installInfo unity.InstallInfo, projectPath, target string) error {
	fmt.Printf("Opening project \"%s\" in version %s\n", projectPath, installInfo.Version.String())
	return installInfo.RunWithTarget(projectPath, target)
}
