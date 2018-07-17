package version

import (
    "path/filepath"
    "os"
    "fmt"
    "github.com/cmcpasserby/unity-loader/internal"
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

    version, err := unityUtils.GetUnityVersion(versionFile)
    if err != nil {
        log.Fatal(err)
    }

    app, err := unityUtils.GetExecutable(version)

    fmt.Printf("version: %s, installed: %t\n", version, app != "")
}
