package run

import (
    "os"
    "fmt"
    "path/filepath"
    "github.com/cmcpasserby/unity-loader/internal"
    "log"
    "os/exec"
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

    version, err := unityUtils.GetUnityVersion(versionFile)
    if err != nil {
        log.Fatal(err)
    }

    app, err := unityUtils.GetExecutable(version)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Opening Unity Version: %s", version)
    exec.Command("open", app)
}

