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

const (
	baseUnityPath       = "/Applications/Unity"
	unityInstallHubPath = baseUnityPath + "/Hub"
	unityHubEditorPath  = unityInstallHubPath + "/Editor"
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

func (info *InstallInfo) NewProject(projectName string) error {
	projectPath, _ := filepath.Abs(projectName)
	app := exec.Command("open", "-a", info.Path, "--args", "-createProject", projectPath)
	fmt.Printf("Creating project at %s\n", projectPath)
	return app.Run()
}

func (info *InstallInfo) GetPlatforms() ([]string, error) {
	dirs, err := ioutil.ReadDir(filepath.Join(info.Path, "PlaybackEngines"))
	if err != nil {
		return nil, err
	}

	platforms := make([]string, 0, len(dirs))
	for _, platform := range dirs {
		name := strings.TrimSuffix(filepath.Base(platform.Name()), "Support")
		platforms = append(platforms, name)
	}

	return platforms, nil
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
	hubInstallPaths, _ := filepath.Glob("/Applications/Unity/Hub/Editor/**/Unity.app")
	unityPaths = append(unityPaths, hubInstallPaths...)

	installs := make([]*InstallInfo, 0, len(unityPaths))
	for _, path := range unityPaths {
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

	installData := InstallInfo{Version: VersionFromString(appInfo.CFBundleVersion), Path: path}
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
	newName := fmt.Sprintf("%s", install.Version.String())
	newPath := filepath.Join(unityHubEditorPath, newName)

	if oldPath == newPath {
		return nil
	}

	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		if err := os.Mkdir(newPath, 0755); err != nil {
			return err
		}
	}

	files, err := ioutil.ReadDir(oldPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name() == "Hub" || file.Name() == ".DS_Store" {
			continue
		}

		oldFilePath := filepath.Join(oldPath, file.Name())
		newFilePath := filepath.Join(newPath, file.Name())

		if err := os.Rename(oldFilePath, newFilePath); err != nil {
			return err
		}
	}

	if oldPath != baseUnityPath {
		if err := os.RemoveAll(oldPath); err != nil {
			return err
		}
	}

	return nil
}

func InstallToUnityDir(install *InstallInfo) error {
	oldPath := filepath.Dir(install.Path)
	newPath := baseUnityPath // TODO switch out for unity-hub path

	if oldPath == newPath {
		return nil
	}

	files, err := ioutil.ReadDir(oldPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}

		oldFilePath := filepath.Join(oldPath, file.Name())
		newFilePath := filepath.Join(newPath, file.Name())

		if err := os.Rename(oldFilePath, newFilePath); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(oldPath); err != nil {
		return err
	}

	return nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
