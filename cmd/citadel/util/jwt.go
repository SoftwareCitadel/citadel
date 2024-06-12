package util

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func StoreJWTToken(jwt string) error {
	_, err := InitConfigDir()
	if err != nil {
		return err
	}

	viper.Set("jwt", jwt)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func InitConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(homeDir, ".citadel")

	if !directoryExists(dir) {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return "", err
		}
	}

	// Initialize viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	// Create config file if it doesn't exist
	if _, err := os.Stat(filepath.Join(dir, "config.yaml")); os.IsNotExist(err) {
		if _, err := os.Create(filepath.Join(dir, "config.yaml")); err != nil {
			return "", err
		}
	}

	return dir, nil
}

func directoryExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func RemoveConfigFile() error {
	configDir, err := InitConfigDir()
	if err != nil {
		return err
	}

	return os.Remove(filepath.Join(configDir, "config.yaml"))
}
