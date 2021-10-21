package unity

import "fmt"

type versionNotFoundError struct {
	version string
}

func (err versionNotFoundError) Error() string {
	return fmt.Sprintf("unity version %s not found", err.version)
}

func IsVersionNotFound(err error) bool {
	_, ok := err.(versionNotFoundError)
	return ok
}
