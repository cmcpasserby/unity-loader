package settings

import "fmt"

type CacheNotFoundError struct {
	path string
}

func (err CacheNotFoundError) Error() string {
	return fmt.Sprintf("package cache not found at %q", err.path)
}
