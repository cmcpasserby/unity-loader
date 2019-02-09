package parsing

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type Releases struct {
	Official PkgSlice `json:"official"`
	Beta     PkgSlice `json:"beta"`
}

var hubUrlRe = regexp.MustCompile(`(unityhub://(\d+\.\d+\.\d+\w\d+)/[0-9a-f]{12})`)

const hubUrl = "https://public-cdn.cloud.unity3d.com/hub/prod/releases-darwin.json"

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

func (r *Releases) Filter(f func(Pkg) bool) PkgSlice {
	result := make([]Pkg, 0)

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

func (r *Releases) First(f func(Pkg) bool) Pkg {
	return r.Filter(f)[0]
}
