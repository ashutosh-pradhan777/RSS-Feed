package config

import (
	"encoding/json"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}


func getConfigFilePath() (string, error) {

	path, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	path += configFileName

	return path, nil

}

func Read() (Config, error) {

	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var outConfig Config

	err2 := json.Unmarshal(data, &outConfig)
	if err2 != nil {
		return Config{}, err2
	}

	return outConfig, nil

}

func (cnf Config) SetUser(current_user_name string) error {

	current_config, err := Read()
	if err != nil {
		return err
	}

	current_config.CurrentUserName = current_user_name

	data, err := json.Marshal(current_config)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err2 := os.WriteFile(path, data, 0644)
	if err2 != nil {
		return err
	}

	return nil

}
