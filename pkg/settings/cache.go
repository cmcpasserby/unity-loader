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
	Timestamp time.Time        `json:"timestamp"`
	Releases  parsing.Releases `json:"releases"`
}

func WriteCache(data *parsing.Releases) error {
	cache := Cache{
		time.Now(),
		*data,
	}

	cachePath, err := getCachePath()
	if err != nil {
		return err
	}

	f, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer closeFile(f)

	if err := json.NewEncoder(f).Encode(cache); err != nil {
		return err
	}

	return err
}

func ReadCache() (*Cache, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return nil, err
	}

	// TODO does cache file exist
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
