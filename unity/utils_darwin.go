package unity

import (
	"fmt"
	"howett.net/plist"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

type appInfoDict struct {
	CFBundleVersion string `plist:"CFBundleVersion"`
}

func GetInstalls() ([]*InstallInfo, error) {
	unityPaths, err := filepath.Glob("/Applications/Unity/Hub/Editor/**/Unity.app") // TODO should be user configurable
	if err != nil {
		return nil, err
	}

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
