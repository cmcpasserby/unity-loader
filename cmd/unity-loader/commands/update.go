package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"sort"
)

func update(args ...string) error {
	fmt.Println("Fetching Unity Hub Versions...")
	hubVersions, err := parsing.GetHubVersions()
	if err != nil {
		return err
	}

	sort.Sort(hubVersions.Official)
	hubOldest := unity.VersionDataFromString(hubVersions.Official[0].Version)

	pkgs, err := parsing.GetArchiveVersions(func (data unity.VersionData) bool {
		return unity.VersionLess(data, hubOldest)
	})
	if err != nil {
		return err
	}

	hubVersions.Official = append(hubVersions.Official, pkgs...)

	if err := settings.WriteCache(hubVersions); err != nil {
		return err
	}

	return nil
}
