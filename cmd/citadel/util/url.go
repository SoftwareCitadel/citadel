package util

import (
	"errors"
	"net/url"

	"github.com/spf13/viper"
)

var UrlValidateFunc = func(s string) error {
	_, err := url.Parse(s)
	if err != nil {
		return errors.New("URL must be a valid URL")
	}

	return nil
}

func StoreConsoleUrl(consoleUrl string) error {
	_, err := InitConfigDir()
	if err != nil {
		return err
	}

	viper.Set("console_url", consoleUrl)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
