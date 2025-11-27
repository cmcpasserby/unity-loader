package main

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/cmcpasserby/scli"
	"github.com/cmcpasserby/unity-loader/serve"
)

func createServeCmd() *scli.Command {
	const (
		defaultPort = 8080
	)

	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	portFlag := fs.Int("port", defaultPort, "port to serve on")

	return &scli.Command{
		Usage:         "serve [targetDirectory]",
		ShortHelp:     "Serves a WebGL build",
		LongHelp:      "Serves a WebGL build",
		FlagSet:       fs,
		ArgsValidator: scli.MaxArgs(1),
		Exec: func(ctx context.Context, args []string) error {
			port := *portFlag
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			absPath, err := filepath.Abs(dir)
			if err != nil {
				return err
			}

			return serve.Serve(absPath, port)
		},
	}
}
