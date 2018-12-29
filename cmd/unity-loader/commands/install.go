package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
	"github.com/cmcpasserby/unity-loader/pkg/sudoer"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/cheggaaa/pb.v1"
	"log"
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
		if module.Selected {
			defaults = append(defaults, moduleString)
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

	isValid, err := validatePkg(&installInfo, unityPath)
	if err != nil {
		return err
	}
	if !isValid {
		return fmt.Errorf("%q was not a valid package, try installing again\n", installInfo.Version)
	}

	for _, module := range selected {
		fmt.Println(module.Id)
		// TODO download selected modules
	}

	return nil
}

func downloadPkg(pkg *parsing.PkgDetails) (string, error) {
	pkgPath, err := settings.GetPkgPath()
	if err != nil {
		return "", err
	}

	downloadPath := filepath.Join(pkgPath, filepath.Base(pkg.DownloadUrl))

	out, err := os.Create(downloadPath)
	if err != nil {
		return "", err
	}
	defer closeFile(out)

	done := make(chan int64)
	go downloadProgress(pkg.DownloadSize, pkg.Version, downloadPath, done)
	// TODO start download

	return downloadPath, nil
}

func validatePkg(pkg *parsing.PkgDetails, path string) (bool, error) {
	return true, nil
}

func downloadModule(mod *parsing.PkgModule) (string, error) {
}

func validateModule(mod *parsing.PkgModule, path string) (bool, error) {
}

func download() error {
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
		}
		time.Sleep(time.Second * 10)
	}
}

func cleanUp(downloadPath string) {
	if err := os.RemoveAll(downloadPath); err != nil {
		log.Fatal(err)
	}
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
