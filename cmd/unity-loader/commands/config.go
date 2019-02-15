package commands

import (
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"os"
	"os/exec"
	"path"
)

func config(args ...string) error {
	dotPath, err := settings.GetPath()
	if err != nil {
		return err
	}

	dotFilePath := path.Join(dotPath, ".config.toml")

	if _, err := os.Stat(dotFilePath); os.IsNotExist(err) {
		if err := settings.CreateDotFile(dotFilePath); err != nil {
			return err
		}
	}

	cmd := exec.Command("vim", dotFilePath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
