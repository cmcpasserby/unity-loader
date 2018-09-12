package versions

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

var verTypeRe = regexp.MustCompile(`[pfba]`)

type Data struct {
    Major int
    Minor int
    Update int
    Type string
    Patch int
}

type ExtendedData struct {
    Data
    Uuid string
}

func (v *Data) String() string {
    return fmt.Sprintf("%d.%d.%d%s%d", v.Major, v.Minor, v.Update, v.Type, v.Patch)
}

func DataFromString(input string) Data {
    separated := strings.Split(input, ".")

    major, _ := strconv.Atoi(separated[0])
    minor, _ := strconv.Atoi(separated[1])

    final := verTypeRe.Split(separated[2], -1)

    update, _ := strconv.Atoi(final[0])
    verType := verTypeRe.FindString(separated[2])
    patch, _ := strconv.Atoi(final[1])

    return Data{major, minor, update, verType, patch}
}
