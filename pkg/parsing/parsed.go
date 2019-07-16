package parsing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"gopkg.in/ini.v1"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	configName        = "unity-%s-osx.ini"
	unitySection      = "Unity"
	unityPathFragment = "{UNITY_PATH}"
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

	firstVersionWithDocsCategory = unity.VersionFromString("2018.2.0a1")
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
			pkg.Lts = v.Major == 4 // TODO might need to find a better way but will just check if its a x.4 release for now
			pkg.DownloadUrl = currentUrl + iniData.Path
			pkg.DownloadSize = int(iniData.Size)
			pkg.InstalledSize = int(iniData.InstalledSize)
			pkg.Checksum = iniData.Md5
			pkg.Modules = make([]PkgModule, 0, len(sectionNames)-1)
		} else {
			pkg.Modules = append(pkg.Modules, PkgModule{
				Id:            section,
				Name:          iniData.Title,
				Description:   iniData.Description,
				DownloadUrl:   currentUrl + iniData.Path,
				Category:      getCategory(section, v),
				DownloadSize:  int(iniData.Size),
				InstalledSize: int(iniData.InstalledSize),
				Checksum:      iniData.Md5,
				Destination:   getDestination(section),
				Visible:       !iniData.Hidden,
				Selected:      iniData.Install,
			})
		}
	}

	return pkg, nil
}

func (v *CacheVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s:%s", v.String(), v.RevisionHash))
}

func (v *CacheVersion) UnmarshalJSON(data []byte) error {
	var dataString string

	if err := json.Unmarshal(data, &dataString); err != nil {
		return err
	}

	split := strings.Split(dataString, ":")

	v.ExtendedVersionData = unity.ExtendedVersionData{
		VersionData:  unity.VersionFromString(split[0]),
		RevisionHash: split[1],
	}

	return nil
}

type CacheVersionSlice []CacheVersion

func (s CacheVersionSlice) Len() int {
	return len(s)
}

func (s CacheVersionSlice) Less(i, j int) bool {
	return s[i].VersionData.Compare(s[j].VersionData) < 0
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

func (s CacheVersionSlice) First(f func(CacheVersion) bool) *CacheVersion {
	versions := s.Filter(f)
	if len(versions) == 0 {
		return nil
	}
	return &versions[0]
}

func (s CacheVersionSlice) Any(f func(CacheVersion) bool) bool {
	for _, ver := range s {
		if f(ver) {
			return true
		}
	}
	return false
}

func getDestination(componentId string) string {
	switch componentId {
	case "mono":
	case "visualstudio":
		return ""
	case "monodevelop":
	case "documentation":
		return unityPathFragment
	case "standardassets":
		return filepath.Join(unityPathFragment, "Standard Assets")
	case "exampleprojects":
	case "example":
		return "/Users/Shared/Unity"
	case "android":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer")
	case "android-sdk-build-tools":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer/SDK/build-tools")
	case "android-sdk-platforms":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer/SDK/platforms")
	case "android-sdk-platform-tools":
	case "android-sdk-ndk-tools":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer/SDK")
	case "android-ndk":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer/NDK")
	case "android-open-jdk":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AndroidPlayer/OpenJDK")
	case "ios":
		return filepath.Join(unityPathFragment, "PlaybackEngines")
	case "tvos":
	case "appletv":
		return filepath.Join(unityPathFragment, "PlaybackEngines/AppleTVSupport")
	case "linux":
		return filepath.Join(unityPathFragment, "PlaybackEngines/LinuxStandaloneSupport")
	case "mac":
	case "mac-il2cpp":
		return filepath.Join(unityPathFragment, "Unity.app/Contents/PlaybackEngines/MacStandaloneSupport")
	case "samsungtv":
	case "samsung-tv":
		return filepath.Join(unityPathFragment, "PlaybackEngines/STVPlayer")
	case "tizen":
		return filepath.Join(unityPathFragment, "PlaybackEngines/TizenPlayer")
	case "vuforia":
	case "vuforia-ar":
		return filepath.Join(unityPathFragment, "PlaybackEngines/VuforiaSupport")
	case "webgl":
		return filepath.Join(unityPathFragment, "PlaybackEngines/WebGLSupport")
	case "windows":
	case "windows-mono":
		return filepath.Join(unityPathFragment, "PlaybackEngines/WindowsStandaloneSupport")
	case "facebook":
	case "facebook-games":
		return filepath.Join(unityPathFragment, "PlaybackEngines/Facebook")
	case "facebookgameroom":
		return ""
	case "lumin":
		return filepath.Join(unityPathFragment, "PlaybackEngines/LuminSupport")
	}

	if strings.HasPrefix(componentId, "language-") {
		return filepath.Join(unityPathFragment, "Unity.app/Contents/Localization")
	}

	return unityPathFragment
}

func getCategory(componentId string, version *CacheVersion) string {
	switch componentId {
	case "monodevelop":
	case "visualstudio":
		return "Dev tools"
	case "mono":
	case "visualstudioprofessionalunityworkload":
	case "visualstudioenterpriseunityworkload":
	case "facebookgameroom":
		return "Plugins"

	case "standardassets":
	case "exampleprojects":
	case "example":
		return "Components"

	case "documentation":
		return "Components" // TODO fix for docs removed in later versions
	}

	if strings.HasPrefix(componentId, "language-") {
		return "Language Packs (Preview)"
	}
	return "Platforms"
}
