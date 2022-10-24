package main

import (
	"context"
	"fmt"
	"github.com/cmcpasserby/scli"
	"github.com/cmcpasserby/unity-loader/unity"
)

func createListCmd() *scli.Command {
	return &scli.Command{
		Usage:         "list",
		ShortHelp:     "List all installed unity versions",
		LongHelp:      "List all installed unity versions",
		ArgsValidator: scli.NoArgs(),
		Exec: func(ctx context.Context, args []string) error {
			config, err := getConfig()
			if err != nil {
				return err
			}

			installs, err := unity.GetInstalls(config.SearchPaths...)
			if err != nil {
				return err
			}

			for _, data := range installs {
				fmt.Println(data.String())
			}

			return nil
		},
	}
}
