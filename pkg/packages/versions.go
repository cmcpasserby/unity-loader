package packages

import (
    "fmt"
    "gopkg.in/ini.v1"
    "io/ioutil"
    "net/http"
    "regexp"
    "strconv"
    "strings"
)

var downloadRe = regexp.MustCompile(`(https?://[\w/.-]+/[0-9a-f]{12}/)[\w/.-]+-(\d+\.\d+\.\d+\w\d+)(?:\.dmg|\.pkg)`)
var versionRe = regexp.MustCompile(`(\d+\.\d+\.\d+\w\d+)`)
var uuidRe = regexp.MustCompile(`[0-9a-f]{12}`)

var verTypeRe = regexp.MustCompile(`[pfba]`)

var archiveUrls = [...]string {
    "https://unity3d.com/get-unity/download/archive",
    "https://unity3d.com/unity/qa/lts-releases",
    "https://unity3d.com/unity/qa/patch-releases",
}

const betaUrl = "https://unity3d.com/unity/beta-download"

type VersionData struct {
    Major int
    Minor int
    Update int
    VerType string
    Patch int
}

type ExtendedVersionData struct {
    VersionData
    VersionUuid string
}

func (v *VersionData) String() string {
    return fmt.Sprintf("%d.%d.%d%s%d", v.Major, v.Minor, v.Update, v.VerType, v.Patch)
}

func VersionDataFromString(input string) VersionData {
    separated := strings.Split(input, ".")

    major, _ := strconv.Atoi(separated[0])
    minor, _ := strconv.Atoi(separated[1])

    final := verTypeRe.Split(separated[2], -1)

    update, _ := strconv.Atoi(final[0])
    verType := verTypeRe.FindString(separated[2])
    patch, _ := strconv.Atoi(final[1])

    return VersionData{major, minor, update, verType, patch}
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

func GetPackages(ver ExtendedVersionData) ([]*Package, error) {
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

func buildConfigUrls(ver ExtendedVersionData) []UrlData {
    urls := make([]UrlData, 0, len(baseUrls))
    for _, baseUrl := range baseUrls {
        urls = append(urls, UrlData{Base:baseUrl, Version: ver})
    }
    return urls
}

func getVersionsFromUrl(url string, ver string, ch chan<- *ExtendedVersionData) {
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
            verData := ExtendedVersionData{VersionDataFromString(verStr), verUuid}
            ch <- &verData
            return
        }
    }
    ch <- nil
}

func GetVersionData(ver string) (ExtendedVersionData, error) {
    if !versionRe.MatchString(ver) {
        return ExtendedVersionData{}, InvalidVersionError{ver}
    }

    ch := make(chan *ExtendedVersionData)

    for _, url := range archiveUrls {
        go getVersionsFromUrl(url, ver, ch)
    }

    resultCount := 0
    for res := range ch {
        resultCount += 1
        if res != nil {
            return *res, nil
        }

        if resultCount >= len(archiveUrls) {
            break
        }
    }
    return ExtendedVersionData{}, VersionNotFoundError{ver}
}
