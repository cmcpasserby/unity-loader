package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"os/user"
	"path/filepath"
)

const configName = ".unity-loader"

type Config struct {
	SearchPaths []string `toml:"search_paths"`
}

func GetConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(usr.HomeDir, configName)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, configContents, 0644); err != nil {
			return nil, err
		}
		fmt.Printf("writing default config to \"%s\"\n", path)
	} else if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	if _, err := toml.DecodeReader(f, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
