package versions

import "testing"

func TestVersionLess(t *testing.T) {
	// test major
	a := Data{5, 6, 4, "f", 1}
	b := Data{2018, 2, 7, "f", 1}

	if !VersionLess(a, b) {
		t.Error("major")
	}

	// test minor
	a = Data{2018, 4, 4, "f", 1}
	b = Data{2018, 6, 7, "f", 1}

	if !VersionLess(a, b) {
		t.Error("minor")
	}

	// test update
	a = Data{2018, 6, 4, "f", 1}
	b = Data{2018, 6, 7, "f", 1}

	if !VersionLess(a, b) {
		t.Error("update")
	}

	// test type
	a = Data{2018, 6, 7, "f", 1}
	b = Data{2018, 6, 7, "p", 1}

	if !VersionLess(a, b) {
		t.Error("type")
	}

	// test patch
	a = Data{2018, 6, 7, "f", 1}
	b = Data{2018, 6, 7, "f", 2}

	if !VersionLess(a, b) {
		t.Error("patch")
	}
}
