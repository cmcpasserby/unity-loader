package main

import (
	"context"
	"fmt"
	"github.com/cmcpasserby/scli"
	"github.com/cmcpasserby/unity-loader/unity"
)

func createSearchCmd() *scli.Command {
	return &scli.Command{
		Usage:         "unity-loader search [partialVersion]",
		ShortHelp:     "Searches for a unity version on the archive site",
		LongHelp:      "Search for a unity version on the archive site, partial numbers can be listed and all matches will be returned",
		ArgsValidator: scli.ExactArgs(1),
		Exec: func(ctx context.Context, args []string) error {
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
