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

type VersionData struct {
	Major   int
	Minor   int
	Update  int
	VerType string
	Patch   int
}

func (v *VersionData) String() string {
	return fmt.Sprintf("%d.%d.%d%s%d", v.Major, v.Minor, v.Update, v.VerType, v.Patch)
}

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

type ExtendedVersionData struct {
	VersionData
	RevisionHash string
}

// ByVersionSorter properly sorts versions numbers
type ByVersionSorter []VersionData

func (s ByVersionSorter) Len() int {
	return len(s)
}

func (s ByVersionSorter) Less(i, j int) bool {
	return s[i].Compare(s[j]) < 0
}

func (s ByVersionSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
