package unity

import (
	"bufio"
	"errors"
	"fmt"
	"howett.net/plist"
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
	absProject, _ := filepath.Abs(project)
	app := exec.Command("open", "-a", info.Path, "--args", "-projectPath", absProject)
	return app.Run()
}

type appInfoDict struct {
	CFBundleVersion string `plist:"CFBundleVersion"`
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

func GetInstalls() ([]*InstallInfo, error) {
	unityPaths, _ := filepath.Glob("/Applications/**/Unity.app")

	installs := make([]*InstallInfo, 0, len(unityPaths))
	for _, path := range unityPaths {
		installData, err := GetInstallFromPath(path)
		if err != nil {
			return nil, err
		}
		installs = append(installs, installData)
	}

	sort.Slice(installs, func(i, j int) bool {
		return !VersionLess(installs[i].Version, installs[j].Version)
	})

	return installs, nil
}

func GetInstallFromPath(path string) (*InstallInfo, error) {
	plistPath := filepath.Join(path, "Contents/info.plist")
	file, err := os.Open(plistPath)
	if err != nil {
		return nil, err
	}
	defer closeFile(file)

	var appInfo appInfoDict

	if err := plist.NewDecoder(file).Decode(&appInfo); err != nil {
		return nil, err
	}

	installData := InstallInfo{Version: VersionDataFromString(appInfo.CFBundleVersion), Path: path}
	return &installData, nil
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

func RepairInstallPath(install *InstallInfo) error {
	oldPath := filepath.Dir(install.Path)
	newName := fmt.Sprintf("Unity %s", install.Version.String())
	newPath := filepath.Join("/Applications/", newName)

	if oldPath == newPath {
		return nil
	}

	fmt.Printf("moving %q to %q\n", oldPath, newPath)
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	return nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
