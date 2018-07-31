package unity

import (
    "path/filepath"
    "os"
    "fmt"
    "log"
)

type Command struct {
    Name string
    HelpText string
    Action CommandAction
}

type CommandAction func(...string) error

func NewCommand(name string, helpText string, action CommandAction) *Command {
    return &Command{Name:name, HelpText:helpText, Action:action}
}

var Commands = map[string]*Command {

    "run": NewCommand(
        "run",
        "run the passed in project with a auto detected version of unity",
        func(args ...string) error {
            path := args[0]

            versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
            if _, err := os.Stat(versionFile); os.IsNotExist(err) {
                fmt.Printf("%q is not a valid unity project\n", path)
            }

            version, err := GetVersionFromProject(versionFile)
            if err != nil {
                log.Fatal(err)
            }

            appInstall, err := GetInstallFromVersion(version)
            if err != nil {
                log.Fatalf("Unity version %q not found", version)
            }

            fmt.Printf("Opening project %q in version: %s\n", path, version)
            err = appInstall.Run(path)
            if err != nil {
                log.Fatalf("Could not execute unity from %q", appInstall.Path)
            }
            return nil
        }),

    "version": NewCommand(
        "version",
        "check what version of unity a project is using",
        func(args ...string) error {
            path := args[0]

            versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
            if _, err := os.Stat(versionFile); os.IsNotExist(err) {
                fmt.Printf("%q is not a valid unity project\n", path)
            }

            version, err := GetVersionFromProject(versionFile)
            if err != nil {
                log.Fatal(err)
            }

            _, err = GetInstallFromVersion(version)
            fmt.Printf("version: %s, installed: %t\n", version, err == nil)

            return nil
        }),

    "list": NewCommand(
        "list",
        "list all installed unity versions",
        func(args ...string) error {
            for _, data := range GetInstalls() {
                fmt.Printf("Version: %s Path: %q", data.Version, data.Path)
            }
            return nil
        }),
}
