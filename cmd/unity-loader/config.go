package main

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

const (
	configName    = ".unity-loader"
	configComment = `# unity-loader config
#
# Multiple search paths can be defined, will search unity hubs default path if not defined`
)

type config struct {
	UnityHubPath string   `toml:"unity_hub_path"`
	SearchPaths  []string `toml:"search_paths"`
}

func getConfig() (*config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(usr.HomeDir, configName)

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return createConfig(path)
		} else {
			return nil, err
		}
	}
	defer f.Close()

	var config config
	if _, err := toml.DecodeReader(f, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func createConfig(path string) (*config, error) {
	defaultConfig, err := getDefault()
	if err != nil {
		return nil, err
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	_, _ = fmt.Fprintln(f, configComment)
	_, _ = fmt.Fprintln(f)

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
