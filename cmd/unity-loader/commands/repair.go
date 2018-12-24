package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
)

func repair(args ...string) error {
	installs, err := unity.GetInstalls()
	if err != nil {
		return err
	}

	fmt.Println("repairing unity install paths")
	for _, install := range installs {
		err := unity.RepairInstallPath(install)
		if err != nil {
			return err
		}
	}
	return nil
}
