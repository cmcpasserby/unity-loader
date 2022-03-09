package unity

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

const archiveUrl = "https://unity3d.com/get-unity/download/archive"

var hubUrlRe = regexp.MustCompile(`unityhub://(\d+\.\d+\.\d+[pfba]\d+)/([0-9a-f]{12})`)

func GetAllVersions() ([]VersionData, error) {
	resp, err := http.Get(archiveUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	matches := hubUrlRe.FindAllSubmatch(body, -1)
	versions := make([]VersionData, len(matches))

	for i, m := range matches {
		versions[i] = VersionFromString(string(m[1]))
		versions[i].RevisionHash = string(m[2])
	}

	return versions, nil
}

func SearchArchive(partialVersion string) ([]VersionData, error) {
	allVersions, err := GetAllVersions()
	if err != nil {
		return nil, err
	}

	results := make([]VersionData, 0)
	for _, ver := range allVersions {
		if strings.HasPrefix(ver.String(), partialVersion) {
			results = append(results, ver)
		}
	}

	return results, nil
}

// func InstallUnity(data VersionData) (InstallInfo, error) {
// }
