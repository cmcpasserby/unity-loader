package parsing

import (
	"encoding/json"
	"log"
	"net/http"
)

const hubUrl = "https://public-cdn.cloud.unity3d.com/hub/prod/releases-darwin.json"

type Package interface {
	Download() (string, error)
	Validate(path string) (bool, error)
	Install(path string) error
}

type Releases struct {
	Official []PkgDetails `json:"official"`
	Beta     []PkgDetails `json:"beta"`
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

func (pkg *PkgDetails) Download() (string, error) {
	return "", nil
}

func (pkg *PkgDetails) Validate(path string) (bool, error) {
	return false, nil
}

func (pkg *PkgDetails) Install(path string) error {
	return nil
}

func (pkg *PkgDetails) downloadProgress(done chan int64) {
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

func (r Releases) Filter(f func(PkgDetails) bool) []PkgDetails {
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
