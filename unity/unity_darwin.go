package unity

import (
	"os/exec"
	"path/filepath"
)

func binFromApp(path string) (string, error) {
	appInfo, err := unmarshalAppInfo(path)
	if err != nil {
		return "", err
	}

	return filepath.Join(path, "Contents", "MacOS", appInfo.CFBundleExecutable), nil
}

func unityGlob(searchPath string) ([]string, error) {
	return filepath.Glob(filepath.Join(searchPath, "**/Unity.app"))
}

func command(path string, args ...string) *exec.Cmd {
	newArgs := append([]string{path, "-W", "-n", "--args"}, args...)
	return exec.Command("open", newArgs...)
}
