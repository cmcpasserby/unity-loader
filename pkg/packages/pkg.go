package packages

import (
    "crypto/md5"
    "encoding/hex"
    "errors"
    "fmt"
    "github.com/cmcpasserby/unity-loader/pkg/sudoer"
    "gopkg.in/cheggaaa/pb.v1"
    "io"
    "log"
    "net/http"
    "os"
    "path"
    "time"
)


type UrlData struct {
    Base string
    Version ExtendedVersionData
}

func (url *UrlData) GetIniUrl() string {
    fileName := fmt.Sprintf(configName, url.Version.String())
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
    RequiresUnity bool `ini:"requires_unity"`
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

func (pkg *Package) Download(tempDir string) error {
    url := pkg.GetDownloadUrl()
    fileName := path.Base(url)
    filePath := path.Join(tempDir, fileName)

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

    return nil
}

func (pkg *Package) Validate() (bool, error) {
    if pkg.filePath == "" {
        return false, errors.New("no downloaded package to validate")
    }

    fmt.Printf("Validating pacakge %q...", pkg.Data.Title)
    file, err := os.Open(pkg.filePath)
    if err != nil {return false, err}
    defer file.Close()

    hash := md5.New()
    _, err = io.Copy(hash, file)
    if err != nil {return false, err}

    sum := hash.Sum(nil)
    isValid := hex.EncodeToString(sum) == pkg.Data.Md5

    fmt.Print("\033[2K") // clears current line
    if isValid {
        fmt.Printf("\rPackage %q is valid\n", pkg.Data.Title)
    } else {
        fmt.Printf("\rPackage %q is not valid\n", pkg.Data.Title)
    }

    return isValid, nil
}

func (pkg *Package) Install(sudo *sudoer.Sudoer) error {
    if pkg.filePath == "" {
        return errors.New("no downloaded package to install")
    }

    fmt.Printf("Installing pacakge %q...", pkg.Data.Title)

    err := sudo.RunAsRoot("installer", "-package", pkg.filePath, "-target", "/")
    if err != nil {return err}

    os.Remove(pkg.filePath)
    pkg.filePath = ""

    fmt.Print("\033[2K") // clears current line
    fmt.Printf("\rInstalled pacakge %q\n", pkg.Data.Title)
    return nil
}

func (pkg Package) downloadProgress(done chan int64) {
    stop := false

    bar := pb.New64(pkg.Data.Size)
    bar.Prefix(pkg.Data.Title)
    bar.ShowSpeed = true
    bar.Width = 120
    bar.SetUnits(pb.U_BYTES)
    bar.Start()

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

            bar.Set64(size)
        }
        if stop {
            bar.Set64(pkg.Data.Size)
            bar.FinishPrint(fmt.Sprintf("Downloaded %q", pkg.Data.Title))
            return
        }
        time.Sleep(time.Second)
    }
}

func Filter(pkgs []*Package, f func(*Package) bool) []*Package {
    newPkgs := make([]*Package, 0)
    for _, pkg := range pkgs {
        if f(pkg) {
            newPkgs = append(newPkgs, pkg)
        }
    }
    return newPkgs
}
