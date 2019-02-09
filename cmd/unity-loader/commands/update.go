package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
)

func update(args ...string) error {
	fmt.Println("Fetching Unity Hub Versions...")
	hubVersions, err := parsing.GetHubVersions()
	if err != nil {
		return err
	}

	pkgs, err := parsing.GetArchiveVersions(func (data unity.VersionData) bool {
		for _, hubVer := range hubVersions.Official {
			if hubVer.Version == data.String() {
				return false
			}
		}
		return true
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
