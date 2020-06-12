package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
)

func list(args ...string) error {
	installs, err := unity.GetInstalls()
	if err != nil {
		return err
	}

	for _, data := range installs {
		fmt.Printf("Version: %q Path: %q\n", data.Version.String(), data.Path)
	}
	return nil
}
