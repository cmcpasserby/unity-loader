package main

import (
    "fmt"
    "os"
    "log"
)

func main() {
    if len(os.Args) == 1 {
        printHelp()
        return
    }

    if len(os.Args) == 2 {
        if _, err := os.Stat(os.Args[1]); err == nil {
            commands["run"].Action(os.Args[1])
            return
        }
    }

    if val, ok := commands[os.Args[1]]; ok {
        err := val.Action(os.Args[2:]...)
        if err != nil{
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

    for _, key := range commandOrder {
        fmt.Printf("  %-12s%s\n", commands[key].Name, commands[key].HelpText)
    }
}
