package main

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/cmd/unity-loader/commands"
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
  unity-loader <command>, [project_path]

commands are:`)

	for _, key := range commands.CommandOrder {
		fmt.Printf("  %-12s%s\n", commands.Commands[key].Name, commands.Commands[key].HelpText)
	}
}
