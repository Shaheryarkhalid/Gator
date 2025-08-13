package config

import "testing"

func TestConfigSetUser(t *testing.T) {
	config, err := Read()
	if err != nil {
		t.Errorf("Error Testing Set User method on Config struct: %v", err)
	}
	err = config.SetUser("Shaheryar")
	if err != nil {
		t.Errorf("Error Testing Set User method on Config struct: %v", err)
		t.FailNow()
	}
	testConfig, err := Read()
	if err != nil {
		t.Errorf("Error Testing Set User method on Config struct: %v", err)
		t.FailNow()
	}
	if testConfig.CurrentUserName == "" {
		t.Error("Test failed as Current User Name does not exist on config file. ")
		t.FailNow()
	}
	if testConfig.DbUrl == "" {
		t.Error("Test failed as database url does not exist on config file. ")
		t.FailNow()
	}
}

func TestRead(t *testing.T) {
	config, err := Read()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if config.DbUrl == "" {
		t.Error("Database url not found in config file.")
		t.FailNow()
	}
}

func TestGetConfigFilePath(t *testing.T) {
	path, err := getConfigFilePath()
	if err != nil {
		t.Errorf("Test getting the config file path failed: %v", err)
		t.FailNow()
	}
	if path == "" {
		t.Errorf("Test getting the config file path failed: Method returned the empty path")
		t.FailNow()
	}
}
