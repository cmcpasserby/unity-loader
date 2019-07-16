package commands

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/sudoer"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const modulesFilename = "modules.json"

var moduleIdRe = regexp.MustCompile(`{(.*?)}`)

type downloadedModule struct {
	parsing.PkgModule
	ModulePath string
}

func install(args ...string) error {
	if err := repairPaths(true); err != nil {
		return err
	}

	cache, err := settings.ReadCache()
	if err != nil {
		return err
	}

	if cache.NeedsUpdate() {
		if err := cache.Update(); err != nil {
			return err
		}
	}

	currentInstalls, err := unity.GetInstalls()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		releases := cache.Releases.Filter(func(cacheVersion parsing.CacheVersion) bool {
			for _, install := range currentInstalls {
				if install.Version.String() == cacheVersion.String() {
					return false
				}
			}
			return true
		})
		sort.Sort(sort.Reverse(releases))

		versionStrs := make([]string, 0, len(releases))
		for _, ver := range releases {
			versionStrs = append(versionStrs, ver.String())
		}

		prompt := &survey.Select{
			Message:  "select version to install:",
			Options:  versionStrs,
			PageSize: 10,
		}

		var result string
		if err := survey.AskOne(prompt, &result, nil); err != nil {
			return err
		}

		selectedVersion := releases.First(func(details parsing.CacheVersion) bool {
			return details.String() == result
		})

		if err := installVersion(*selectedVersion, false); err != nil {
			return err
		}
		return nil
	} else {
		// TODO handle passed in version number
	}

	return nil
}

func installVersion(version parsing.CacheVersion, modulesOnly bool) error {
	config, err := settings.ParseDotFile()
	if err != nil {
		return err
	}

	if modulesOnly {
		if _, err := unity.GetInstallFromVersion(version.String()); err != nil {
			return err
		}
	} else {
		if _, err := unity.GetInstallFromVersion(version.String()); err == nil {
			return errors.New("version is already installed") // TODO replace with proper error in future
		}
	}

	sudo := new(sudoer.Sudoer)
	if err := sudo.AskPass(); err != nil {
		return err
	}

	installInfo, err := version.GetPkg()
	if err != nil {
		return err
	}

	titles := make([]string, 0, len(installInfo.Modules))
	defaults := make([]string, 0, len(installInfo.Modules))

	for _, module := range installInfo.Modules {
		if !module.Visible {
			continue
		}

		moduleString := fmt.Sprintf("%s {%s}", module.Name, module.Id)

		titles = append(titles, moduleString)

		for _, item := range config.ModuleDefaults {
			if item == module.Id {
				defaults = append(defaults, moduleString)
				break
			}
		}
	}

	prompt := &survey.MultiSelect{
		Message:  "select modules to install",
		Options:  titles,
		Default:  defaults,
		PageSize: 10,
	}

	var results []string
	if err := survey.AskOne(prompt, &results, nil); err != nil {
		return err
	}

	selectedModulesCount := 0
	for _, mod := range installInfo.Modules {
		matchFound := false

		for _, resultStr := range results {
			modId := moduleIdRe.FindStringSubmatch(resultStr)[1]
			if modId == mod.Id {
				matchFound = true
				mod.Selected = true
				selectedModulesCount++
				break
			}
		}
		if !matchFound {
			mod.Selected = false
		}
	}

	defer cleanUp()

	var unityPath string
	if !modulesOnly {
		unityPath, err := downloadPkg(&installInfo)
		if err != nil {
			return err
		}

		isValid, err := validate(&installInfo, unityPath)
		if err != nil {
			return err
		}
		if !isValid {
			return fmt.Errorf("%q was not a valid package, try installing again\n", installInfo.Version)
		}
	} else {
		install, err := unity.GetInstallFromVersion(version.String())
		if err != nil {
			return err
		}

		if err := unity.InstallToUnityDir(install); err != nil {
			return err
		}
	}

	modulePaths := make([]downloadedModule, 0, selectedModulesCount)

	for _, module := range installInfo.Modules {
		if !module.Selected {
			continue
		}

		modPath, err := downloadModule(&module)
		if err != nil {
			return err
		}

		isValid, err := validate(&module, modPath)
		if err != nil {
			return err
		}

		if !isValid {
			return fmt.Errorf("%q was not a valid package, try installing again\n", module.Name)
		}
		modulePaths = append(modulePaths, downloadedModule{module, modPath})
	}

	if err := installEditor(&installInfo, unityPath, sudo); err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)

	installedInfo, err := unity.GetInstallFromVersion(version.String())
	if err != nil {
		// TODO maybe make this error more user friendly for this use case
		return err
	}
	baseUnityPath := filepath.Dir(installedInfo.Path)

	for _, modPath := range modulePaths {
		if strings.HasSuffix(strings.ToLower(modPath.DownloadUrl), ".pkg") {
			if err := installModule(&modPath, sudo); err != nil {
				return err
			}
		} else if strings.HasSuffix(strings.ToLower(modPath.DownloadUrl), ".zip") {
			if err := installZip(&modPath, baseUnityPath, sudo); err != nil {
				return err
			}
		} else if strings.HasSuffix(strings.ToLower(modPath.DownloadUrl), ".dmg") {
			if err := installDmg(&modPath, baseUnityPath, sudo); err != nil {
				return err
			}
		}
	}

	time.Sleep(500 * time.Millisecond)

	// TODO create modules list here and dump to json
	if err := writeModulesFile(baseUnityPath, installInfo.Modules); err != nil {
		return err
	}

	if installInfo, err := unity.GetInstallFromVersion(version.String()); err == nil {
		if err := unity.RepairInstallPath(installInfo); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func download(url, name string, size int) (string, error) {
	pkgPath, err := settings.GetPkgPath()
	if err != nil {
		return "", err
	}

	downloadPath := filepath.Join(pkgPath, filepath.Base(url))

	out, err := os.Create(downloadPath)
	if err != nil {
		return "", err
	}
	defer closeFile(out)

	done := make(chan int64)
	go downloadProgress(size, name, downloadPath, done)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer closeResponse(resp)

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	done <- n

	return downloadPath, nil
}

// TODO can combine these 2 to use the PkgGeneric
func downloadPkg(pkg *parsing.Pkg) (string, error) {
	downloadPath, err := download(pkg.DownloadUrl, pkg.Version, pkg.DownloadSize)
	if err != nil {
		return "", err
	}
	return downloadPath, nil
}

func downloadModule(mod *parsing.PkgModule) (string, error) {
	downloadPath, err := download(mod.DownloadUrl, mod.Name, mod.DownloadSize)
	if err != nil {
		return "", err
	}
	return downloadPath, nil
}

func downloadProgress(downloadSize int, name, path string, done <-chan int64) {
	stop := false

	bar := pb.New64(int64(downloadSize))
	bar.Prefix(name)
	bar.ShowSpeed = true
	bar.Width = 120
	bar.SetUnits(pb.U_BYTES)
	bar.Start()

	for {
		select {
		case <-done:
			stop = true
		default:
			fi, err := os.Stat(path)
			if err != nil {
				log.Fatal(err)
			}

			size := fi.Size()
			if size == 0 {
				size = 1
			}

			bar.Set64(size)
		}
		if stop {
			bar.Set64(int64(downloadSize))
			bar.FinishPrint(fmt.Sprintf("Downloaded %q", name))
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func validate(pkg parsing.PkgGeneric, path string) (bool, error) {
	fmt.Printf("Validating pacakge %q...", pkg.PkgName())

	isValid := false

	if pkg.Md5() == "" {
		f, err := os.Open(path)
		if err != nil {
			return false, err
		}
		defer closeFile(f)

		fi, err := f.Stat()
		if err != nil {
			return false, err
		}

		isValid = fi.Size() == int64(pkg.Size())
	} else {
		f, err := os.Open(path)
		if err != nil {
			return false, err
		}
		defer closeFile(f)

		hash := md5.New()

		_, err = io.Copy(hash, f)
		if err != nil {
			return false, err
		}

		sum := hash.Sum(nil)
		isValid = hex.EncodeToString(sum) == pkg.Md5()
	}

	fmt.Print("\033[2K") // clears current line
	if isValid {
		fmt.Printf("\rPackage %q is valid\n", pkg.PkgName())
	} else {
		fmt.Printf("\rPackage %q is not valid\n", pkg.PkgName())
	}

	return isValid, nil
}

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

func writeModulesFile(baseUnityPath string, modules []parsing.PkgModule) error {
	modsFilePath := filepath.Join(baseUnityPath, modulesFilename)
	f, err := os.Create(modsFilePath)
	if err != nil {
		return err
	}
	defer closeFile(f)

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")

	return enc.Encode(modules)
}

func cleanUp() {
	downloadPath, err := settings.GetPkgPath() // todo will need to make this per unity version later on, as to not break multiple installs at once
	if err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll(downloadPath); err != nil {
		log.Fatal(err)
	}
}

func closeResponse(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
