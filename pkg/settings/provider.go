package settings

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

const (
	settingsDir = ".unityLoader"
	configName = "config.toml"
	packages = "packages"
)


const configHelpString = `# config.toml
#
# help text
`

type Settings struct {
	ProjectDirectory string          `toml:"ProjectDirectory"`
	ModuleDefaults   map[string]bool `toml:"ModuleDefaults"`
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

	data := Settings{
		ProjectDirectory: "",
		ModuleDefaults:   map[string]bool{},
	}

	b := new(bytes.Buffer)

	if err := toml.NewEncoder(b).Encode(data); err != nil {
		return err
	}

	if _, err := f.WriteString(fmt.Sprintf("%s\n", configHelpString)); err != nil {
		return err
	}

	if _, err := f.Write(b.Bytes()); err != nil {
		return err
	}

	return nil
}

func GetPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	settingsPath := filepath.Join(usr.HomeDir, settingsDir)

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		if err := os.Mkdir(settingsPath, 0755); err != nil {
			return "", err
		}
	}

	return settingsPath, nil
}

func GetPkgPath() (string, error) {
	dotPath, err := GetPath()
	if err != nil {
		return "", nil
	}

	pkgPath := filepath.Join(dotPath, packages)

	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		if err := os.Mkdir(pkgPath, 0755); err != nil {
			return "", err
		}
	}

	return pkgPath, nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
