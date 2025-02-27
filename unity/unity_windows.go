package unity

import (
	"os"
	"os/exec"
	"path/filepath"
)

func binFromApp(path string) (string, error) {
	return path, nil
}

func unityGlob(searchPath string) ([]string, error) {
	items, err := filepath.Glob(filepath.Join(searchPath, "**/Editor/Unity.exe"))
	if err != nil {
		return nil, err
	}

	directPath := filepath.Join(searchPath, "Editor/Unity.exe")
	if _, err = os.Stat(directPath); err == nil {
		items = append(items, directPath)
	}

	return items, nil
}

func command(path string, args ...string) *exec.Cmd {
	return exec.Command(path, args...)
}
