package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) RemoveUser() error {
	c.CurrentUserName = ""
	return writeConfig(c)
}

func (c *Config) SetUser(currentUserName string) error {
	currentUserName = strings.Trim(currentUserName, " ")
	if currentUserName == "" {
		return fmt.Errorf("User name cannot be empty.")
	}
	c.CurrentUserName = currentUserName
	return writeConfig(c)
}
func writeConfig(cnfg *Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(cnfg)
	if err != nil {
		return fmt.Errorf("Error marshaling the data: %w", err)
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("Error writing the config file: %w", err)
	}
	return nil
}
