package versions

import "fmt"

type VersionNotFoundError struct {
    version string
}

func (err VersionNotFoundError) Error() string {
    return fmt.Sprintf("unity version %q not found", err.version)
}

type InvalidVersionError struct {
    version string
}

func (err InvalidVersionError) Error() string {
    return fmt.Sprintf("unity version %q is not a valid unity version", err.version)
}
