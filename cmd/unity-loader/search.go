package main

import (
	"context"
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func createSearchCmd() *ffcli.Command {
	return &ffcli.Command{
		Name:       "search",
		ShortUsage: "unity-loader search [partialVersion]",
		ShortHelp:  "Searches for a unity version on the archive site",
		LongHelp:   "Search for a unity version on the archive site, partial numbers can be listed and all matches will be returned",
		Exec: func(ctx context.Context, args []string) error {
			argc := len(args)
			if argc == 0 || argc > 1 {
				return fmt.Errorf("search expected 1 argumenet got %d", argc)
			}

			results, err := unity.SearchArchive(args[0])
			if err != nil {
				return err
			}

			for _, ver := range results {
				fmt.Printf("%s (%s)\n", ver.String(), ver.RevisionHash)
			}
			return nil
		},
	}
}
