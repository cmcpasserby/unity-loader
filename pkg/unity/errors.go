package unity

import "fmt"

type VersionNotFoundError struct {
    version string
}

func (err VersionNotFoundError) Error() string {
    return fmt.Sprintf("unity version %s not found", err.version)
}
