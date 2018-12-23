package commands

import (
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
)

func update(args ...string) error {
	data, err := parsing.GetHubVersions()
	if err != nil {
		return err
	}

	if err := settings.WriteCache(data); err != nil {
		return err
	}

	return nil
}
