package unity

import (
	"path/filepath"
	"sort"
)

func (info *InstallInfo) RunWithTarget(project, target string) error {
	panic("not implemented yet")
}

func (info *InstallInfo) NewProject(projectName string) error {
	panic("not implemented yet")
}

func GetInstalls() ([]InstallInfo, error) { // TODO might be able to make this not per platform and just provide the glob pattern
	unityPaths, err := filepath.Glob("C:/Program Files/Unity/Hub/Editor/**/Editor/Unity.exe") // TODO should be user configurable
	if err != nil {
		return nil, err
	}

	installs := make([]InstallInfo, 0, len(unityPaths))
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

func GetInstallFromPath(path string) (InstallInfo, error) {
	panic("not implemented yet")
}
