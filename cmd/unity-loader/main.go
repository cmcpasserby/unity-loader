package main

import (
    "github.com/cmcpasserby/unity-loader/pkg/unity"
    "fmt"
    )

func main() {
    versions, _ := unity.ParseVersions(unity.UnityDownloads)
    fmt.Println(versions["2017.3.1f1"].GetAndroidSupportUrl())

    // if len(os.Args) == 1 {
    //     printHelp(&unity.Commands)
    //     return
    // }
    //
    // if len(os.Args) == 2 {
    //     if _, err := os.Stat(os.Args[1]); err == nil {
    //         unity.Commands["run"].Action(os.Args[1])
    //         return
    //     }
    // }
    //
    // if val, ok := unity.Commands[os.Args[1]]; ok {
    //     val.Action(os.Args[2:]...)
    // } else {
    //     fmt.Printf("%q is not a valid command\n", os.Args[1])
    //     fmt.Println()
    //     printHelp(&unity.Commands)
    // }
}

func printHelp(commands *map[string]unity.Command) {
    fmt.Println("usage: unity-loader <commands>, [project_path]")
    fmt.Println()
    fmt.Println("commands are:")

    for _, cmd := range *commands {
        fmt.Printf("  %-12s%s\n", cmd.Name, cmd.HelpText)
    }
}
