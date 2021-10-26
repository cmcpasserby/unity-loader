package unity

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type InstallInfo struct {
	Path    string
	Version VersionData
}

func (info *InstallInfo) Run(project string) error {
	return info.RunWithTarget(project, "")
}

func (info *InstallInfo) String() string {
	return fmt.Sprintf("Version: %q Path: %q", info.Version.String(), info.Path)
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

func GetInstallFromVersion(version string) (InstallInfo, error) {
	installs, err := GetInstalls()
	if err != nil {
		return InstallInfo{}, err
	}

	for _, install := range installs {
		if version == install.Version.String() {
			return install, nil
		}
	}

	return InstallInfo{}, versionNotFoundError{version}
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
