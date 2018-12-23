package settings

import (
	"github.com/BurntSushi/toml"
	"os"
	"os/user"
	"path"
)

const settingsDir = ".unityLoader"
const configName = "config.toml"

type Settings struct {
	ProjectDirectory string `toml:"ProjectDirectory"`
}

func ParseDotFile() (*Settings, error) {
	dotPath, err := getPath()
	if err != nil {
		return nil, err
	}
	configPath := path.Join(dotPath, configName)

	f, err := os.Open(configPath)
	if os.IsNotExist(err) {
		if err := createDotFile(configPath); err != nil {
			return nil, err
		}
		return &Settings{}, nil
	} else if err != nil {
		return nil, err
	}
	defer f.Close()

	var data Settings

	if _, err := toml.DecodeReader(f, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func createDotFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data := Settings{}

	if err := toml.NewEncoder(f).Encode(data); err != nil {
		return err
	}
	return nil
}

func getPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, settingsDir), nil
}
