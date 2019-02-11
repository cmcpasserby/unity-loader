package settings

import (
	"encoding/json"
	"github.com/cmcpasserby/unity-loader/pkg/parsing"
	"os"
	"path"
	"time"
)

const fileName = "cache.json"

type Cache struct {
	Timestamp time.Time `json:"timestamp"`
	Releases parsing.CacheVersionSlice `json:"releases"`
}

func (c *Cache) NeedsUpdate() bool {
	return time.Now().After(c.Timestamp.Add(time.Hour * 24))
}

func (c *Cache) Update() error {
	versions, err := parsing.GetVersions()
	if err != nil {
		return err
	}

	c.Timestamp = time.Now()
	c.Releases = versions

	cachePath, err := getCachePath()
	if err != nil {
		return err
	}

	f, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer closeFile(f)

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	if err := enc.Encode(c); err != nil {
		return err
	}

	return nil
}

func ReadCache() (*Cache, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		cache := new(Cache)
		if err := cache.Update(); err != nil {
			return nil, err
		}
		return cache, nil
	}

	f, err := os.Open(cachePath)
	if err != nil {
		return nil, err
	}

	var cache Cache

	if err := json.NewDecoder(f).Decode(&cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

func getCachePath() (string, error) {
	var cachePath string
	if tempPath, err := GetPath(); err != nil {
		return "", err
	} else {
		cachePath = path.Join(tempPath, fileName)
	}
	return cachePath, nil
}
