package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
)

func update(args ...string) error {
	data, err := parsing.GetHubVersions()
	if err != nil {
		return err
	}
	fmt.Println(data)
	return nil
}
