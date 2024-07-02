package util

import (
	"os"

	"github.com/spf13/viper"
)

func RetrieveTokenFromConfig() (string, error) {
	configDir, err := InitConfigDir()
	if err != nil {
		return "", err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	token := viper.GetString("jwt")

	return token, nil
}

func RetrieveAppSlugFromConfig() (string, error) {
	viper.SetConfigName("citadel")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	appSlug := viper.GetString("application_slug")

	return appSlug, nil
}

func IsAlreadyInitialized() bool {
	vi := viper.New()
	vi.SetConfigName("citadel")
	vi.AddConfigPath(".")
	vi.SetConfigType("toml")

	err := vi.ReadInConfig()

	return err == nil
}

func RetrieveReleaseCommandFromProjectConfig() (string, error) {
	viper.SetConfigName("citadel")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	releaseCommand := viper.GetString("release_command")

	return releaseCommand, nil
}

func InitializeConfigFile(orgId, appSlug string) error {
	vi := viper.New()
	vi.SetConfigName("citadel")
	vi.AddConfigPath(".")
	vi.SetConfigType("toml")

	vi.Set("organization_id", orgId)
	vi.Set("application_slug", appSlug)

	var fileExists bool

	if _, err := os.Stat("citadel.toml"); err == nil {
		fileExists = true
	} else if os.IsNotExist(err) {
		fileExists = false
	} else {
		return err
	}

	if fileExists {
		if err := vi.MergeInConfig(); err != nil {
			return err
		}
	} else {
		if _, err := os.Create("citadel.toml"); err != nil {
			return err
		}
	}

	if err := vi.WriteConfig(); err != nil {
		return err
	}

	return nil
}
