package unity

import (
	"os/exec"
)

func GetInstallFromPath(path string) (InstallInfo, error) {
	cmd := exec.Command(path, "-batchmode", "-version")
	result, err := cmd.Output()
	if err != nil {
		return InstallInfo{}, err
	}
	return InstallInfo{Path: path, Version: VersionFromString(string(result))}, nil
}

func unityGlob(searchPath string) ([]string, error) {
	return filepath.Glob(fmt.Sprintf("%s/**/Editor/Unity.exe"))
}

func command(path string, args ...string) *exec.Cmd {
	newArgs := append([]string{"", path}, args...)
	return exec.Command("start", newArgs...)
}
