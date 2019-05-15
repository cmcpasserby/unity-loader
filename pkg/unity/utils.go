package unity

import (
	"bufio"
	"errors"
	"fmt"
	"howett.net/plist"
	"io/ioutil"
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

func (info *InstallInfo) RunWithTarget(project, target string) error {
	absProject, _ := filepath.Abs(project)

	var app *exec.Cmd
	if target == "" {
		app = exec.Command("open", "-a", info.Path, "--args", "-projectPath", absProject)
	} else {
		app = exec.Command("open", "-a", info.Path, "--args", "-projectPath", absProject, "-buildTarget", target)
	}

	return app.Run()
}

func (info *InstallInfo) Run(project string) error {
	return info.RunWithTarget(project, "")
}

func (info *InstallInfo) GetPlatforms() ([]string, error) {
	return nil, errors.New("not implemented")
}

type appInfoDict struct {
	CFBundleVersion string `plist:"CFBundleVersion"`
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

func InstallToUnityDir(install *InstallInfo) error {
	oldPath := filepath.Dir(install.Path)
	newPath := "/Applications/Unity"

	if oldPath == newPath {
		return nil
	}

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
