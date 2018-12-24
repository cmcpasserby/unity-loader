package commands

import "github.com/cmcpasserby/unity-loader/pkg/settings"

func cleanup(args ...string) error {
	var paths []string

	if len(args) >= 1 {
		paths = args
	} else {
		config, err := settings.ParseDotFile()
		if err != nil {
			return err
		}
		// TODO ensure project directory is defined
		paths = []string{config.ProjectDirectory}
	}

	return nil
}
