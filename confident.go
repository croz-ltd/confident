package confident

import (
	"encoding/json"
	"fmt"
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

func (k *Confident) Read() {
	configFilePath := k.Path + "/" + k.Name + "." + k.Type
	_, err := os.Stat(configFilePath)
	if err != nil {
		// Config file does not exists, skip reading
		return
	}

	configFileBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("fatal: fail to read configuration file")
		os.Exit(1)
	}

	err = json.Unmarshal(configFileBytes, &k.Config)
	if err != nil {
		fmt.Println("fatal: fail to read configuration file to provided stucture")
		os.Exit(1)
	}

	k.InitialHash = CalculateHash(k.Config)
}

func (k *Confident) Persist() {
	hash := CalculateHash(k.Config)

	if hash != k.InitialHash {
		configPath := k.Path + "/" + k.Name + "." + k.Type

		b, err := json.MarshalIndent(k.Config, "", " ")
		if err != nil {
			fmt.Println("fatal: error persisting config change")
			os.Exit(1)
		}

		file, err := os.Create(configPath)
		if err != nil {
			fmt.Println("fatal: error persisting config change to file")
			os.Exit(1)
		}
		defer file.Close()

		err = file.Chmod(k.Permission)
		if err != nil {
			fmt.Println("fatal: fail to set wanted permissions to config file")
		}

		_, err = file.Write(b)
		if err != nil {
			fmt.Println("fatal: error to persisting config change to file")
			os.Exit(1)
		}
	}
}
