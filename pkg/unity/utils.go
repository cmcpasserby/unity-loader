package unity

import (
	"howett.net/plist"
	"os"
	"os/exec"
	"path/filepath"
)

type InstallInfo struct {
	Path    string
	Version VersionData
}

func (info *InstallInfo) Run(project string) error {
	absProject, _ := filepath.Abs(project)
	app := exec.Command("open", "-a", info.Path, "--args", "-projectPath", absProject)
	return app.Run()
}

type appInfoDict struct {
	CFBundleVersion string `plist:"CFBundleVersion"`
}

func GetInstalls() []InstallInfo {
	unityPaths, _ := filepath.Glob("/Applications/**/Unity.app")

	installs := make([]InstallInfo, 0, len(unityPaths))
	for _, path := range unityPaths {
		installData := GetInstallFromPath(path)
		installs = append(installs, installData)
	}

	// TODO sort installs

	return installs
}

func GetInstallFromPath(path string) InstallInfo {
	plistPath := filepath.Join(path, "Contents/info.plist")
	file, _ := os.Open(plistPath)

	var appInfo appInfoDict

	decoder := plist.NewDecoder(file)
	decoder.Decode(&appInfo) // TODO handle error

	installData := InstallInfo{Version: VersionDataFromString(appInfo.CFBundleVersion), Path: path}
	return installData
}

func GetInstallFromVersion(version string) (InstallInfo, error) {
	installs := GetInstalls()

	for _, install := range installs {
		if version == install.Version.String() {
			return install, nil
		}
	}

	return InstallInfo{}, VersionNotFoundError{version}
}
