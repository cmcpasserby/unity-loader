package unity

import (
	"fmt"
	"github.com/gonutz/w32/v2"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetInstallFromPath(path string) (InstallInfo, error) {
	size := w32.GetFileVersionInfoSize(path)
	if size <= 0 {
		return InstallInfo{}, fmt.Errorf("GetFUleVersionInfoSIzeFiaed")
	}

	info := make([]byte, size)
	if ok := w32.GetFileVersionInfo(path, info); !ok{
		return InstallInfo{}, fmt.Errorf("GetFileVersionInfo Failed")
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
	return filepath.Glob(fmt.Sprintf("%s/**/Editor/Unity.exe", searchPath))
}

func command(path string, args ...string) *exec.Cmd {
	return exec.Command(path, args...)
}
