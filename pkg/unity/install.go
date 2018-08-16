package unity

import (
    "io/ioutil"
    "fmt"
    "os"
    "time"
    "net/http"
    "path"
    "io"
    "os/exec"
    "errors"
    "log"
)

var tempDir string

func Install(version string) error {
    if os.Getuid() != 0 {
        return errors.New("admin is required to install packages, try running with sudo")
    }

    versionData, err := GetVersionData(version)
    if err != nil {return err}

    packages, err := getPackages(versionData)
    if err != nil {return err}

    defer cleanUp()

    pkgPath, err := download(packages["Unity"])
    if err != nil {
        return err
    }
    fmt.Println(pkgPath)

    err = installPkg(pkgPath)
    if err != nil {return err}

    return nil
}

func download(pkg *Package) (string, error) {
    pkgDirectory, err := ioutil.TempDir("", "unitypacakges_")
    if err != nil {return "", err}

    tempDir = pkgDirectory

    url := pkg.GetDownloadUrl()
    fileName := path.Base(url)
    filePath :=  path.Join(pkgDirectory, fileName)

    start := time.Now()

    out, err := os.Create(filePath)
    if err != nil {return "", err}
    defer out.Close()

    done := make(chan int64)

    go downloadProgress(done, filePath, pkg.Size)

    response, err := http.Get(pkg.GetDownloadUrl())
    if err != nil {return "", err}
    defer response.Body.Close()

    n, err := io.Copy(out, response.Body)
    if err != nil {return "", err}

    done <- n

    fmt.Printf("Download completed in %s\n", time.Since(start))
    return filePath, nil
}

func installPkg(filePath string) error {
    process := exec.Command("installer", "-package", filePath, "-target", "/")
    return process.Run()
}

func cleanUp() {
    if tempDir == "" {
        return
    }

    dirRead, _ := os.Open(tempDir)
    dirFiles, _ := dirRead.Readdir(0)

    for i := range dirFiles {
        f := dirFiles[i]
        path := path.Join(tempDir, f.Name())
        os.Remove(path)
    }

    os.Remove(tempDir)
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

            fmt.Printf("\rDownloading Unity Editor %.0f%%", percent)
        }
        if stop {
            fmt.Printf("\r100")
            fmt.Println("%")
            return
        }
        time.Sleep(time.Second)
    }
}

