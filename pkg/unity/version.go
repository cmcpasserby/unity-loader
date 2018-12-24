package unity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var verTypeRe = regexp.MustCompile(`[pfba]`)

type VersionData struct {
	Major   int
	Minor   int
	Update  int
	VerType string
	Patch   int
}

func (v *VersionData) String() string {
	return fmt.Sprintf("%d.%d.%d%c%d", v.Major, v.Minor, v.Update, v.VerType, v.Patch)
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
