package commands

import (
	"errors"
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/sudoer"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/AlecAivazis/survey.v1"
	"path"
)

func uninstall(args ...string) error {
	versions := make([]string, 0)

	if len(args) == 0 {
		installs, err := unity.GetInstalls()
		if err != nil {
			return err
		}

		options := make([]string, 0, len(installs))
		for _, install := range installs {
			options = append(options, install.Version.String())
		}

		prompt := &survey.MultiSelect{
			Message:  "Select versions to uninstall",
			Options:  options,
			PageSize: len(options),
		}
		if err := survey.AskOne(prompt, &versions, nil); err != nil {
			return err
		}
	} else {
		versions = args
	}

	validInstalls := make([]*unity.InstallInfo, 0, len(versions))
	for _, ver := range versions {
		if install, err := unity.GetInstallFromVersion(ver); err == nil {
			validInstalls = append(validInstalls, install)
		}
	}

	if len(validInstalls) == 0 {
		return errors.New("nothing to uninstall")
	}

	sudo := new(sudoer.Sudoer)
	if err := sudo.AskPass(); err != nil {
		return err
	}

	for _, install := range validInstalls {
		fmt.Printf("Uninstalling Unity Version %q\n", install.Version.String())
		if err := sudo.RunAsRoot("rm", "-rf", path.Dir(install.Path)); err != nil {
			return fmt.Errorf("error uninstalling %q, Error: %q", install.Path, err)
		}
	}

	// TODO Display Freed up space

	return nil
}
