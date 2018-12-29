package parsing

import (
	"encoding/json"
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"log"
	"net/http"
)

const hubUrl = "https://public-cdn.cloud.unity3d.com/hub/prod/releases-darwin.json"

type PkgDetailsSlice []PkgDetails

func (pkg PkgDetailsSlice) Len() int {
	return len(pkg)
}

func (pkg PkgDetailsSlice) Less(i, j int) bool {
	a := unity.VersionDataFromString(pkg[i].Version)
	b := unity.VersionDataFromString(pkg[j].Version)
	return unity.VersionLess(a, b)
}

func (pkg PkgDetailsSlice) Swap(i, j int) {
	pkg[i], pkg[j] = pkg[j], pkg[i]
}

type Releases struct {
	Official PkgDetailsSlice `json:"official"`
	Beta     PkgDetailsSlice `json:"beta"`
}

type PkgDetails struct {
	Version       string      `json:"version"`
	Lts           bool        `json:"lts"`
	DownloadUrl   string      `json:"downloadUrl"`
	DownloadSize  int         `json:"downloadSize"`
	InstalledSize int         `json:"installedSize"`
	Checksum      string      `json:"checksum"`
	Modules       []PkgModule `json:"modules"`
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

func GetHubVersions() (*Releases, error) {
	resp, err := http.Get(hubUrl)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var data Releases
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *Releases) Filter(f func(PkgDetails) bool) []PkgDetails {
	result := make([]PkgDetails, 0)

	for _, pkg := range r.Official {
		if f(pkg) {
			result = append(result, pkg)
		}
	}

	for _, pkg := range r.Beta {
		if f(pkg) {
			result = append(result, pkg)
		}
	}

	return result
}

func (r *Releases) First(f func(PkgDetails) bool) PkgDetails {
	return r.Filter(f)[0]
}
