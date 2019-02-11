package commands

import (
	"errors"
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"path/filepath"
)

func run(args ...string) error {
	// TODO add a --force flag to let you user define what version to open in

	var path string

	if len(args) == 0 {
		path, _ = os.Getwd()
	} else {
		path = args[0]
	}

	expandedPath, _ := filepath.Abs(path)

	version, err := unity.GetVersionFromProject(path)
	if err != nil {
		return err
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

	if err := runInstallVersion(appInstall, expandedPath); err != nil {
		return err
	}

	return nil
}

func runInstallVersion(installInfo *unity.InstallInfo, projectPath string) error {
	fmt.Printf("Opening project %q in version: %s\n", projectPath, installInfo.Version.String())
	return installInfo.Run(projectPath)
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

	if cacheVersion := cache.Releases.First(func (ver parsing.CacheVersion) bool {return ver.String() == version}); cacheVersion != nil {
		installUnity := false

		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Do you want to install Unity %s", version),
		}

		if err := survey.AskOne(prompt, &installUnity, nil); err != nil {
			return nil, err
		}

		if installUnity {
			if err := installVersion(*cacheVersion); err != nil {
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
