package packages

import (
    "strconv"
    "strings"
)

var releaseLetterSort = map[string]int{"p": 4, "f": 3, "b": 2, "a": 1}

type ByVersionSorter []VersionData

func (s ByVersionSorter) Len() int {
    return len(s)
}

func (s ByVersionSorter) Less(i, j int) bool {
}

func (s ByVersionSorter) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func splitVersionString(ver string) (int, int, int, string, int) {
    seperated := strings.Split(ver, ".")
    main, _ := strconv.Atoi(seperated[0])
    major, _ := strconv.Atoi(seperated[1])
}
