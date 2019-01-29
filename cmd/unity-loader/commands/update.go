package commands

import (
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"sort"
)

func update(args ...string) error {
	hubVersions, err := parsing.GetHubVersions()
	if err != nil {
		return err
	}

	sort.Sort(hubVersions.Official)
	hubOldest := unity.VersionDataFromString(hubVersions.Official[0].Version)

	if err := parsing.GetArchiveVersions(func (data unity.VersionData) bool {
		return unity.VersionLess(data, hubOldest)
	}); err != nil {
		return err
	}

	// if err := settings.WriteCache(data); err != nil {
	// 	return err
	// }

	return nil
}
