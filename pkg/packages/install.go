package packages

import "fmt"

func Install(version string) error {
    // if os.Getuid() != 0 {
    //     return errors.New("admin is required to install packages, try running with sudo")
    // }

    versionData, err := GetVersionData(version)
    if err != nil {return err}

    packages, err := getPackages(versionData)
    if err != nil {return err}

    packages = filter(packages, func(pkg *Package) bool {return !pkg.Data.Hidden})

    for _, pkg := range packages {
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
