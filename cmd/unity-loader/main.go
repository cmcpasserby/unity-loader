package main

import (
    "fmt"
    "os"
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
        val.Action(os.Args[2:]...)
    } else {
        fmt.Printf("%q is not a valid command\n", os.Args[1])
        fmt.Println()
        printHelp(&Commands)
    }
}

func printHelp(commands *map[string]Command) {
    fmt.Println("usage: unity-loader <commands>, [project_path]")
    fmt.Println()
    fmt.Println("commands are:")

    for _, cmd := range *commands {
        fmt.Printf("  %-12s%s\n", cmd.Name, cmd.HelpText)
    }
}
