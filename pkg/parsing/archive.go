package parsing

import (
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"io/ioutil"
	"net/http"
	"regexp"
)


const archiveUrl   = "https://unity3d.com/get-unity/download/archive"

var (
	unityHubRe = regexp.MustCompile(`(unityhub://(\d+\.\d+\.\d+\w\d+)/[0-9a-f]{12})`)
	downloadRe = regexp.MustCompile(`(https?://[\w/.-]+/[0-9a-f]{12}/)[\w/.-]+-(\d+\.\d+\.\d+\w\d+)(?:\.dmg|\.pkg)`)
	versionRe  = regexp.MustCompile(`(\d+\.\d+\.\d+\w\d+)`)
	revisionHashRe = regexp.MustCompile(`[0-9a-f]{12}`)
)

func GetVersions() ([]CacheVersion, error) {
	versions := make([]CacheVersion, 0)

	resp, err := http.Get(archiveUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	matches := unityHubRe.FindAllString(string(contents), -1)
	matches = append(matches, downloadRe.FindAllString(string(contents), -1)...)
	dupMap := make(map[unity.VersionData]bool)

	for _, match := range matches {
		verStr := versionRe.FindString(match)
		revisionHash := revisionHashRe.FindString(match)
		verData := unity.ExtendedVersionData{
			VersionData: unity.VersionDataFromString(verStr),
			RevisionHash: revisionHash,
		}

		if _, value := dupMap[verData.VersionData]; !value {
			dupMap[verData.VersionData] = true
			versions  = append(versions, CacheVersion{verData})
		}
	}

	return versions, nil
}
