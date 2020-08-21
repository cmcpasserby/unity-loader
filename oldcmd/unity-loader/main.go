package unity_loader

import (
	"flag"
	"fmt"
	"github.com/cmcpasserby/unity-loader/oldcmd/unity-loader/commands"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	if len(os.Args) == 2 {
		if _, err := os.Stat(os.Args[1]); err == nil {
			if val, ok := commands.Commands["run"]; ok {
				if err := val.Action(os.Args[1:]...); err != nil {
					log.Fatal(err)
				}
			}
			return
		}
	}

	if val, ok := commands.Commands[os.Args[1]]; ok {
		if err := val.Action(os.Args[2:]...); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("%q is not a valid command\n", os.Args[1])
		fmt.Println()
		printHelp()
	}
}

func printHelp() {
	fmt.Println(
		`Tool for loading unity projects with their respective unity versions and installing the proper version if required

usage:
  unity-loader <command> [flags] [project_path]

commands are:`)

	maxNameLen := 0
	maxDescLen := 0

	for _, key := range commands.CommandOrder {
		cmd := commands.Commands[key]
		if len(cmd.Name) > maxNameLen {
			maxNameLen = len(cmd.Name)
		}

		if len(cmd.HelpText) > maxDescLen {
			maxDescLen = len(cmd.HelpText)
		}
	}

	maxNameLen += 2
	maxDescLen += 2

	for _, key := range commands.CommandOrder {
		cmd := commands.Commands[key]
		fmt.Printf("  %-*s%-*sflags: [", maxNameLen, cmd.Name, maxDescLen, cmd.HelpText)

		hasFlags := false

		if cmd.Flags != nil {
			cmd.Flags.VisitAll(func(flag *flag.Flag) {
				fmt.Printf("--%s, ", flag.Name)
				hasFlags = true
			})
		}

		if hasFlags {
			fmt.Printf("\033[2D]")
		} else {
			fmt.Printf("]")
		}

		fmt.Println()
	}
}
