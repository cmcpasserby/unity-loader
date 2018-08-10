package unity

import (
    "net/http"
    "io/ioutil"
    "gopkg.in/ini.v1"
    "fmt"
    "path/filepath"
)

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

type UrlData struct {
    Base string
    Version VersionData
}

func (url *UrlData) GetIniUrl() string {
    fileName := fmt.Sprintf(configName, url.Version.VersionString)
    return fmt.Sprintf(url.Base, url.Version.VersionUuid) + fileName;
}

type Package struct {
    Title string `ini:"title"`
    Description string `ini:"description"`
    Path string `ini:"url"`
    Install bool `ini:"install"`
    Size int64 `ini:"size"`
    InstalledSize int64 `ini:"installedsize"`
    Version string `ini:"version"`
    Md5 string `ini:"md5"`
    Hidden bool `ini:"hidden"`
    Extension string `ini:"extension"`
    Url UrlData
}

func (pkg *Package) GetDownloadUrl() string {
    base := fmt.Sprintf(pkg.Url.Base, pkg.Url.Version.VersionUuid)
    return filepath.Join(base, pkg.Path)
}

func getPackages(ver VersionData) (map[string]Package, error) {
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

    packages := make(map[string]Package)

    testIgnored := func(item string) bool {
        for _, name := range ignoredSections {
            if item == name {
                return true
            }
        }
        return false
    }

    for _, name := range cfg.SectionStrings() {
        if testIgnored(name) {
            continue
        }

        pkg := new(Package)
        cfg.Section(name).MapTo(pkg)
        pkg.Url = currentUrl
        packages[name] = *pkg
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
