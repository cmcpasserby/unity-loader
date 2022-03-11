package unity

import (
	"howett.net/plist"
	"os"
	"os/exec"
	"path/filepath"
)

type appInfoDict struct {
	CFBundleExecutable string `plist:"CFBundleExecutable"`
	CFBundleVersion    string `plist:"CFBundleVersion"`
}

func getFromInstallPathInternal(path string) (InstallInfo, error) {
	plistPath := filepath.Join(path, "Contents", "Info.plist")
	file, err := os.Open(plistPath)
	if err != nil {
		return InstallInfo{}, err
	}
	defer closeFile(file)

	var appInfo appInfoDict
	if err = plist.NewDecoder(file).Decode(&appInfo); err != nil {
		return InstallInfo{}, err
	}

	ver, err := VersionFromString(appInfo.CFBundleVersion)
	if err != nil {
		return InstallInfo{}, err
	}

	installData := InstallInfo{Path: path, Version: ver}
	return installData, nil
}

func unityGlob(searchPath string) ([]string, error) {
	return filepath.Glob(filepath.Join(searchPath, "**/Unity.app"))
}

func command(path string, args ...string) *exec.Cmd {
	newArgs := append([]string{path, "-W", "-n", "--args"}, args...)
	return exec.Command("open", newArgs...)
}

func binFromApp(path string) (string, error) {
	plistPath := filepath.Join(path, "Contents", "Info.plist")
	file, err := os.Open(plistPath)
	if err != nil {
		return "", err
	}
	defer closeFile(file)

	var appInfo appInfoDict
	if err = plist.NewDecoder(file).Decode(&appInfo); err != nil {
		return "", err
	}

	return filepath.Join(path, "Contents", "MacOS", appInfo.CFBundleExecutable), nil
}
