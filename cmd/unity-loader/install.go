package main

import (
    "errors"
    "fmt"
    "github.com/cmcpasserby/unity-loader/pkg/packages"
    "github.com/cmcpasserby/unity-loader/pkg/sudoer"
    "github.com/cmcpasserby/unity-loader/pkg/unity"
    "gopkg.in/AlecAivazis/survey.v1"
    "io/ioutil"
    "os"
)

const baseInstallPath = "/Applications/Unity/Unity.app"

func Install(version string) error {
    versionData, err := packages.GetVersionData(version)
    if err != nil {return err}

    sudo := new(sudoer.Sudoer)
    if !sudo.AskPass() {
        return errors.New("invalid admin password\n")
    }

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

    // if a unity install exists in the base path move it before a new install starts
    if _, err := os.Stat(baseInstallPath); err == nil {
        installInfo := unity.GetInstallFromPath(baseInstallPath)
        unity.RepairInstallPath(installInfo)
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

        err = pkg.Install(sudo)
        if err != nil {return err}
    }

    // after a install do no leave it in the base install path, but move to versioned folder
    if _, err := os.Stat(baseInstallPath); err == nil {
        installInfo := unity.GetInstallFromPath(baseInstallPath)
        unity.RepairInstallPath(installInfo)
    }

    return nil
}

func cleanUp(tempDir string) {
    err := os.RemoveAll(tempDir)
    if err != nil {
        panic(err)
    }
}
