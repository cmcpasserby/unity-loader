package commands

import (
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/sudoer"
	"gopkg.in/AlecAivazis/survey.v1"
	"sort"
)

func install(args ...string) error {
	// TODO check cache timestamp and maybe update

	cache, err := settings.ReadCache()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		sort.Sort(sort.Reverse(cache.Releases.Official))

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
	sudo := new(sudoer.Sudoer)
	if err := sudo.AskPass(); err != nil {
		return err
	}

	installInfo := cache.Releases.First(func(details parsing.PkgDetails) bool {
		return details.Version == version
	})

	titles := make([]string, 0, len(installInfo.Modules))
	defaults := make([]string, 0, len(installInfo.Modules))

	for _, module := range installInfo.Modules {
		titles = append(titles, module.Name)
		if module.Selected {
			defaults = append(defaults, module.Name)
		}
	}

	prompt := &survey.MultiSelect{
		Message:  "select platforms to install",
		Options:  titles,
		Default:  defaults,
		PageSize: len(titles),
	}

	var results []string
	if err := survey.AskOne(prompt, &results, nil); err != nil {
		return err
	}

	return nil
}
