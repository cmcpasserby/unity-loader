package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/sudoer"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"path"
)

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

	projects, err := unity.GetProjectsInPath(paths[0])
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		return nil
	}

	installs, err := unity.GetInstalls()
	if err != nil {
		return err
	}

	toRemove := make([]*unity.InstallInfo, 0, len(installs))

	for _, install := range installs {
		isUsed := false

		for _, proj := range projects {
			projVersion, err := unity.GetVersionFromProject(proj)
			if err != nil {
				continue
			}

			if projVersion == install.Version.String() {
				isUsed = true
			}
		}

		if !isUsed {
			toRemove = append(toRemove, install)
		}
	}

	if len(toRemove) == 0 {
		return nil
	}

	sudo := new(sudoer.Sudoer)
	if err := sudo.AskPass(); err != nil {
		return err
	}

	for _, install := range toRemove {
		fmt.Printf("Uninstalling Unity Version %q\n", install.Version.String())
		if err := sudo.RunAsRoot("rm", "-rf", path.Dir(install.Path)); err != nil {
			return fmt.Errorf("error uninstalling %q, Error: %q", install.Path, err)
		}
	}

	return nil
}
