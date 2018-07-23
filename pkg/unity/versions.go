package unity

import (
    "net/http"
    "log"
    "io/ioutil"
    "fmt"
    "regexp"
)

const UnityDownloads = "https://unity3d.com/get-unity/download/archive"
const UnityLtsDownloads = "https://unity3d.com/unity/qa/lts-releases"
const UnityPatches = "https://unity3d.com/unity/qa/patch-releases"
const UnityBetas = "https://unity3d.com/unity/beta-download"

const downloadMatchRe = `(https?://[\w/.-]+/[0-9a-f]{12}/)[\w/.-]+-(\d+\.\d+\.\d+\w\d+)(?:\.dmg|\.pkg)`
const VersionRE = `^(\d+)?(?:\.(\d+)(?:\.(\d+))?)?(?:(\w)(?:(\d+))?)?$`

func ParseVersions(url string, version string) {
    downloadRe := regexp.MustCompile(downloadMatchRe)
    // versionRe := regexp.MustCompile(VersionRE)

    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
        return
    }
    defer response.Body.Close()

    contents, _ := ioutil.ReadAll(response.Body)
    matches := downloadRe.FindAllString(string(contents), -1)

    for _, m := range matches {
        fmt.Println(m)
    }
}
