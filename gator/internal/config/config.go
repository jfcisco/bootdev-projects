package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

func Read() Config {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		panic(err)
	}

	configFile, err := os.Open(configFilePath)
	if errors.Is(err, os.ErrNotExist) {
		panic("Config file not found. Please create a .gatorconfig.json file in your home directory.")
	} else if err != nil {
		panic(err)
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		panic(errors.Join(errors.New("failed to parse config file: "), err))
	}
	return config
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	configFile, err := os.Create(configFilePath)
	if err != nil {
		return err
	}

	_, err = configFile.Write(data)
	if err != nil {
		return err
	}

	return nil
}
