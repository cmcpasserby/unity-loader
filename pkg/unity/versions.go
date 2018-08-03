package unity

import (
    "net/http"
    "log"
    "io/ioutil"
    "regexp"
    "fmt"
)

const (
    UnityDownloads = "https://unity3d.com/get-unity/download/archive"
    UnityLtsDownloads = "https://unity3d.com/unity/qa/lts-releases"
    UnityPatches = "https://unity3d.com/unity/qa/patch-releases"
    UnityBetas = "https://unity3d.com/unity/beta-download"
)

const (
    macEditorInstaller = "http://netstorage.unity3d.com/unity/%s/MacEditorInstaller/Unity.pkg"
)

const (
    downloadMatchRe = `(https?://[\w/.-]+/[0-9a-f]{12}/)[\w/.-]+-(\d+\.\d+\.\d+\w\d+)(?:\.pkg)`
    versionMatchRe = `(\d+\.\d+\.\d+\w\d+)`
    uuidMatchRe = `[0-9a-f]{12}`
)


type VersionData struct {
    VersionString string
    VersionUuid string
}

func ParseVersions(url string) (map[string]VersionData, error) {
    downloadRe := regexp.MustCompile(downloadMatchRe)
    versionRe := regexp.MustCompile(versionMatchRe)
    uuidRe := regexp.MustCompile(uuidMatchRe)

    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    contents, _ := ioutil.ReadAll(response.Body)
    matches := downloadRe.FindAllString(string(contents), -1)

    versions := map[string]VersionData{}

    for _, url := range matches {
        ver := versionRe.FindString(url)

        if _, ok := versions[ver]; !ok {
            versions[ver] = VersionData{ver,uuidRe.FindString(url)}
        }
    }

    keys := make([]string, 0, len(versions))
    for k := range versions {
        keys = append(keys, k)
    }

    fmt.Printf(macEditorInstaller + "\n", versions["2018.1.9f1"].VersionUuid)

    return versions, nil
}
