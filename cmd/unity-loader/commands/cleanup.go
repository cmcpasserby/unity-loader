package commands

import (
	"errors"
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/sudoer"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"path/filepath"
)

func cleanup(args ...string) error {
	// TODO also clean any temp files from the packages folder
	// TODO display and calculate freed space

	config, err := settings.ParseDotFile()
	if err != nil {
		return err
	}
	path := config.ProjectDirectory

	if path == "" {
		return errors.New("projects path is not defined in config.toml")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	projects, err := unity.GetProjectsInPath(path)
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		return nil
	}

	installs, err := unity.GetInstalls()
	if err != nil {
		return err
	}

	toRemove := make([]*unity.InstallInfo, 0, len(installs))

	for _, install := range installs {
		isUsed := false

		for _, proj := range projects {
			projVersion, err := unity.GetVersionFromProject(proj)
			if err != nil {
				continue
			}

			if projVersion == install.Version.String() {
				isUsed = true
			}
		}

		if !isUsed {
			toRemove = append(toRemove, install)
		}
	}

	if len(toRemove) == 0 {
		return nil
	}

	promptTitles := make([]string, 0, len(toRemove))
	for _, item := range toRemove {
		promptTitles = append(promptTitles, item.Version.String())
	}

	prompt := &survey.MultiSelect{
		Message:  "select installs to remove",
		Options:  promptTitles,
		Default:  promptTitles,
		PageSize: 10,
	}

	var promptResults []string
	if err := survey.AskOne(prompt, &promptResults, nil); err != nil {
		return err
	}

	sudo := new(sudoer.Sudoer)
	if err := sudo.AskPass(); err != nil {
		return err
	}

	for _, installVersion := range promptTitles {
		fmt.Printf("Uninstalling Unity Version %q\n", installVersion)

		installInfo, err := unity.GetInstallFromVersion(installVersion)
		if err != nil {
			return err
		}

		if err := sudo.RunAsRoot("rm", "-rf", filepath.Dir(installInfo.Path)); err != nil {
			fmt.Printf("error uninstalling %q, error: %q", installInfo.Path, err)
		}
	}

	return nil
}
