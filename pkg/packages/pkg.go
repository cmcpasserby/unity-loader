package packages

import (
    "fmt"
    "io/ioutil"
    "path"
    "time"
    "os"
    "net/http"
    "io"
    "os/exec"
    "errors"
    "log"
    "crypto/md5"
    "encoding/hex"
)


type UrlData struct {
    Base string
    Version VersionData
}

func (url *UrlData) GetIniUrl() string {
    fileName := fmt.Sprintf(configName, url.Version.VersionString)
    return fmt.Sprintf(url.Base, url.Version.VersionUuid) + fileName
}


type PkgData struct {
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

type Package struct {
    Data PkgData
    Url UrlData
    filePath string
}

func (pkg *Package) GetDownloadUrl() string {
    base := fmt.Sprintf(pkg.Url.Base, pkg.Url.Version.VersionUuid)
    return base + pkg.Data.Path
}

func (pkg *Package) DownloadPkg() error {
    pkgDirectory, err := ioutil.TempDir("", "unitypackages_")
    if err != nil {return err}

    url := pkg.GetDownloadUrl()
    fileName := path.Base(url)
    filePath := path.Join(pkgDirectory, fileName)

    start := time.Now()

    out, err := os.Create(filePath)
    if err != nil {return err}
    defer out.Close()

    pkg.filePath = filePath

    done := make(chan int64)
    go pkg.downloadProgress(done)

    response, err := http.Get(url)
    if err != nil {return err}
    defer response.Body.Close()

    n, err := io.Copy(out, response.Body)
    if err != nil {return err}

    done <- n

    fmt.Printf("Download completed in %s\n", time.Since(start))
    return nil
}

func (pkg *Package) ValidatePkg() (bool, error) {
    if pkg.filePath == "" {
        return false, errors.New("no downloaded package to install")
    }

    file, err := os.Open(pkg.filePath)
    if err != nil {return false, err}
    defer file.Close()

    hash := md5.New()
    _, err = io.Copy(hash, file)
    if err != nil {return false, err}

    sum := hash.Sum(nil)
    isValid := hex.EncodeToString(sum) == pkg.Data.Md5

    return isValid, nil
}

func (pkg *Package) InstallPkg() error {
    if pkg.filePath == "" {
        return errors.New("no downloaded package to install")
    }

    process := exec.Command("installer", "-package", pkg.filePath, "-target", "/")
    err := process.Run()
    if err != nil {return err}

    os.Remove(pkg.filePath)
    pkg.filePath = ""

    return nil
}

func (pkg *Package) downloadProgress(done chan int64) {
    stop := false

    for {
        select {
        case <- done:
            stop = true
        default:
            file, err := os.Open(pkg.filePath)
            if err != nil {log.Fatal(err)}


            fi, err := file.Stat()
            if err != nil {log.Fatal(err)}

            size := fi.Size()

            if size == 0 {
                size = 1
            }

            percent := (float64(size) / float64(pkg.Data.Size)) * 100
            fmt.Printf("\rDownloading %s, %.0f%%", pkg.Data.Title, percent)
        }
        if stop {
            fmt.Printf("\r100")
            fmt.Println("%")
            return
        }
        time.Sleep(time.Second)
    }
}
