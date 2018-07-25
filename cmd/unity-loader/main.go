package main

import (
    "os"
    "fmt"
    "path/filepath"
    "log"
    "github.com/cmcpasserby/unity-loader/pkg/unity"
)

func main() {
    if len(os.Args) == 1 {
        fmt.Println("usage: unity-loader <command>, [project_path]")
        fmt.Println("commands are: ")
        fmt.Println("  run         run the passed in project with a auto detected version of unity")
        fmt.Println("  version     check what version of unity a project is using")
        fmt.Println("  list        list all installed unity versions")
    }

    switch os.Args[1] {
    case "run":
        runUnity(os.Args[2])
    case "version":
        printVersion(os.Args[2])
    case "list":
        listVersions()
    default:
        log.Fatalf("%q is not a valid command.\n", os.Args[1])
    }
}

func runUnity(path string) {
    versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
    if _, err := os.Stat(versionFile); os.IsNotExist(err) {
        fmt.Printf("%q is not a valid unity project\n", path)
    }

    version, err := unity.GetVersionFromProject(versionFile)
    if err != nil {
        log.Fatal(err)
    }

    appInstall, err := unity.GetInstallFromVersion(version)
    if err != nil {
        log.Fatalf("Unity version %q not found", version)
    }

    fmt.Printf("Opening project %q in version: %s\n", path, version)
    err = appInstall.Run(path)
    if err != nil {
        log.Fatalf("Could not execute unity from %q", appInstall.Path)
    }
}

func printVersion(path string) {
    versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
    if _, err := os.Stat(versionFile); os.IsNotExist(err) {
        fmt.Printf("%q is not a valid unity project\n", path)
    }

    version, err := unity.GetVersionFromProject(versionFile)
    if err != nil {
        log.Fatal(err)
    }

    _, err = unity.GetInstallFromVersion(version)
    fmt.Printf("version: %s, installed: %t\n", version, err == nil)
}

func listVersions() {
    for _, data := range unity.GetInstalls() {
        fmt.Printf("Version: %s Path: %q\n", data.Version, data.Path)
    }
}
