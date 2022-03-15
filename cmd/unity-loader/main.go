package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3/ffcli"
	"os"
)

var version = "3.0.1" // left as a var so it can be updated via ldflags when built from action

func main() {
	fs := flag.NewFlagSet("root", flag.ExitOnError)
	versionFlag := fs.Bool("v", false, "prints unity-loader's version")

	cmd := &ffcli.Command{
		Name:       "unity-loader",
		ShortUsage: "unity-loader <subcommand>",
		ShortHelp:  "Tool for loading unity projects with their respective unity versions",
		LongHelp:   "Tool for loading unity projects with their respective unity versions",
		FlagSet:    fs,
		Subcommands: []*ffcli.Command{
			createRunCmd(),
			createVersionCmd(),
			createListCmd(),
			createSearchCmd(),
			createInstallCmd(),
		},
		Exec: func(ctx context.Context, args []string) error {
			if !*versionFlag {
				return flag.ErrHelp
			}

			fmt.Printf("unity-loader version %s\n", version)
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
