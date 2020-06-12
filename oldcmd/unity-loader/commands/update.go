package commands

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/settings"
)

func update(args ...string) error {
	fmt.Printf("Updating Package Cache...")
	cache := new(settings.Cache)

	if err := cache.Update(); err != nil {
		return err
	}

	fmt.Print("\033[2K") // clears current line
	fmt.Printf("\rPacakge Cache Updated, found %v Unity versions\n", cache.Releases.Len())

	return nil
}
