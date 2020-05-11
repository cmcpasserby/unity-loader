package commands

import (
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
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

	var versionStr string

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

		if err := survey.AskOne(prompt, &versionStr, nil); err != nil {
			return err
		}
	} else {
		// TODO handle passed in version number
	}

	version := cache.Releases.First(func(details parsing.CacheVersion) bool {
		return details.String() == versionStr
	})

	if err := installVersion(*version, true); err != nil {
		return err
	}

	return nil
}
