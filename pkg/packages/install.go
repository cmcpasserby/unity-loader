package packages

import (
    "os"
    "path"
    "fmt"
)

var tempDir string

func Install(version string) error {
    // if os.Getuid() != 0 {
    //     return errors.New("admin is required to install packages, try running with sudo")
    // }

    versionData, err := GetVersionData(version)
    if err != nil {return err}

    packages, err := getPackages(versionData)
    if err != nil {return err}
    defer cleanUp()

    pkg := packages[0]

    err = pkg.DownloadPkg()
    if err != nil {return err}

    isValid, err := pkg.ValidatePkg()
    if err != nil {return err}
    fmt.Printf("is valid: %v", isValid)

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
        fp := path.Join(tempDir, f.Name())
        os.Remove(fp)
    }

    os.Remove(tempDir)
}
