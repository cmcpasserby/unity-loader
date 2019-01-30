package parsing

import (
	"errors"
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type IniData struct {
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

const (
	archiveUrl   = "https://unity3d.com/get-unity/download/archive"
	configName   = "unity-%s-osx.ini"
	unitySection = "Unity"
)

var (
	downloadRe = regexp.MustCompile(`(https?://[\w/.-]+/[0-9a-f]{12}/)[\w/.-]+-(\d+\.\d+\.\d+\w\d+)(?:\.dmg|\.pkg)`)
	versionRe  = regexp.MustCompile(`(\d+\.\d+\.\d+\w\d+)`)
	uuidRe     = regexp.MustCompile(`[0-9a-f]{12}`)

	baseUrls = [...]string{
		"https://netstorage.unity3d.com/unity/%s/",
		"https://download.unity3d.com/download_unity/%s/",
		"https://beta.unity3d.com/download/%s/",
		"https://files.unity3d.com/bootstrapper/%s/",
	}

	ignoredSections = [...]string{
		"DEFAULT",
		"VisualStudio",
		"Mono",
	}
)

func GetArchiveVersions(filter func (version unity.VersionData) bool) error {
	versions, err := getArchiveVersionData(filter)
	if err != nil {
		return err
	}

	pkgs := make(PkgSlice, 0)

	for _, ver := range versions {
		if pkg, err := getInstallData(ver); err == nil {
			pkgs = append(pkgs, pkg)
		} else {
			continue
		}
	}

	fmt.Printf("%v\n", pkgs)

	return nil
}

func getArchiveVersionData(filter func (version unity.VersionData) bool) ([]unity.ExtendedVersionData, error) {
	versions := make([]unity.ExtendedVersionData, 0)

	resp, err := http.Get(archiveUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	matches := downloadRe.FindAllString(string(contents), -1)

	for _, match := range matches {
		verStr := versionRe.FindString(match)
		verUuid := uuidRe.FindString(match)
		verData := unity.ExtendedVersionData{
			VersionData: unity.VersionDataFromString(verStr),
			VersionUuid: verUuid,
		}

		if filter(verData.VersionData) {
			versions = append(versions, verData)
		}
	}

	return versions, nil
}

func getInstallData(versionData unity.ExtendedVersionData) (Pkg, error) {
	fileName := fmt.Sprintf(configName, versionData.String())

	var currentUrl string // save for building unity and module urls
	var resp *http.Response
	var err error

	for _, baseUrl := range baseUrls {
		currentUrl = fmt.Sprintf(baseUrl, versionData.VersionUuid) + fileName
		resp, err = http.Get(currentUrl)
		if err == nil {
			break
		}
	}

	if resp == nil || err != nil {
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

		iniData := new(IniData)

		if err := cfg.Section(section).MapTo(iniData); err != nil {
			return Pkg{}, err
		}

		if section == unitySection {

			version := iniData.Version
			if version == "" {
				version = strings.Replace(iniData.Title, "Unity ", "", -1)
			}

			pkg.Version = version
			pkg.Lts = false
			pkg.DownloadUrl = iniData.Path // TODO get url func
			pkg.DownloadSize = int(iniData.Size)
			pkg.InstalledSize = int(iniData.InstalledSize)
			pkg.Checksum = iniData.Md5
			pkg.Modules = make([]PkgModule, 0, len(sectionNames) - 1)
		} else {
			pkg.Modules = append(pkg.Modules, PkgModule{
				Id: section,
				Name: iniData.Title,
				Description: iniData.Description,
				DownloadUrl: iniData.Path,
				Category: "Archive",
				InstalledSize: int(iniData.InstalledSize),
				DownloadSize: int(iniData.Size),
				Checksum: iniData.Md5,
				Destination: "{UNITY_PATH}",
				Visible: !iniData.Hidden,
				Selected: iniData.Install,
			})
		}
	}
	return pkg, nil
}
