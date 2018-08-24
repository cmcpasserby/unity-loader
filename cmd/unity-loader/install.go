package main

import (
    "github.com/cmcpasserby/unity-loader/pkg/packages"
    "gopkg.in/AlecAivazis/survey.v1"
    "os"
    "errors"
    "fmt"
    "io/ioutil"
    "path/filepath"
)

func Install(version string) error {
    if os.Getuid() != 0 {
        return errors.New("admin is required to install pkgs, try running with sudo")
    }

    versionData, err := packages.GetVersionData(version)
    if err != nil {return err}

    pkgs, err := packages.GetPackages(versionData)
    if err != nil {return err}

    pkgs = packages.Filter(pkgs, func(pkg *packages.Package) bool {return !pkg.Data.Hidden})

    titles := make([]string, 0, len(pkgs))
    defaults := make([]string, 0)

    for _, pkg := range pkgs {
        titles = append(titles, pkg.Data.Title)
        if pkg.Data.Install {
            defaults = append(defaults, pkg.Data.Title)
        }
    }

    prompt := &survey.MultiSelect{
        Message: "Select Platforms to install:",
        Options: titles,
        Default: defaults,
        PageSize: len(titles),
    }

    var resultStrings []string
    survey.AskOne(prompt, &resultStrings, nil)

    resultPackages := make([]packages.Package, 0, len(resultStrings))
    for _, pkg := range pkgs {
        for _, resultStr := range resultStrings {
            if pkg.Data.Title == resultStr {
                resultPackages = append(resultPackages, *pkg)
            }
        }
    }

    tempDir, err := ioutil.TempDir("", "unitypackage_")
    if err != nil {return err}
    defer cleanUp(tempDir)

    for _, pkg := range resultPackages {
        err = pkg.Download(tempDir)
        if err != nil {return err}

        isValid, err := pkg.Validate()
        if err != nil {return err}
        if !isValid {
            return fmt.Errorf("%q was not a valid package, try installing again\n", pkg.Data.Title)
        }

        err = pkg.Install()
        if err != nil {return err}
    }

    appPath := "/Applications/Unity/"
    if _, err := os.Stat(appPath); os.IsExist(err) {
        newName := fmt.Sprintf("Unity %s", version)
        newPath := filepath.Join("/Applications/", newName)
        err = os.Rename(appPath, newPath)
        if err != nil {return err}
        fmt.Printf("Installed Unity %s to %q\n", version, newPath)
    }

    return nil
}

func cleanUp(tempDir string) {
    err := os.RemoveAll(tempDir)
    if err != nil {
        panic(err)
    }
}
