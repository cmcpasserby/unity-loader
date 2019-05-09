package commands

import (
	"errors"
	"flag"
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"path/filepath"
)

func run(args ...string) error {
	flagSet := flag.NewFlagSet("run", flag.ExitOnError)
	forceFlag := flagSet.Bool("force", false, "force a certain version to be used for running")
	targetFlag := flagSet.String("buildTarget", "", "Allows the selection of an active build target before loading a project")

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	var path string

	if len(flagSet.Args()) == 0 {
		path, _ = os.Getwd()
	} else {
		path = flagSet.Args()[0]
	}

	expandedPath, _ := filepath.Abs(path)

	var version string
	var err error

	if forceFlag != nil && *forceFlag {
		installs, err := unity.GetInstalls()
		if err != nil {
			return err
		}

		options := make([]string, 0, len(installs))
		for _, install := range installs {
			options = append(options, install.Version.String())
		}

		prompt := &survey.Select{
			Message:  "Select Unity version to run project with",
			Options:  options,
			PageSize: 10,
		}

		if err := survey.AskOne(prompt, &version, nil); err != nil {
			return err
		}
	} else {
		version, err = unity.GetVersionFromProject(path)
		if err != nil {
			return err
		}
	}

	appInstall, err := unity.GetInstallFromVersion(version)
	if err != nil {
		if _, ok := err.(unity.VersionNotFoundError); ok {
			if appInstall, err = installAndGetInfo(version); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err := runInstallVersion(appInstall, expandedPath, *targetFlag); err != nil {
		return err
	}

	return nil
}

func runInstallVersion(installInfo *unity.InstallInfo, projectPath, target string) error {
	fmt.Printf("Opening project %q in version: %s\n", projectPath, installInfo.Version.String())
	return installInfo.RunWithTarget(projectPath, target)
}

func installAndGetInfo(version string) (*unity.InstallInfo, error) {
	cache, err := settings.ReadCache()
	if err != nil {
		return nil, err
	}

	if cache.NeedsUpdate() {
		if err := cache.Update(); err != nil {
			return nil, err
		}
	}

	if cacheVersion := cache.Releases.First(func(ver parsing.CacheVersion) bool { return ver.String() == version }); cacheVersion != nil {
		installUnity := false

		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Do you want to install Unity %s", version),
		}

		if err := survey.AskOne(prompt, &installUnity, nil); err != nil {
			return nil, err
		}

		if installUnity {
			if err := installVersion(*cacheVersion, false); err != nil {
				return nil, err
			}

			appInstall, err := unity.GetInstallFromVersion(version)
			if err != nil {
				return nil, err
			}
			return appInstall, nil
		} else {
			return nil, errors.New(fmt.Sprintf("unity %s not found or installed", version))
		}
	}
	return nil, errors.New(fmt.Sprintf("unity %s not found for download", version))
}
