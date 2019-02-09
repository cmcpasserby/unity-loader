package parsing

import (
	"errors"
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/ini.v1"
	"net/http"
)

const (
	configName   = "unity-%s-osx.ini"
	unitySection = "Unity"
)

var (
	baseUrls = [...]string{
		"https://download.unity3d.com/download_unity/%s/",
		"https://netstorage.unity3d.com/unity/%s/",
		"https://beta.unity3d.com/download/%s/",
		"https://files.unity3d.com/bootstrapper/%s/",
	}

	ignoredSections = [...]string{
		"DEFAULT",
		"VisualStudio",
		"Mono",
	}
)

type iniData struct {
	Title         string `ini:"title"`
	Description   string `ini:"description"`
	Path          string `ini:"url"`
	Install       bool   `ini:"install"`
	Size          int64  `ini:"size"`
	InstalledSize int64  `ini:"installedsize"`
	Version       string `ini:"version"`
	Md5           string `ini:"md5"`
	Hidden        bool   `ini:"hidden"`
	Extension     string `ini:"extension"`
	RequiresUnity bool   `ini:"requires_unity"`
}

type CacheVersion struct {
	unity.ExtendedVersionData
}

func (v *CacheVersion) GetPkg() (Pkg, error) {
	fileName := fmt.Sprintf(configName, v.String())

	var currentUrl string
	var resp *http.Response
	var err error

	for _, baseUrl := range baseUrls {
		currentUrl = fmt.Sprintf(baseUrl, v.RevisionHash)
		resp, err = http.Get(currentUrl + fileName)
		if err == nil && resp.StatusCode == 200 {
			break
		}
	}

	if resp == nil || resp.StatusCode != 200 || err != nil {
		return Pkg{}, errors.New("connection error")
	}

	defer resp.Body.Close()

	cfg, err := ini.Load(resp.Body)
	if err != nil {
		return Pkg{}, err
	}

	testIgnored := func(item string) bool {
		for _, name := range ignoredSections {
			if item == name {
				return true
			}
		}
		return false
	}

	sectionNames := cfg.SectionStrings()
	pkg := Pkg{}

	for _, section := range sectionNames {
		if testIgnored(section) {
			continue
		}

		iniData := new(iniData)

		if err := cfg.Section(section).MapTo(iniData); err != nil {
			return Pkg{}, err
		}

		if section == unitySection {
			pkg.Version = v.String()
			pkg.Lts = false
			pkg.DownloadUrl = currentUrl + iniData.Path
			pkg.DownloadSize = int(iniData.Size)
			pkg.InstalledSize = int(iniData.InstalledSize)
			pkg.Checksum = iniData.Md5
			pkg.Modules = make([]PkgModule, 0, len(sectionNames) - 1)
		} else {
			pkg.Modules = append(pkg.Modules, PkgModule{
				Id: section,
				Name: iniData.Title,
				Description: iniData.Description,
				DownloadUrl: currentUrl + iniData.Path,
				Category: "Archive",
				DownloadSize: int(iniData.Size),
				InstalledSize: int(iniData.InstalledSize),
				Checksum: iniData.Md5,
				Destination: "{UNITY_PATH}",
				Visible: !iniData.Hidden,
				Selected: iniData.Install,
			})
		}
	}

	return pkg, nil
}

type CacheVersionSlice []CacheVersion

func (s CacheVersionSlice) Len() int {
	return len(s)
}

func (s CacheVersionSlice) Less(i, j int) bool {
	return unity.VersionLess(s[i].VersionData, s[j].VersionData)
}

func (s CacheVersionSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s CacheVersionSlice) Filter(f func(CacheVersion) bool) CacheVersionSlice {
	result := make(CacheVersionSlice, 0)

	for _, ver := range s {
		if f(ver) {
			result = append(result, ver)
		}
	}

	return result
}

func (s CacheVersionSlice) First(f func(CacheVersion) bool) CacheVersion {
	return s.Filter(f)[0]
}
