package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"path/filepath"
)

func run(args ...string) error {
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
			fmt.Printf("Unity %s not installed\n", version)
			installUnity := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Do you want to install Unity %s?", version),
			}
			if err := survey.AskOne(prompt, &installUnity, nil); err != nil {
				return err
			}
			if installUnity {
				// TODO install unity with version, wait then refresh appInstall
			}
		}
	}

	fmt.Printf("Opening project %q in version: %s\n", expandedPath, version)
	if err := appInstall.Run(path); err != nil {
		return fmt.Errorf("could not execute unity from %q", appInstall.Path)
	}

	return nil
}
