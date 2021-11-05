package unity

import (
	"errors"
	"fmt"
)

type versionNotFoundError struct {
	version string
}

func (err versionNotFoundError) Error() string {
	return fmt.Sprintf("unity version %s not found", err.version)
}

// IsVersionNotFound used for checking an error to confirm of it is a version not found error
func IsVersionNotFound(err error) bool {
	var notFound versionNotFoundError
	return errors.As(err, &notFound)
}
