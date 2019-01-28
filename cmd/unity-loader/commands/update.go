package commands

import (
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
)

func update(args ...string) error {
	compareVersion := unity.VersionData{Major: 2017, Minor: 1, Update: 5, VerType: "f", Patch: 1}

	err := parsing.GetArchiveVersions(func (data unity.VersionData) bool {
		return unity.VersionLess(data, compareVersion)
	})
	if err != nil {
		return err
	}

	// data, err := parsing.GetHubVersions()
	// if err != nil {
	// 	return err
	// }
	//
	// if err := settings.WriteCache(data); err != nil {
	// 	return err
	// }

	return nil
}
