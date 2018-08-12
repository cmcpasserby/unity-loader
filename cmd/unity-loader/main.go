package main

import (
    "fmt"
    "os"
    "log"
)

func main() {
    if len(os.Args) == 1 {
        printHelp(&Commands)
        return
    }

    if len(os.Args) == 2 {
        if _, err := os.Stat(os.Args[1]); err == nil {
            Commands["run"].Action(os.Args[1])
            return
        }
    }

    if val, ok := Commands[os.Args[1]]; ok {
        err := val.Action(os.Args[2:]...)
        if err != nil{
            log.Fatal(err)
        }
    } else {
        fmt.Printf("%q is not a valid command\n", os.Args[1])
        fmt.Println()
        printHelp(&Commands)
    }
}

func printHelp(commands *map[string]Command) {
    fmt.Println(
`Tool for loading unity projects with their respective unity versions and installing hte proper version if required

usage:
  unity-loader <command>, [project_path]

commands are:`)

    for _, cmd := range *commands {
        fmt.Printf("  %-12s%s\n", cmd.Name, cmd.HelpText)
    }
}
