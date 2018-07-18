package unity

import (
    "os"
    "bufio"
    "strings"
    "errors"
    "path/filepath"
    "log"
    "howett.net/plist"
    "fmt"
)

type appInfoDict struct {
    CFBundleVersion string `plist:"CFBundleVersion"`
}

func GetUnityVersion(versionFile string) (string, error) {
    file, _ := os.Open(versionFile)
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)

    for scanner.Scan() {
        text := scanner.Text()
        if strings.HasPrefix(text, "m_EditorVersion:") {
            return strings.TrimSpace(strings.Split(text, ":")[1]), nil
        }
    }
    return "", errors.New("invalid ProjectVersion.txt")
}

func GetExecutable(version string) (string, error) {
    unityPaths, err := filepath.Glob("/Applications/**/Unity.app")
    if err != nil {
        log.Fatal(err)
    }

    var execPath string

    for _, path := range unityPaths {
        plistPath := filepath.Join(path, "Contents/info.plist")
        file, _ := os.Open(plistPath)

        var appInfo appInfoDict

        decoder := plist.NewDecoder(file)
        err := decoder.Decode(&appInfo)
        if err != nil {
            log.Fatal(err)
        }

        if appInfo.CFBundleVersion == version {
            execPath = path
            file.Close()
            break
        }
        file.Close()
    }

    if execPath != "" {
        return execPath, nil
    }

    return "", errors.New(fmt.Sprintf("unity version %s not found", version))
}
