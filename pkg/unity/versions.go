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
    downloadMatchRe = `(https?://[\w/.-]+/[0-9a-f]{12}/)[\w/.-]+-(\d+\.\d+\.\d+\w\d+)(?:\.pkg)`
    versionMatchRe = `(\d+\.\d+\.\d+\w\d+)`
    uuidMatchRe = `[0-9a-f]{12}`
)

type VersionData struct {
    VersionString string
    VersionUuid string
}

func (v VersionData) GetEditorUrl() string {
    return fmt.Sprintf("http://netstorage.unity3d.com/unity/%s/MacEditorInstaller/Unity.pkg", v.VersionUuid)
}

func (v VersionData) GetAndroidSupportUrl() string {
    return fmt.Sprintf(
        "http://netstorage.unity3d.com/unity/%s/MacEditorTargetInstaller/UnitySetup-Android-Support-for-Editor-%s.pkg",
        v.VersionUuid,
        v.VersionString)
}

func (v VersionData) GetIosSupportUrl() string {
    return fmt.Sprintf(
        "http://netstorage.unity3d.com/unity/%s/MacEditorTargetInstaller/UnitySetup-iOS-Support-for-Editor-%s.pkg",
        v.VersionUuid,
        v.VersionString)
}

func (v VersionData) GetWebGlSupportUrl() string {
    return fmt.Sprintf(
        "http://netstorage.unity3d.com/unity/%s/MacEditorTargetInstaller/UnitySetup-WebGL-Support-for-Editor-%s.pkg",
        v.VersionUuid,
        v.VersionString)
}

func (v VersionData) GetWindowsSupportUrl() string {
    return fmt.Sprintf(
        "http://netstorage.unity3d.com/unity/%s/MacEditorTargetInstaller/UnitySetup-Windows-Support-for-Editor-%s.pkg",
        v.VersionUuid,
        v.VersionString)
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
            versions[ver] = VersionData{
                VersionString: ver,
                VersionUuid:uuidRe.FindString(url),
            }
        }
    }

    return versions, nil
}
