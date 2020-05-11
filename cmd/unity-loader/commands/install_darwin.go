package commands

import (
	"errors"
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/sudoer"
	"os"
	"strings"
)

func installEditor(pkg *parsing.Pkg, unityPath string, sudo *sudoer.Sudoer) error {
	// TODO use xar and untar instead of running the installer package and apply needed copying

	if unityPath == "" {
		return errors.New("no downloaded package to install")
	}

	fmt.Printf("Installing Unity %q...", pkg.Version)

	if err := sudo.RunAsRoot("installer", "-package", unityPath, "-target", "/"); err != nil {
		return err
	}

	if err := os.Remove(unityPath); err != nil {
		return err
	}

	fmt.Print("\033[2k") // clears current line
	fmt.Printf("\rInstalled Unity %q\n", pkg.Version)
	return nil
}

func installModule(pkg *downloadedModule, sudo *sudoer.Sudoer) error {
	// TODO use xar and untar instead of running the installer package

	if pkg.ModulePath == "" {
		return errors.New("no downloaded package to install")
	}

	fmt.Printf("Installing Unity %q...", pkg.PkgName())

	err := sudo.RunAsRoot("installer", "-package", pkg.ModulePath, "-target", "/")
	if err != nil {
		return err
	}

	if err := os.Remove(pkg.ModulePath); err != nil {
		return err
	}

	fmt.Print("\033[2K") // clears current line
	fmt.Printf("\rInstalled Unity %q\n", pkg.PkgName())
	return nil
}

func installOther(pkg *downloadedModule, unityPath string, sudo *sudoer.Sudoer) error {
	typeString := strings.TrimSuffix(pkg.Category, "s")

	fmt.Printf("Installing %s %q...", typeString, pkg.PkgName())

	targetPath := strings.Replace(pkg.Destination, "{UNITY_PATH}", unityPath, 1)

	if err := sudo.RunAsRoot("cp", pkg.ModulePath, targetPath); err != nil {
		return err
	}

	fmt.Print("\033[2K") // clears current line
	fmt.Printf("\rInstalled %s %q\n", typeString, pkg.PkgName())
	return nil
}

func installZip(pkg *downloadedModule, unityPath string, sudo *sudoer.Sudoer) error {
	typeString := strings.TrimSuffix(pkg.Category, "s")

	fmt.Printf("Installing %s %q...", typeString, pkg.PkgName())

	targetPath := strings.Replace(pkg.Destination, "{UNITY_PATH}", unityPath, 1)

	err := sudo.RunAsRoot("unzip", pkg.ModulePath, "-d", targetPath)
	if err != nil {
		return err
	}

	if err := os.Remove(pkg.ModulePath); err != nil {
		return err
	}

	fmt.Print("\033[2K") // clears current line
	fmt.Printf("\rInstalled %s %q\n", typeString, pkg.PkgName())
	return nil
}

func installDmg(pkg *downloadedModule, unityPath string, sudo *sudoer.Sudoer) error {
	fmt.Println(`Ignoring "dmg" support not implemented yet`)
	fmt.Printf("Path %q\n", pkg.ModulePath)
	return nil
}
