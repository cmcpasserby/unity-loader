package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

const (
	configFolderName = "unity-loader"
	configFileName   = "config.toml"
)

//go:embed config_header.toml
var configHeader string

type config struct {
	UnityHubPath   string   `toml:"unity_hub_path"`
	SearchPaths    []string `toml:"search_paths"`
	DefaultModules []string `toml:"default_modules"`
}

func getConfig() (*config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(configDir, configFolderName, configFileName)

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return createConfig(path)
		}
		return nil, err
	}
	defer f.Close()

	var config config
	if _, err := toml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func createConfig(path string) (*config, error) {
	fmt.Printf("no config found, creating default config in \"%s\"\n", path)

	defaultConfig, err := getDefault()
	if err != nil {
		return nil, err
	}

	if err = os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return nil, err
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	_, _ = fmt.Fprintf(f, "%s\n", configHeader)

	err = toml.NewEncoder(f).Encode(*defaultConfig)
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}

func getDefault() (*config, error) {
	switch runtime.GOOS {
	case "darwin":
		return &config{
			UnityHubPath: "/Applications/Unity Hub.app",
			SearchPaths:  []string{"/Applications/Unity/Hub/Editor"},
		}, nil
	case "windows":
		return &config{
			UnityHubPath: "C:/Program Files/Unity Hub/Unity Hub.exe",
			SearchPaths:  []string{"C:/Program Files/Unity/Hub/Editor"},
		}, nil
	default:
		return nil, fmt.Errorf("no default config for GOOS: %s", runtime.GOOS)
	}
}
