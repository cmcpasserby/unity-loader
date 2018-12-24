package commands

import (
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"gopkg.in/AlecAivazis/survey.v1"
)

func install(args ...string) error {
	// TODO check cache timestamp and maybe update

	cache, err := settings.ReadCache()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		versionStrs := make([]string, 0, len(cache.Releases.Official))
		for _, ver := range cache.Releases.Official {
			versionStrs = append(versionStrs, ver.Version)
		}

		prompt := &survey.Select{
			Message:  "select version to install:",
			Options:  versionStrs,
			PageSize: 10,
		}

		var result string
		if err := survey.AskOne(prompt, &result, nil); err != nil {
			return err
		}

		if err := installVersion(result, cache); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func installVersion(version string, cache *settings.Cache) error {
	return nil
}
