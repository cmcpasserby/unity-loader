package commands

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/sudoer"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

var moduleIdRe = regexp.MustCompile(`{(.*?)}`)

func install(args ...string) error {
	// TODO check cache timestamp and maybe update

	cache, err := settings.ReadCache()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		sort.Sort(sort.Reverse(cache.Releases.Official))

		versionStrs := make([]string, 0, len(cache.Releases.Official))
		for _, ver := range cache.Releases.Official {
			versionStrs = append(versionStrs, ver.Version)
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

		if err := installVersion(result, cache); err != nil {
			return err
		}
		return nil
	}

	return nil
}

func installVersion(version string, cache *settings.Cache) error {
	config, err := settings.ParseDotFile()
	if err != nil {
		return err
	}

	sudo := new(sudoer.Sudoer)
	if err := sudo.AskPass(); err != nil {
		return err
	}

	installInfo := cache.Releases.First(func(details parsing.PkgDetails) bool {
		return details.Version == version
	})

	titles := make([]string, 0, len(installInfo.Modules))
	defaults := make([]string, 0, len(installInfo.Modules))

	for _, module := range installInfo.Modules {
		moduleString := fmt.Sprintf("%s {%s}", module.Name, module.Id)

		titles = append(titles, moduleString)

		if value, ok := config.ModuleDefaults[module.Id]; ok {
			if value {
				defaults = append(defaults, moduleString)
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

	selected := installInfo.FilterModules(func(mod parsing.PkgModule) bool {
		for _, resultStr := range results {
			modId := moduleIdRe.FindStringSubmatch(resultStr)[1]
			if modId == mod.Id {
				return true
			}
		}
		return false
	})

	unityPath, err := downloadPkg(&installInfo)
	if err != nil {
		return err
	}

	isValid, err := validate(installInfo.Version, installInfo.Checksum, unityPath)
	if err != nil {
		return err
	}
	if !isValid {
		return fmt.Errorf("%q was not a valid package, try installing again\n", installInfo.Version)
	}

	modulePaths := make([]string, 0, len(selected))

	for _, module := range selected {
		modPath, err := downloadModule(&module)
		if err != nil {
			return err
		}

		isValid, err := validate(module.Name, module.Checksum, modPath)
		if err != nil {
			return err
		}

		if !isValid {
			return fmt.Errorf("%q was not a valid package, try installing again\n", module.Name)
		}
		modulePaths = append(modulePaths, modPath)
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

func downloadPkg(pkg *parsing.PkgDetails) (string, error) {
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

func downloadProgress(downloadSize int, name, path string, done chan int64) {
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

func validate(name, checksum, path string) (bool, error) {
	fmt.Printf("Validating pacakge %q...", name)

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
	isValid := hex.EncodeToString(sum) == checksum

	fmt.Print("\033[2K") // clears current line
	if isValid {
		fmt.Printf("\rPackage %q is valid\n", name)
	} else {
		fmt.Printf("\rPackage %q is not valid\n", name)
	}

	return isValid, nil
}

func cleanUp(downloadPath string) {
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
