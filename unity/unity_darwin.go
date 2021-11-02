package unity

import (
	"howett.net/plist"
	"os"
	"os/exec"
	"path/filepath"
)

type appInfoDict struct {
	CFBundleVersion string `plist:"CFBundleVersion"`
}

func GetInstallFromPath(path string) (InstallInfo, error) {
	plistPath := filepath.Join(path, "Contents/info.plist")
	file, err := os.Open(plistPath)
	if err != nil {
		return InstallInfo{}, err
	}
	defer closeFile(file)

	var appInfo appInfoDict
	if err = plist.NewDecoder(file).Decode(&appInfo); err != nil {
		return InstallInfo{}, err
	}

	installData := InstallInfo{Path: path, Version: VersionFromString(appInfo.CFBundleVersion)}
	return installData, nil
}

func unityGlob(searchPath string) ([]string, error) {
	return filepath.Glob(filepath.Join(searchPath, "**/Unity.app"))
}

func command(path string, args ...string) *exec.Cmd {
	newArgs := append([]string{"-a", path, "-n", "--args"}, args...)
	return exec.Command("open", newArgs...)
}
