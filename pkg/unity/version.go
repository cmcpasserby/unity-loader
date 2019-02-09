package unity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	verTypeRe = regexp.MustCompile(`[pfba]`)
	releaseTypeSort = map[string]int{"p": 4, "f": 3, "b": 2, "a": 1}
)

type VersionData struct {
	Major   int
	Minor   int
	Update  int
	VerType string
	Patch   int
}

type ExtendedVersionData struct {
	VersionData
	RevisionHash string
}

func (v *VersionData) String() string {
	return fmt.Sprintf("%d.%d.%d%s%d", v.Major, v.Minor, v.Update, v.VerType, v.Patch)
}

func VersionDataFromString(input string) VersionData {
	separated := strings.Split(input, ".")

	major, _ := strconv.Atoi(separated[0])
	minor, _ := strconv.Atoi(separated[1])

	final := verTypeRe.Split(separated[2], -1)

	update, _ := strconv.Atoi(final[0])
	verType := verTypeRe.FindString(separated[2])
	patch, _ := strconv.Atoi(final[1])

	return VersionData{major, minor, update, verType, patch}
}

// Version Sorting
type ByVersionSorter []VersionData

func (s ByVersionSorter) Len() int {
	return len(s)
}

func (s ByVersionSorter) Less(i, j int) bool {
	return VersionLess(s[i], s[j])
}

func (s ByVersionSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func VersionLess(a, b VersionData) bool {
	if a.Major != b.Major {
		return a.Major < b.Major
	}

	if a.Minor != b.Minor {
		return a.Minor < b.Minor
	}

	if a.Update != b.Update {
		return a.Update < b.Update
	}

	aType := releaseTypeSort[a.VerType]
	bType := releaseTypeSort[b.VerType]
	if aType != bType {
		return aType < bType
	}

	if a.Patch != b.Patch {
		return a.Patch < b.Patch
	}

	return false
}
