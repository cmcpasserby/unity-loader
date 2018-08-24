package main

import (
    "github.com/cmcpasserby/unity-loader/pkg/unity"
    "path/filepath"
    "os"
    "fmt"
    "log"
    "errors"
)

type command struct {
    Name string
    HelpText string
    Action func(...string) error
}

var commandOrder = [...]string{"run", "version", "list", "install", "repair"}

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
                return err
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
                fmt.Printf("Version: %s Path: %q\n", data.Version, data.Path)
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

    "repair": {
        "repair",
        "fix paths to unity installs",
        func(args ...string) error {
            fmt.Println("repairing unity install paths")
            for _, install := range unity.GetInstalls() {
                oldPath := filepath.Dir(install.Path)
                newName := fmt.Sprintf("Unity %s", install.Version)
                newPath := filepath.Join("/Applications/", newName)

                if oldPath == newPath {continue}
                fmt.Printf("moveing %q to %q\n", oldPath, newPath)
                err := os.Rename(oldPath, newPath)
                if err != nil {return err}
            }
            return nil
        },
    },
}
