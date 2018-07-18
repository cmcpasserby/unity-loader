package main

import (
    "os"
    "fmt"
    "path/filepath"
    "log"
    "github.com/cmcpasserby/unity-loader/pkg/unity"
    "os/exec"
)

func main() {
    if len(os.Args[1:]) > 1 {
        switch os.Args[1] {
        case "run":
            runUnity(os.Args[2])
        case "version":
            printVersion(os.Args[2])
        default:
            fmt.Printf("%q is not a valid command.\n", os.Args[1])
        }
    } else {
        runUnity(os.Args[1])
    }
}

func runUnity(path string) {
    versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
    if _, err := os.Stat(versionFile); os.IsNotExist(err) {
        fmt.Printf("%q is not a valid unity project\n", path)
    }

    version, err := unity.GetVersionFromProject(versionFile)
    if err != nil {
        log.Fatal(err)
    }

    appPath, err := unity.GetExecutableFromVersion(version)
    if err != nil {
        log.Fatal(err)
    }

    absProjectPath, _ := filepath.Abs(path)


    fmt.Printf("Opening Unity Version: %s\n", version)

    app := exec.Command("open", "-a", appPath, "--args", "-projectPath", absProjectPath)

    err = app.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func printVersion(path string) {
    versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
    if _, err := os.Stat(versionFile); os.IsNotExist(err) {
        fmt.Printf("%q is not a valid unity project\n", path)
    }

    version, err := unity.GetVersionFromProject(versionFile)
    if err != nil {
        log.Fatal(err)
    }

    app, err := unity.GetExecutableFromVersion(version)

    fmt.Printf("version: %s, installed: %t\n", version, app != "")
}
