package main

import (
    "os"
    "fmt"
    "path/filepath"
    "log"
    "os/exec"
    "github.com/cmcpasserby/unity-loader/pkg/loader"
)

func main() {
    if len(os.Args[1:]) > 1 {
        runUnity(os.Args[2])
    } else {
        runUnity(os.Args[1])
    }
}

func runUnity(path string) {
    versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
    if _, err := os.Stat(versionFile); os.IsNotExist(err) {
        fmt.Printf("%q is not a valid unity project\n", path)
    }

    version, err := loader.GetUnityVersion(versionFile)
    if err != nil {
        log.Fatal(err)
    }

    app, err := loader.GetExecutable(version)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Opening Unity Version: %s", version)
    exec.Command("open", app)
}

