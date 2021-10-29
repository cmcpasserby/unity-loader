package unity

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func (info *InstallInfo) RunWithTarget(project, target string) error {
	panic("not implemented yet")
}

func (info *InstallInfo) NewProject(projectName string) error {
	panic("not implemented yet")
}

func unityGlob(searchPath string) ([]string, error) {
	return filepath.Glob(fmt.Sprintf("%s/**/Editor/Unity.exe"))
}

func GetInstallFromPath(path string) (InstallInfo, error) {
	cmd := exec.Command(path, "-batchmode", "-version")
	result, err := cmd.Output()
	if err != nil {
		return InstallInfo{}, err
	}
	return InstallInfo{Path: path, Version: VersionFromString(string(result))}, nil
}
