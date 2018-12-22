package main

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/cmd/unity-loader/settings"
	"path/filepath"
)

func cleanup(args ...string) error {
	var path string

	if len(args) == 0 {
		dotFile, err := settings.ParseDotFile()
		if err != nil {
			return err
		}
		path = dotFile.ProjectFolder
	} else {
		path = args[0]
	}

	versions, err := getUsedVersions(path)
	if err != nil {
		return err
	}

	for _, path := range versions {
		fmt.Println(path)
	}

	return nil
}

func getUsedVersions(path string) ([]string, error) {
	globString := filepath.Join(path, "**/ProjectSettings/ProjectVersion.txt")
	projectPaths, _ := filepath.Glob(globString)

	versions := make([]string, 0, len(projectPaths))

	for _, path := range projectPaths {
		path = filepath.Dir(filepath.Dir(path))
		versions = append(versions, path)
	}

	return versions, nil
}
