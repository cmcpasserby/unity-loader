package packages

import (
    "net/http"
    "io/ioutil"
    "gopkg.in/ini.v1"
    "fmt"
    "regexp"
)

var downloadRe = regexp.MustCompile(`(https?://[\w/.-]+/[0-9a-f]{12}/)[\w/.-]+-(\d+\.\d+\.\d+\w\d+)(?:\.dmg|\.pkg)`)
var versionRe = regexp.MustCompile(`(\d+\.\d+\.\d+\w\d+)`)
var uuidRe = regexp.MustCompile(`[0-9a-f]{12}`)

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
var ignoredSections = [...]string {
    "DEFAULT",
    "VisualStudio",
    "Mono",
}

const configName = "unity-%s-osx.ini"

var baseUrls = [...]string {
    "https://netstorage.unity3d.com/unity/%s/",
    "https://download.unity3d.com/download_unity/%s/",
    "https://beta.unity3d.com/download/%s/",
    "https://files.unity3d.com/bootstrapper/%s/",
}

func getPackages(ver VersionData) ([]*Package, error) {
    var response *http.Response
    var err error
    var currentUrl UrlData

    for _, url := range buildConfigUrls(ver) {
        response, err = http.Get(url.GetIniUrl())
        if err == nil {
            currentUrl = url
            break
        }
    }
    defer response.Body.Close()

    contents, err := ioutil.ReadAll(response.Body)
    if err != nil {return nil, err}

    cfg, err := ini.Load(contents)
    if err != nil {return nil, err}

    testIgnored := func(item string) bool {
        for _, name := range ignoredSections {
            if item == name {return true}
        }
        return false
    }

    sectionNames := cfg.SectionStrings()
    packages := make([]*Package, 0, len(sectionNames))

    for _, name := range sectionNames {
        if testIgnored(name) {continue}

        pkg := new(Package)

        cfg.Section(name).MapTo(&pkg.Data)
        pkg.Url = currentUrl

        packages = append(packages, pkg)
    }
    return packages, nil
}

func buildConfigUrls(ver VersionData) []UrlData {
    urls := make([]UrlData, 0, len(baseUrls))
    for _, baseUrl := range baseUrls {
        urls = append(urls, UrlData{Base:baseUrl, Version: ver})
    }
    return urls
}

func getVersionsFromUrl(url string, ver string, ch chan<- *VersionData) {
    response, err := http.Get(url)
    if err != nil {
        ch <- nil
        return
    }
    defer response.Body.Close()

    contents, _ := ioutil.ReadAll(response.Body)
    matches := downloadRe.FindAllString(string(contents), -1)

    for _, m := range matches {
        verStr := versionRe.FindString(m)
        if verStr == ver {
            verUuid := uuidRe.FindString(m)
            ch <- &VersionData{verStr, verUuid}
            return
        }
    }
    ch <- nil
}

func GetVersionData(ver string) (VersionData, error) {
    if !versionRe.MatchString(ver) {
        return VersionData{}, fmt.Errorf("unity version %q is not a valid unity version", ver)
    }

    ch := make(chan *VersionData)

    for _, url := range archiveUrls {
        go getVersionsFromUrl(url, ver, ch)
    }

    for res := range ch {
        if res != nil {
            return *res, nil
        }
    }

    return VersionData{}, fmt.Errorf("unity Version %q not found", ver)
}
