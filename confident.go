package confident

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Confident struct {
	Name        string
	Type        string
	Path        string
	Permission  os.FileMode
	Config      interface{}
	InitialHash [16]byte
}

func New() *Confident {
	k := new(Confident)
	k.Name = "config"
	k.Type = "json"
	k.Path = "."
	k.Permission = os.FileMode(0644)
	return k
}

func (k *Confident) WithConfiguration(config interface{}) {
	k.Config = config
	k.InitialHash = CalculateHash(config)
}

func (k *Confident) GetConfig() interface{} {
	return k.Config
}

func (k *Confident) Read() error {
	configFilePath := k.Path + "/" + k.Name + "." + k.Type
	_, err := os.Stat(configFilePath)
	if err != nil {
		// Config file does not exists, skip reading
		return nil
	}

	configFileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return ConfidentFileNotReadableError{Path: configFilePath}
	}

	switch k.Type {
	case "json":
		err = json.Unmarshal(configFileBytes, &k.Config)
		break
	case "yaml", "yml":
		err = yaml.Unmarshal(configFileBytes, &k.Config)
		break
	}

	if err != nil {
		return ConfidentUnmarshallingError{Path: configFilePath, UnmarshalError: err}
	}

	k.InitialHash = CalculateHash(k.Config)
	return nil
}

func (k *Confident) PersistConfiguration(force bool) error {
	hash := CalculateHash(k.Config)

	if (hash != k.InitialHash) || force {
		configPath := k.Path + "/" + k.Name + "." + k.Type

		var b []byte
		var err error

		fmt.Println(k.Type)
		switch k.Type {
		case "json":
			b, err = json.MarshalIndent(k.Config, "", " ")
			break
		case "yaml", "yml":
			b, err = yaml.Marshal(k.Config)
			break
		}
		if err != nil {
			return ConfidentMarshallingError{Path: configPath, MarshalError: err}
		}

		file, err := os.Create(configPath)
		if err != nil {
			return ConfidentFileCreationError{Path: configPath, CreationError: err}
		}
		defer file.Close()

		err = file.Chmod(k.Permission)
		if err != nil {
			fmt.Println("warning: fail to set wanted permissions to config file")
		}

		_, err = file.Write(b)
		if err != nil {
			return ConfidentWriteError{Path: configPath, WriteError: err}
		}
	}

	return nil
}

func (k *Confident) Persist() error {
	return k.PersistConfiguration(false)
}
