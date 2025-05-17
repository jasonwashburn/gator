package config

import (
	"encoding/json"
	"os"
)

type ConfigFile struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return userHomeDir + "/" + configFileName, nil
}

func Read() (ConfigFile, error) {
	filepath, err := getConfigFilePath()
	if err != nil {
		return ConfigFile{}, err
	}
	file, err := os.Open(filepath)
	if err != nil {
		return ConfigFile{}, err
	}
	defer file.Close()

	var config ConfigFile
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return ConfigFile{}, err
	}
	return config, nil
}

func (c *ConfigFile) write() error {
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigFile) SetUser(user string) error {
	c.CurrentUserName = user

	err := c.write()
	if err != nil {
		return err
	}
	return nil
}
