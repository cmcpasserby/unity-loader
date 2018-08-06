package unity

import (
    "net/http"
    "io/ioutil"
    "regexp"
    "fmt"
)

const (
    downloadMatchRe = `(https?://[\w/.-]+/[0-9a-f]{12}/)[\w/.-]+-(\d+\.\d+\.\d+\w\d+)(?:\.dmg|.pkg)`
    versionMatchRe = `(\d+\.\d+\.\d+\w\d+)`
    uuidMatchRe = `[0-9a-f]{12}`
)

var archiveUrls = [...]string {
    "https://unity3d.com/get-unity/download/archive",
    "https://unity3d.com/unity/qa/lts-releases",
    "https://unity3d.com/unity/qa/patch-releases",
    "https://unity3d.com/unity/beta-download",
}


type VersionData struct {
    VersionString string
    VersionUuid string
}

func GetVersionData(ver string) (VersionData, error) {
    downloadRe := regexp.MustCompile(downloadMatchRe)
    versionRe := regexp.MustCompile(versionMatchRe)
    uuidRe := regexp.MustCompile(uuidMatchRe)

    if !versionRe.MatchString(ver) {
        return VersionData{}, fmt.Errorf("unity version %q is not a valid unity version", ver)
    }

    for _, url := range archiveUrls {
        response, err := http.Get(url)
        if err != nil {return VersionData{}, err}

        contents, _ := ioutil.ReadAll(response.Body)
        matches := downloadRe.FindAllString(string(contents), -1)

        for _, m := range matches {
            verStr := versionRe.FindString(m)
            if verStr == ver {
                verUuid := uuidRe.FindString(m)
                return VersionData{verStr, verUuid}, nil
            }
        }

        response.Body.Close()
    }
    return VersionData{}, fmt.Errorf("unity Version %q not found", ver)
}

func Install(version string) error {
    versionData, err := GetVersionData(version)
    if err != nil {return err}

    packages, err := getPackages(versionData)
    if err != nil {return err}

    fmt.Println(packages["Unity"].Title)
    return nil
}
