package settings

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

const settingsDir = ".unityLoader"
const configName = "config.toml"
const packages = "packages"

type Settings struct {
	ProjectDirectory string `toml:"ProjectDirectory"`
}

func ParseDotFile() (*Settings, error) {
	dotPath, err := GetPath()
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

	defer closeFile(f)

	var data Settings

	if _, err := toml.DecodeReader(f, &data); err != nil {
		return nil, err
	}

	if strings.HasPrefix(data.ProjectDirectory, "~") {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		data.ProjectDirectory = filepath.Join(usr.HomeDir, data.ProjectDirectory[1:])
	}

	return &data, nil
}

func createDotFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer closeFile(f)

	data := Settings{}

	if err := toml.NewEncoder(f).Encode(data); err != nil {
		return err
	}
	return nil
}

func GetPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, settingsDir), nil
}

func GetPkgPath() (string, error) {
	dotPath, err := GetPath()
	if err != nil {
		return "", nil
	}
	return path.Join(dotPath, packages), nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
