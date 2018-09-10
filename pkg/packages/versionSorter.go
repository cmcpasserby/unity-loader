package packages

var releaseTypeSort = map[string]int{"p": 4, "f": 3, "b": 2, "a": 1}

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
    if a.Major < b.Major {
        return true
    }

    if a.Minor < b.Minor {
        return true
    }

    if a.Update < b.Update {
        return true
    }

    if releaseTypeSort[a.VerType] < releaseTypeSort[b.VerType] {
        return true
    }

    if a.Patch < b.Patch {
        return true
    }

    return false
}
