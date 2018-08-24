package unity

import (
    "os"
    "bufio"
    "strings"
    "errors"
    "path/filepath"
    "howett.net/plist"
    "os/exec"
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

    var installs []InstallInfo

    for _, path := range unityPaths {
        plistPath := filepath.Join(path, "Contents/info.plist")
        file, _ := os.Open(plistPath)

        var appInfo appInfoDict

        decoder := plist.NewDecoder(file)
        decoder.Decode(&appInfo)

        installData := InstallInfo{Version: appInfo.CFBundleVersion, Path: path}
        installs = append(installs, installData)
    }
    return installs
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
