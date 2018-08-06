package unity

import (
    "net/http"
    "io/ioutil"
    "gopkg.in/ini.v1"
    "fmt"
)

var configUrls = [...]string {
    "https://netstorage.unity3d.com/unity/%s/unity-%s-osx.ini",
    "https://download.unity3d.com/download_unity/%s/unity-%s-osx.ini",
    "https://beta.unity3d.com/download/%s/unity-%s-osx.ini",
    "https://files.unity3d.com/bootstrapper/%s/unity-%s-osx.ini",
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
}

func getPackages(ver VersionData) (map[string]Package, error) {
    var response *http.Response
    var err error

    for _, url := range configUrls {
        response, err = http.Get(fmt.Sprintf(url, ver.VersionUuid, ver.VersionString))
        if err == nil {break}
    }
    defer response.Body.Close()

    contents, err := ioutil.ReadAll(response.Body)
    if err != nil {return nil, err}

    cfg, err := ini.Load(contents)
    if err != nil {return nil, err}

    packages := make(map[string]Package)
    for _, name := range cfg.SectionStrings() {
        pkg := new(Package)
        cfg.Section(name).MapTo(pkg)
        packages[name] = *pkg
    }
    return packages, nil
}
