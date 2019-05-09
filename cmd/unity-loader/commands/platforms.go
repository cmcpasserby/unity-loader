package commands

import (
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/AlecAivazis/survey.v1"
)

func platforms(args ...string) error {
	if err := repairPaths(true); err != nil {
		return err
	}

	cache, err := settings.ReadCache()
	if err != nil {
		return err
	}

	if cache.NeedsUpdate() {
		if err := cache.Update(); err != nil {
			return err
		}
	}

	var version string

	if len(args) == 0 {
		installs, err := unity.GetInstalls()
		if err != nil {
			return err
		}

		options := make([]string, 0, len(installs))
		for _, install := range installs {
			options = append(options, install.Version.String())
		}

		prompt := &survey.Select{
			Message: "select unity version to install modules on",
			Options: options,
			PageSize: 10,
		}

		if err := survey.AskOne(prompt, &version, nil); err != nil {
			return err
		}
	} else {
		// get version from args
	}



	return nil
}
