package unity

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// InstallInfo represents a runnable unity install
type InstallInfo struct {
	Path    string
	Version VersionData
}

// Run launches this Unity installs with a given project
func (info *InstallInfo) Run(project string) error {
	return info.RunWithTarget(project, "")
}

// RunWithTarget launches this unity install with the given project and target
func (info *InstallInfo) RunWithTarget(project, target string) error {
	absProject, _ := filepath.Abs(project)

	args := []string{"-projectPath", absProject}
	if target != "" {
		args = append(args, "-buildTarget", target)
	}

	return command(info.Path, args...).Start()
}

// String prints version and path for this InstallInfo
func (info *InstallInfo) String() string {
	return fmt.Sprintf("Version: %s Path: \"%s\"", info.Version.String(), info.Path)
}

// GetInstallFromPath returns an InstallInfo for a given path
func GetInstallFromPath(path string) (InstallInfo, error) {
	return getFromInstallPathInternal(path)
}

// GetVersionFromProject finds the Unity version used in a given project path
func GetVersionFromProject(projectPath string) (VersionData, error) {
	versionFile := filepath.Join(projectPath, "ProjectSettings", "ProjectVersion.txt")
	file, err := os.Open(versionFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return VersionData{}, fmt.Errorf("\"%s\" is not a valid unity project", projectPath)
		}
		return VersionData{}, err
	}
	defer closeFile(file)

	versionData, err := readProjectVersion(file)
	if err != nil {
		return VersionData{}, err
	}
	return versionData, nil
}

// GetInstalls lists all found Unity installs for a given set of search paths
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

// GetInstallFromVersion tries to find the appropriate Unity install for a given version
func GetInstallFromVersion(version VersionData, searchPaths ...string) (InstallInfo, error) {
	installs, err := GetInstalls(searchPaths...)
	if err != nil {
		return InstallInfo{}, err
	}

	for _, install := range installs {
		if install.Version.Compare(version) == 0 {
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
