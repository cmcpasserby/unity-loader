package unity

import (
	"fmt"
	"github.com/gonutz/w32/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetInstallFromPath returns an InstallInfo for a given path
func GetInstallFromPath(path string) (InstallInfo, error) {
	size := w32.GetFileVersionInfoSize(path)
	if size <= 0 {
		return InstallInfo{}, fmt.Errorf("GetFileVersionInfoSize failed")
	}

	info := make([]byte, size)
	if ok := w32.GetFileVersionInfo(path, info); !ok {
		return InstallInfo{}, fmt.Errorf("GetFileVersionInfo failed")
	}

	translations, ok := w32.VerQueryValueTranslations(info)
	if !ok {
		return InstallInfo{}, fmt.Errorf("VerQueryValueTranslations failed")
	}

	if len(translations) == 0 {
		return InstallInfo{}, fmt.Errorf("no translations found")
	}

	t := translations[0]
	productVersion, ok := w32.VerQueryValueString(info, t, w32.ProductVersion)
	if !ok {
		return InstallInfo{}, fmt.Errorf("cannot get company name")
	}

	ver := strings.Split(productVersion, "_")[0]
	ver = strings.TrimSpace(ver)

	return InstallInfo{Path: path, Version: VersionFromString(ver)}, nil
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
