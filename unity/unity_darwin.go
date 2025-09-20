package unity

import (
	"os"
	"os/exec"
	"path/filepath"

	"howett.net/plist"
)

type appInfoDict struct {
	CFBundleExecutable string `plist:"CFBundleExecutable"`
	CFBundleVersion    string `plist:"CFBundleVersion"`
}

func unmarshalAppInfo(path string) (appInfoDict, error) {
	plistPath := filepath.Join(path, "Contents", "Info.plist")
	f, err := os.Open(plistPath)
	if err != nil {
		return appInfoDict{}, err
	}
	defer closeFile(f)

	var appInfo appInfoDict
	if err = plist.NewDecoder(f).Decode(&appInfo); err != nil {
		return appInfoDict{}, err
	}

	return appInfo, nil
}

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

func command(path string, args ...string) (*exec.Cmd, error) {
	executablePath, err := binFromApp(path)
	if err != nil {
		return nil, err
	}
	return exec.Command(executablePath, args...), nil
}
