package parsing

import (
	"github.com/cmcpasserby/unity-loader/pkg/unity"
	"io/ioutil"
	"net/http"
	"regexp"
)

const archiveUrl = "https://unity3d.com/get-unity/download/archive"

var downloadRe = regexp.MustCompile(`(https?://[\w/.-]+/[0-9a-f]{12}/)[\w/.-]+-(\d+\.\d+\.\d+\w\d+)(?:\.dmg|\.pkg)`)
var versionRe = regexp.MustCompile(`(\d+\.\d+\.\d+\w\d+)`)
var uuidRe = regexp.MustCompile(`[0-9a-f]{12}`)

func getArchiveVersionData() ([]unity.ExtendedVersionData, error) {
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
		versions = append(versions, verData)
	}

	return versions, nil
}
