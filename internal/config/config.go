package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Config struct{
	DbUrl string `json:"db_url"`
	CurrentUsername string `json:"current_username"`
}

func (c *Config) Read()( Config, error){
	filePath, err := getFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("Error: Trying to read config file:\n %v\n", err)
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return Config{}, fmt.Errorf("Error: Trying to unmarshal config file:\n %v\n", err)
	}
	return  *c, nil
}
func (c *Config) SetUser(currentUsername string)error{
	c.CurrentUsername = currentUsername
	filePath, err := getFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return fmt.Errorf("Error: Trying to Marshal config file.\n%v\n", err)
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("Error: Trying to write config file.\n%v\n", err)
	}
	return nil
}

func getFilePath()(string, error){
	filePath := ".gatorconfig.json"
	homeDirPath, err := os.UserHomeDir() 
	if err != nil {
		return "", fmt.Errorf("Error: Trying to get home dir of user:\n %v\n", err)
	}
	return path.Join(homeDirPath, filePath), nil
}
