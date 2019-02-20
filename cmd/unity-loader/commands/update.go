package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/pkg/settings"
)

func update(args ...string) error {
	fmt.Println("Updating Package Cache")

	cache := new(settings.Cache)
	return cache.Update()
}
