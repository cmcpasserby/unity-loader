package unity

import (
	"fmt"
	"howett.net/plist"
	"os"
	"os/exec"
	"path/filepath"
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

func unityGlob(searchPath string) ([]string, error) {
	return filepath.Glob(fmt.Sprintf("%s/**/Unity.app", searchPath))
}

func GetInstallFromPath(path string) (InstallInfo, error) {
	plistPath := filepath.Join(path, "Contents/info.plist")
	file, err := os.Open(plistPath)
	if err != nil {
		return InstallInfo{}, err
	}
	defer closeFile(file)

	var appInfo appInfoDict
	if err = plist.NewDecoder(file).Decode(&appInfo); err != nil {
		return InstallInfo{}, err
	}

	installData := InstallInfo{Path: path, Version: VersionFromString(appInfo.CFBundleVersion)}
	return installData, nil
}
