package unity

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	unityInstallHubPath = baseUnityPath + "/Hub"
	unityHubEditorPath  = unityInstallHubPath + "/Editor"
)

type InstallInfo struct {
	Path    string
	Version VersionData
}

func (info *InstallInfo) Run(project string) error {
	return info.RunWithTarget(project, "")
}

func GetProjectsInPath(projectPath string) ([]string, error) {
	projects := make([]string, 0)

	folders, err := ioutil.ReadDir(projectPath)
	if err != nil {
		return nil, err
	}

	for _, f := range folders {
		if !f.IsDir() {
			continue
		}

		projectVersionPath := filepath.Join(projectPath, f.Name(), "ProjectSettings", "ProjectVersion.txt")
		if _, err := os.Stat(projectVersionPath); !os.IsNotExist(err) {
			projects = append(projects, filepath.Join(projectPath, f.Name()))
		}
	}

	return projects, nil
}

func GetVersionFromProject(path string) (string, error) {
	versionFile := filepath.Join(path, "ProjectSettings", "ProjectVersion.txt")
	if _, err := os.Stat(versionFile); os.IsNotExist(err) {
		return "", fmt.Errorf("%q is not a valid unity project\n", path)
	}

	file, err := os.Open(versionFile)
	if err != nil {
		return "", err
	}
	defer closeFile(file)

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

func GetInstallFromVersion(version string) (*InstallInfo, error) {
	installs, err := GetInstalls()
	if err != nil {
		return nil, err
	}

	for _, install := range installs {
		if version == install.Version.String() {
			return install, nil
		}
	}

	return nil, VersionNotFoundError{version}
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
