package main

import (
	"context"
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func createListCmd() *ffcli.Command {
	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "list",
		ShortHelp:  "List all installed unity versions",
		LongHelp:   "List all installed unity versions",
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
