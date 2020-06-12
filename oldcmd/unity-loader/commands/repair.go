package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
)

func repair(args ...string) error {
	return repairPaths(false)
}

func repairPaths(silent bool) error {
	installs, err := unity.GetInstalls()
	if err != nil {
		return err
	}

	if !silent {
		fmt.Println("Repairing unity install paths")
	}

	for _, install := range installs {
		err := unity.RepairInstallPath(install)
		if err != nil {
			return err
		}
	}
	return nil
}
