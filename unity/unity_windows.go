package unity

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetInstallFromPath(path string) (InstallInfo, error) {
	cmd := exec.Command(path, "-batchmode", "-version")
	result, err := cmd.Output()
	if err != nil {
		return InstallInfo{}, err
	}
	verString := strings.TrimSpace(string(result))
	return InstallInfo{Path: path, Version: VersionFromString(verString)}, nil
}

func unityGlob(searchPath string) ([]string, error) {
	return filepath.Glob(fmt.Sprintf("%s/**/Editor/Unity.exe", searchPath))
}

func command(path string, args ...string) *exec.Cmd {
	return exec.Command(path, args...)
}
