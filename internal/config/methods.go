package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func Read() (config Config, err error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("Error reading the config file: %w", err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("Error reading the config file: %w", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("Error unmarshaling the config file: %w", err)
	}
	return config, nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error getting user home directory: %w", err)
	}
	return filepath.Join(homePath, configFileName), nil
}
