package unity

import (
    // "fmt"
    "os"
    "path"
    // "errors"
    "errors"
)

var tempDir string

func Install(version string) error {
    if os.Getuid() != 0 {
        return errors.New("admin is required to install packages, try running with sudo")
    }

    versionData, err := GetVersionData(version)
    if err != nil {return err}

    packages, err := getPackages(versionData)
    if err != nil {return err}
    defer cleanUp()

    return nil
}

func cleanUp() {
    if tempDir == "" {
        return
    }

    dirRead, _ := os.Open(tempDir)
    dirFiles, _ := dirRead.Readdir(0)

    for i := range dirFiles {
        f := dirFiles[i]
        path := path.Join(tempDir, f.Name())
        os.Remove(path)
    }

    os.Remove(tempDir)
}
