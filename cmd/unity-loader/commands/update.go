package commands

import (
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
)

func update(args ...string) error {
	versions, err := parsing.GetVersions()
	if err != nil {
		return err
	}

	if err := settings.WriteCache(versions); err != nil {
		return err
	}
	return nil
}
