package main

import (
    "path/filepath"
    "os"
    "fmt"
    "github.com/cmcpasserby/unity-loader/pkg/loader"
    "log"
)

func main() {
    printVersion(os.Args[2])
}

func printVersion(path string) {
    versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
    if _, err := os.Stat(versionFile); os.IsNotExist(err) {
        fmt.Printf("%q is not a valid unity project\n", path)
    }

    version, err := loader.GetUnityVersion(versionFile)
    if err != nil {
        log.Fatal(err)
    }

    app, err := loader.GetExecutable(version)

    fmt.Printf("version: %s, installed: %t\n", version, app != "")
}
