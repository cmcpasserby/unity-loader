package unity

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type InstallInfo struct {
	Path    string
	Version VersionData
}

func (info *InstallInfo) Run(project string) error {
	return info.RunWithTarget(project, "")
}

func (info *InstallInfo) RunWithTarget(project, target string) error {
	absProject, _ := filepath.Abs(project)

	var app *exec.Cmd
	if target == "" {
		app = command(info.Path, "-projectPath", absProject)
	} else {
		app = exec.Command(info.Path, "-projectPath", absProject, "-buildTarget", target)
	}
	return app.Run()
}

func (info *InstallInfo) String() string {
	return fmt.Sprintf("Version: %s Path: \"%s\"", info.Version.String(), info.Path)
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

func GetInstalls(searchPaths ...string) ([]InstallInfo, error) {
	installPaths := make([]string, 0)
	for _, path := range searchPaths {
		globed, err := unityGlob(path)
		if err != nil {
			return nil, err
		}
		installPaths = append(installPaths, globed...)
	}

	installs := make([]InstallInfo, 0, len(installPaths))
	for _, path := range installPaths {
		installData, err := GetInstallFromPath(path)
		if err != nil {
			return nil, err
		}
		installs = append(installs, installData)
	}

	sort.Slice(installs, func(i, j int) bool {
		return installs[i].Version.Compare(installs[j].Version) > 0
	})

	return installs, nil
}

func GetInstallFromVersion(version string, searchPaths ...string) (InstallInfo, error) {
	installs, err := GetInstalls(searchPaths...)
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
