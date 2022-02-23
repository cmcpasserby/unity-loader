package unity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	verTypeRe       = regexp.MustCompile(`[pfba]`)
	releaseTypeSort = map[string]int{"p": 4, "f": 3, "b": 2, "a": 1}
)

// VersionData represents a Unity version in a comparable format
type VersionData struct {
	Major   int
	Minor   int
	Update  int
	VerType string
	Patch   int
}

// String outputs version in string format(major.minor.update verType patch)
func (v *VersionData) String() string {
	return fmt.Sprintf("%d.%d.%d%s%d", v.Major, v.Minor, v.Update, v.VerType, v.Patch)
}

// Compare comparison function for versions
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
		return v.Patch - v.Patch
	}

	return 0
}

// VersionFromString parses a string and returns a VersionData
func VersionFromString(input string) VersionData {
	separated := strings.Split(input, ".")

	major, _ := strconv.Atoi(separated[0])
	minor, _ := strconv.Atoi(separated[1])

	final := verTypeRe.Split(separated[2], -1)

	update, _ := strconv.Atoi(final[0])
	verType := verTypeRe.FindString(separated[2])
	patch, _ := strconv.Atoi(final[1])

	return VersionData{major, minor, update, verType, patch}
}

// ExtendedVersionData extends VersionData with a RevisionHash
type ExtendedVersionData struct {
	VersionData
	RevisionHash string
}
