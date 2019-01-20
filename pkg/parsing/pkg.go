package parsing

import (
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"strings"
)

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

type PkgGeneric interface {
	PkgName() string
	Md5() string
	Size() int
	IsModule() bool
	IsPkgFile() bool
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

func (pkg *Pkg) PkgName() string {
	return pkg.Version
}

func (pkg *Pkg) Md5() string {
	return pkg.Checksum
}

func (pkg *Pkg) Size() int {
	return pkg.DownloadSize
}

func (pkg *Pkg) IsModule() bool {
	return false
}

func (pkg *Pkg) IsPkgFile() bool {
	return true
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

func (pkg *PkgModule) PkgName() string {
	return pkg.Name
}

func (pkg *PkgModule) Md5() string {
	return pkg.Checksum
}

func (pkg *PkgModule) Size() int {
	return pkg.DownloadSize
}

func (pkg *PkgModule) IsModule() bool {
	return true
}

func (pkg *PkgModule) IsPkgFile() bool {
	return strings.HasSuffix(pkg.DownloadUrl, ".pkg")
}
