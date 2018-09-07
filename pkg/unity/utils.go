package unity

import (
    "bufio"
    "errors"
    "fmt"
    "howett.net/plist"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

type InstallInfo struct {
    Version string
    Path string
}

func (info *InstallInfo) Run(project string) error {
    absProject, _ := filepath.Abs(project)
    app := exec.Command("open", "-a", info.Path, "--args", "-projectPath", absProject)
    return app.Run()
}

type appInfoDict struct {
    CFBundleVersion string `plist:"CFBundleVersion"`
}

func GetVersionFromProject(versionFile string) (string, error) {
    file, _ := os.Open(versionFile)
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)

    for scanner.Scan() {
        text := scanner.Text()
        if strings.HasPrefix(text, "m_EditorVersion:") {
            return strings.TrimSpace(strings.Split(text, ":")[1]), nil
        }
    }
    return "", errors.New("invalid ProjectVersion.txt")
}

func GetInstalls() []InstallInfo {
    unityPaths, _ := filepath.Glob("/Applications/**/Unity.app")

    installs := make([]InstallInfo, 0, len(unityPaths))
    for _, path := range unityPaths {
        installData := GetInstallFromPath(path)
        installs = append(installs, installData)
    }
    return installs
}

func GetInstallFromPath(path string) InstallInfo {
    plistPath := filepath.Join(path, "Contents/info.plist")
    file, _ := os.Open(plistPath)

    var appInfo appInfoDict

    decoder := plist.NewDecoder(file)
    decoder.Decode(&appInfo)

    installData := InstallInfo{Version: appInfo.CFBundleVersion, Path: path}
    return installData
}

func GetInstallFromVersion(version string) (InstallInfo, error) {
    Installs := GetInstalls()

    for _, install := range Installs {
        if version == install.Version {
            return install, nil
        }
    }
    return InstallInfo{}, VersionNotFoundError{version}
}

func RepairInstallPath(install InstallInfo) error {
    oldPath := filepath.Dir(install.Path)
    newName := fmt.Sprintf("Unity %s", install.Version)
    newPath := filepath.Join("/Applications/", newName)

    if oldPath == newPath {
        return nil
    }

    fmt.Printf("moving %q to %q\n", oldPath, newPath)
    err := os.Rename(oldPath, newPath)
    if err != nil {return err}

    return nil
}
