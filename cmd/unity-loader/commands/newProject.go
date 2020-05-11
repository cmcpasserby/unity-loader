package commands

import (
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
)

func newProject(args ...string) error {
	var projectName string

	if len(args) == 0 {
		prompt := &survey.Input{ // TODO add validation here to ensure a valid project name
			Message: "enter project name",
		}

		if err := survey.AskOne(prompt, &projectName, nil); err != nil {
			return err
		}
	} else {
		projectName = os.Args[0]
	}


	installs, err := unity.GetInstalls()
	if err != nil {
		return err
	}

	options := make([]string, 0, len(installs))
	for _, install := range installs {
		options = append(options, install.Version.String())
	}

	prompt := &survey.Select{
		Message: "Select unity version to create project with",
		Options: options,
		PageSize: 10,
	}

	var version string
	if err := survey.AskOne(prompt, &version, nil); err != nil {
		return err
	}

	appInstall, err := unity.GetInstallFromVersion(version)
	if err != nil {
		return err
	}

	if err := appInstall.NewProject(projectName); err != nil {
		return err
	}

	return nil
}
