package main

import (
    "errors"
    "fmt"
    "github.com/cmcpasserby/unity-loader/pkg/sudoer"
    "github.com/cmcpasserby/unity-loader/pkg/unity"
    "gopkg.in/AlecAivazis/survey.v1"
    "log"
    "os"
    "path"
    "path/filepath"
    "time"
)

type command struct {
    Name string
    HelpText string
    Action func(...string) error
}

var commandOrder = [...]string{"run", "version", "list", "install", "uninstall", "repair"}

var commands = map[string]command {

    "run": {
        "run",
        "run the passed in project with an auto detected version of unity",
        func(args ...string) error {
            var path string

            if len(args) == 0 {
                path, _ = os.Getwd()
            } else {
                path = args[0]
            }

            versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
            if _, err := os.Stat(versionFile); os.IsNotExist(err) {
                fmt.Printf("%q is not a valid unity project\n", path)
            }

            version, err := unity.GetVersionFromProject(versionFile)
            if err != nil {
                return err
            }

            appInstall, err := unity.GetInstallFromVersion(version)
            if err != nil {
                if _, ok := err.(unity.VersionNotFoundError); ok {
                    fmt.Printf("Unity %s not installed\n", version)
                    installUnity := false
                    prompt := &survey.Confirm{
                        Message: fmt.Sprintf("Do you want to install Unity %s?", version),
                    }
                    survey.AskOne(prompt, &installUnity, nil)
                    if installUnity {
                        Install(version)
                        time.Sleep(time.Second)
                        appInstall, _ = unity.GetInstallFromVersion(version)
                    }
                } else {
                    return err
                }
            }

            fmt.Printf("Opening project %q in version: %s\n", path, version)
            err = appInstall.Run(path)
            if err != nil {
                return fmt.Errorf("could not execute unity from %q", appInstall.Path)
            }
            return nil
        },
    },

    "version": {
        "version",
        "check what version of unity a project is using",
        func(args ...string) error {
            var path string

            if len(args) == 0 {
                path, _ = os.Getwd()
            } else {
                path = args[0]
            }

            versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
            if _, err := os.Stat(versionFile); os.IsNotExist(err) {
                return fmt.Errorf("%q is not a valid unity project\n", path)
            }

            version, err := unity.GetVersionFromProject(versionFile)
            if err != nil {
                return err
            }

            _, err = unity.GetInstallFromVersion(version)
            fmt.Printf("version: %s, installed: %t\n", version, err == nil)

            return nil
        },
    },

    "list": {
        "list",
        "list all installed unity versions",
        func(args ...string) error {
            for _, data := range unity.GetInstalls() {
                fmt.Printf("Version: %q Path: %q\n", data.Version.String(), data.Path)
            }
            return nil
        },
    },

    "install": {
        "install",
        "installed the specified version of unity",
        func(args ...string) error {
            if len(args) == 0 {
                return errors.New("no version specified")
            }

            version := args[0]
            err := Install(version)
            if err != nil {
                log.Fatal("ERROR: ", err)
            }
            return nil
        },
    },

    "uninstall": {
        "uninstall",
        "uninstall one or multiple versions of Unity",
        func(args ...string) error {
            versions := make([]string, 0)

            if len(args) == 0 {
                installs := unity.GetInstalls()

                options := make([]string, 0, len(installs))
                for _, install := range installs {
                    options = append(options, install.Version.String())
                }

                prompt := &survey.MultiSelect{
                    Message: "Select versions to uninstall",
                    Options: options,
                    PageSize:len(options),
                }

                survey.AskOne(prompt, &versions, nil)
            } else {
                versions = args
            }

            validInstalls := make([]unity.InstallInfo, 0, len(versions))
            for _, ver := range versions {
                install, err := unity.GetInstallFromVersion(ver)
                if err != nil {continue}
                validInstalls = append(validInstalls, install)
            }

            if len(validInstalls) == 0 {
                return errors.New("nothing to uninstall")
            }

            sudo := new(sudoer.Sudoer)
            if !sudo.AskPass() {
                return errors.New("invalid admin password\n")
            }

            for _, install := range validInstalls {
                fmt.Printf("Uninstalling Unity Version %q\n", install.Version.String())
                sudo.RunAsRoot("rm", "-rf", path.Dir(install.Path))
            }

            return nil
        },
    },

    "repair": {
        "repair",
        "fix paths to unity installs",
        func(args ...string) error {
            fmt.Println("repairing unity install paths")
            for _, install := range unity.GetInstalls() {
                err := unity.RepairInstallPath(install)
                if err != nil {return err}
            }
            return nil
        },
    },
}
