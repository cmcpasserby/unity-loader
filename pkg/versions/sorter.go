package versions

var releaseTypeSort = map[string]int{"p": 4, "f": 3, "b": 2, "a": 1}

type ByVersionSorter []Data

func (s ByVersionSorter) Len() int {
    return len(s)
}

func (s ByVersionSorter) Less(i, j int) bool {
    return VersionLess(s[i], s[j])
}

func (s ByVersionSorter) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

type ByExtendedVersionSorter []ExtendedData

func (s ByExtendedVersionSorter) Len() int {
    return len(s)
}

func (s ByExtendedVersionSorter) Less(i, j int) bool {
    return VersionLess(s[i].Data, s[j].Data)
}

func (s ByExtendedVersionSorter) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func VersionLess(a, b Data) bool {
    if a.Major != b.Major {
        return a.Major < b.Major
    }

    if a.Minor != b.Minor {
        return a.Minor < b.Minor
    }

    if a.Update != b.Update {
        return a.Update < b.Update
    }

    aType := releaseTypeSort[a.Type]
    bType := releaseTypeSort[b.Type]
    if aType != bType {
        return aType < bType
    }

    if a.Patch != b.Patch {
        return a.Patch < b.Patch
    }

    return false
}