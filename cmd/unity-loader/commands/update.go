package commands

import "github.com/cmcpasserby/unity-loader/pkg/settings"

func update(args ...string) error {
	cache := new(settings.Cache)
	return cache.Update()
}
