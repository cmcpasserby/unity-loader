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
    downloadMatchRe = `(https?://[\w/.-]+/[0-9a-f]{12}/)[\w/.-]+-(\d+\.\d+\.\d+\w\d+)(?:\.dmg|\.pkg)`
    versionMatchRe = `(\d+\.\d+\.\d+\w\d+)`
    hashMatchRe = `[0-9a-f]{12}`
)


type VersionData struct {
    VersionString string
    Packages []string
}

func ParseVersions(url string, version string) {
    downloadRe := regexp.MustCompile(downloadMatchRe)
    versionRe := regexp.MustCompile(versionMatchRe)

    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
        return
    }
    defer response.Body.Close()

    contents, _ := ioutil.ReadAll(response.Body)
    matches := downloadRe.FindAllString(string(contents), -1)

    packages := make([]string, 0, 5)
    for _, m := range matches {
        ver := versionRe.FindString(m)

        if ver == version {
            packages = append(packages, m)
        }
    }

    for _, url := range packages {
        fmt.Println(url)
    }
}
