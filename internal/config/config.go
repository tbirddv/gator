package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if homeDir == "" {
		return "", os.ErrNotExist
	}
	return filepath.Join(homeDir, configFileName), nil
}

func Read() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	if configPath == "" {
		return nil, os.ErrNotExist
	}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func write(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	if configPath == "" {
		return os.ErrNotExist
	}
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return err
	}
	return nil
}

func (c *Config) SetUser(UserName string) error {
	c.CurrentUserName = UserName
	return write(c)
}

func (c *Config) ClearUser() error {
	c.CurrentUserName = ""
	return write(c)
}
