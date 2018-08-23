package unity

import "fmt"

type VersionNotFound struct {
    version string
}

func (err VersionNotFound) Error() string {
    return fmt.Sprintf("unity version %s not found", err.version)
}
