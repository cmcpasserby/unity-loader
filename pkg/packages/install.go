package packages

func Install(version string) error {
    // if os.Getuid() != 0 {
    //     return errors.New("admin is required to install packages, try running with sudo")
    // }

    versionData, err := GetVersionData(version)
    if err != nil {return err}

    packages, err := getPackages(versionData)
    if err != nil {return err}

    pkg := packages[0]

    err = pkg.DownloadPkg()
    if err != nil {return err}

    _, err = pkg.ValidatePkg()
    if err != nil {return err}

    err = pkg.CleanupPkg()
    if err != nil {return err}

    return nil
}
