package unity

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const archiveURL = "https://unity3d.com/get-unity/download/archive"

var hubURLRe = regexp.MustCompile(`unityhub://(\d+\.\d+\.\d+[pfba]\d+)/([0-9a-f]{12})`)

func GetAllVersions() ([]VersionData, error) {
	resp, err := http.Get(archiveURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	matches := hubURLRe.FindAllSubmatch(body, -1)
	versions := make([]VersionData, len(matches))

	for i, m := range matches {
		versions[i], err = VersionFromString(string(m[1]))
		if err != nil {
			return nil, err
		}
		versions[i].RevisionHash = string(m[2])
	}

	return versions, nil
}

func SearchArchive(partialVersion string) ([]VersionData, error) {
	allVersions, err := GetAllVersions()
	if err != nil {
		return nil, err
	}

	results := make([]VersionData, 0)
	for _, ver := range allVersions {
		if strings.HasPrefix(ver.String(), partialVersion) {
			results = append(results, ver)
		}
	}

	return results, nil
}

func InstallFromArchive(ver VersionData, hubPath string, modules, searchPaths []string) (InstallInfo, error) {
	hubPath, err := binFromApp(hubPath)
	if err != nil {
		return InstallInfo{}, err
	}

	args := []string{"--", "--headless", "-v", ver.String(), "--changeset", ver.RevisionHash}
	for _, mod := range modules {
		args = append(args, "-m", mod)
	}

	cmd := exec.Command(hubPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return InstallInfo{}, err
	}

	return GetInstallFromVersion(ver, searchPaths...)
}
