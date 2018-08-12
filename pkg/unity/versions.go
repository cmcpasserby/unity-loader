package unity

import (
    "io/ioutil"
    "regexp"
    "fmt"
    "os"
    "log"
    "time"
    "net/http"
    "path"
    "io"
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

func Install(version string) error {
    versionData, err := GetVersionData(version)
    if err != nil {return err}

    packages, err := getPackages(versionData)
    if err != nil {return err}

    download(packages["Unity"])

    return nil
}

func download(pkg *Package) error {
    pkgDirectory, err := ioutil.TempDir("", "unitypacakges_")
    if err != nil {return err}

    url := pkg.GetDownloadUrl()
    fileName := path.Base(url)
    filePath :=  path.Join(pkgDirectory, fileName)

    start := time.Now()

    out, err := os.Create(filePath)
    if err != nil {return err}
    defer out.Close()

    done := make(chan int64)

    go downloadProgress(done, filePath, pkg.Size)

    response, err := http.Get(pkg.GetDownloadUrl())
    if err != nil {return err}
    defer response.Body.Close()

    n, err := io.Copy(out, response.Body)
    if err != nil {return err}

    done <- n

    fmt.Printf("Download completed in %s", time.Since(start))
    return nil
}

func downloadProgress(done chan int64, path string, total int64) {
    stop := false

    for {
        select {
        case <- done:
            stop = true
        default:
            file, err := os.Open(path)
            if err != nil {log.Fatal(err)}

            fi, err := file.Stat()
            if err != nil {log.Fatal(err)}

            size := fi.Size()

            if size == 0 {
                size = 1
            }

            percent := float64(size) / float64(total) * 100

            fmt.Printf("%.0f", percent)
            fmt.Println("%")
        }
        if stop {
            fmt.Printf("100")
            fmt.Println("%")
            return
        }
        time.Sleep(time.Second)
    }
}
