package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/cmcpasserby/scli"
	"github.com/cmcpasserby/verinfo"
	"os"
)

func main() {
	fs := flag.NewFlagSet("root", flag.ExitOnError)
	versionFlag := fs.Bool("v", false, "prints unity-loader's version")

	cmd := &scli.Command{
		Usage:         "unity-loader <subcommand>",
		ShortHelp:     "Tool for loading unity projects with their respective unity versions",
		LongHelp:      "Tool for loading unity projects with their respective unity versions",
		FlagSet:       fs,
		ArgsValidator: scli.NoArgs,
		Subcommands: []*scli.Command{
			createRunCmd(),
			createVersionCmd(),
			createListCmd(),
			// createSearchCmd(),
			// createInstallCmd(),
		},
		Exec: func(ctx context.Context, args []string) error {
			if !*versionFlag {
				return flag.ErrHelp
			}

			info, err := verinfo.Get()
			if err != nil {
				return err
			}

			fmt.Printf("unity-loader version %s (%s)\n", info.Version, info.Revision)
			return nil
		},
	}

	if err := cmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}

		fmt.Println(err)
		os.Exit(1)
	}
}
