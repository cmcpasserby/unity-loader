package packages

import "testing"

func TestVersionLess(t *testing.T) {
    // test major
    a := VersionData{5, 6, 4, "f", 1}
    b := VersionData{2018, 2, 7, "f", 1}

    if !VersionLess(a, b) {
        t.Error("major")
    }

    // test minor
    a = VersionData{2018, 4, 4, "f", 1}
    b = VersionData{2018, 6, 7, "f", 1}

    if !VersionLess(a, b) {
        t.Error("minor")
    }

    // test update
    a = VersionData{2018, 6, 4, "f", 1}
    b = VersionData{2018, 6, 7, "f", 1}

    if !VersionLess(a, b) {
        t.Error("update")
    }

    // test type
    a = VersionData{2018, 6, 7, "f", 1}
    b = VersionData{2018, 6, 7, "p", 1}

    if !VersionLess(a, b) {
        t.Error("type")
    }

    // test patch
    a = VersionData{2018, 6, 7, "f", 1}
    b = VersionData{2018, 6, 7, "f", 2}

    if !VersionLess(a, b) {
        t.Error("patch")
    }
}
