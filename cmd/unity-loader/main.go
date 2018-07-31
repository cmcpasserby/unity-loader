package main

import (
    "github.com/cmcpasserby/unity-loader/pkg/unity"
    "fmt"
    "os"
)

func main() {
    commands := unity.GetCommands()

    if len(os.Args) == 1 {
        printHelp(&commands)
        return
    }

    if len(os.Args) == 2 {
        if _, err := os.Stat(os.Args[1]); err == nil {
            commands["run"].Action(os.Args[1])
            return
        }
    }

    if val, ok := commands[os.Args[1]]; ok {
        val.Action(os.Args[2:]...)
    } else {
        fmt.Printf("%q is not a valid command\n", os.Args[1])
        fmt.Println()
        printHelp(&commands)
    }
}

func printHelp(cmds *map[string]*unity.Command)  {
    fmt.Println("usage: unity-loader <command>, [project_path]")
    fmt.Println()
    fmt.Println("commands are:")

    for _, cmd := range *cmds {
        fmt.Printf("  %-12s%s\n", cmd.Name, cmd.HelpText)
    }
}
