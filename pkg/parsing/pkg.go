package parsing

import "github.com/cmcpasserby/unity-loader/pkg/unity"

type PkgSlice []Pkg

func (pkg PkgSlice) Len() int {
	return len(pkg)
}

func (pkg PkgSlice) Less(i, j int) bool {
	a := unity.VersionDataFromString(pkg[i].Version)
	b := unity.VersionDataFromString(pkg[j].Version)
	return unity.VersionLess(a, b)
}

func (pkg PkgSlice) Swap(i, j int) {
	pkg[i], pkg[j] = pkg[j], pkg[i]
}

type Releases struct {
	Official PkgSlice `json:"official"`
	Beta     PkgSlice `json:"beta"`
}

type Pkg struct {
	Version       string      `json:"version"`
	Lts           bool        `json:"lts"`
	DownloadUrl   string      `json:"downloadUrl"`
	DownloadSize  int         `json:"downloadSize"`
	InstalledSize int         `json:"installedSize"`
	Checksum      string      `json:"checksum"`
	Modules       []PkgModule `json:"modules"`
}

func (pkg *Pkg) FilterModules(f func(mod PkgModule) bool) []PkgModule {
	result := make([]PkgModule, 0, len(pkg.Modules))

	for _, mod := range pkg.Modules {
		if f(mod) {
			result = append(result, mod)
		}
	}

	return result
}

type PkgModule struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DownloadUrl   string `json:"downloadUrl"`
	Category      string `json:"category"`
	InstalledSize int    `json:"installedSize"`
	DownloadSize  int    `json:"downloadSize"`
	Visible       bool   `json:"visible"`
	Selected      bool   `json:"selected"`
	Destination   string `json:"destination"`
	Checksum      string `json:"checksum"`
}
