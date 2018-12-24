package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"os"
)

func version(args ...string) error {
	var path string

	if len(args) == 0 {
		path, _ = os.Getwd()
	} else {
		path = args[0]
	}

	version, err := unity.GetVersionFromProject(path)
	if err != nil {
		return err
	}

	_, err = unity.GetInstallFromVersion(version)
	fmt.Printf("version: %q installed: %t", version, err == nil)

	return nil
}
