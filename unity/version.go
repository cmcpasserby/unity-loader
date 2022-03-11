package unity

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	verTypeRe               = regexp.MustCompile(`[pfba]`)
	editorVersionRe         = regexp.MustCompile(`(\d+\.\d+\.\d+[pfba]\d+)`)
	editorVersionRevisionRe = regexp.MustCompile(`(\d+\.\d+\.\d+[pfba]\d+) \(([0-9a-f]{12})\)`)
	releaseTypeSort         = map[string]int{"p": 4, "f": 3, "b": 2, "a": 1}
)

// VersionData represents a Unity version in a comparable format
type VersionData struct {
	Major        int
	Minor        int
	Update       int
	VerType      string
	Patch        int
	RevisionHash string
}

func (v VersionData) HasRevisionHash() bool {
	return v.RevisionHash != ""
}

// String outputs version in string format "major.minor.update verType patch"
func (v VersionData) String() string {
	return fmt.Sprintf("%d.%d.%d%s%d", v.Major, v.Minor, v.Update, v.VerType, v.Patch)
}

// Compare comparison function for versions, ignores RevisionHash
func (v VersionData) Compare(other VersionData) int {
	if v.Major != other.Major {
		return v.Major - other.Major
	}

	if v.Minor != other.Minor {
		return v.Minor - other.Minor
	}

	if v.Update != other.Update {
		return v.Update - other.Update
	}

	aType := releaseTypeSort[v.VerType]
	bType := releaseTypeSort[other.VerType]
	if aType != bType {
		return aType - bType
	}

	if v.Patch != other.Patch {
		return v.Patch - other.Patch
	}

	return 0
}

func convertSegments(parts ...string) ([]int, error) {
	result := make([]int, len(parts))
	for i, str := range parts {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		result[i] = num
	}
	return result, nil
}

// VersionFromString parses a string and returns a VersionData
func VersionFromString(input string) (VersionData, error) {
	separated := strings.Split(input, ".")

	majorMinor, err := convertSegments(separated[:2]...)
	if err != nil {
		return VersionData{}, fmt.Errorf("invalid version number format: %w", err)
	}

	final := verTypeRe.Split(separated[2], -1)
	verType := verTypeRe.FindString(separated[2])

	updatePatch, err := convertSegments(final...)
	if err != nil {
		return VersionData{}, fmt.Errorf("invalid version number format: %w", err)
	}

	return VersionData{Major: majorMinor[0], Minor: majorMinor[1], Update: updatePatch[0], VerType: verType, Patch: updatePatch[1]}, nil
}

func readProjectVersion(reader io.Reader) (VersionData, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	versionData := VersionData{}
	var err error
	found := false

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "m_EditorVersion:") {
			versionStr := editorVersionRe.FindString(text)
			versionData, err = VersionFromString(versionStr)
			if err != nil {
				return VersionData{}, err
			}
			found = true
		} else if strings.HasPrefix(text, "m_EditorVersionWithRevision:") {
			groups := editorVersionRevisionRe.FindStringSubmatch(text)
			versionData, err = VersionFromString(groups[1])
			if err != nil {
				return VersionData{}, err
			}
			versionData.RevisionHash = groups[2]
			found = true
		}
	}

	if !found {
		return VersionData{}, errors.New("invalid ProjectVersion.txt")
	}
	return versionData, nil
}
