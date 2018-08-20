package main

import (
    "fmt"
    "github.com/cmcpasserby/unity-loader/pkg/packages"
)

func Install(version string) error {
    // if os.Getuid() != 0 {
    //     return errors.New("admin is required to install pkgs, try running with sudo")
    // }

    versionData, err := packages.GetVersionData(version)
    if err != nil {return err}

    pkgs, err := packages.GetPackages(versionData)
    if err != nil {return err}

    pkgs = packages.Filter(pkgs, func(pkg *packages.Package) bool {return !pkg.Data.Hidden})

    for _, pkg := range pkgs {
        fmt.Println(pkg.Data.Description)
    }

    // err = pkg.DownloadPkg()
    // if err != nil {return err}
    //
    // _, err = pkg.ValidatePkg()
    // if err != nil {return err}
    //
    // err = pkg.CleanupPkg()
    // if err != nil {return err}

    return nil
}
