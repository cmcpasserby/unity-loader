package unity

import (
	"fmt"
	"howett.net/plist"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	baseUnityPath = "Applications/Unity"
)

func (info *InstallInfo) RunWithTarget(project, target string) error {
	absProject, _ := filepath.Abs(project)

	var app *exec.Cmd
	if target == "" {
		app = exec.Command("open", "-a", info.Path, "-n", "--args", "-projectPath", absProject)
	} else {
		app = exec.Command("open", "-a", info.Path, "-n", "--args", "-projectPath", absProject, "-buildTarget", target)
	}

	return app.Run()
}

func (info *InstallInfo) NewProject(projectName string) error {
	projectPath, _ := filepath.Abs(projectName)
	app := exec.Command("open", "-a", info.Path, "--args", "-createProject", projectPath)
	fmt.Printf("Creating project at %s\n", projectPath)
	return app.Run()
}

// TODO do we even need this?
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
