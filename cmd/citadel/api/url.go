package api

import "os"

const (
	API_BASE_URL = "https://console.softwarecitadel.com"
)

func RetrieveApiBaseUrl() string {
	if os.Getenv("API_BASE_URL") != "" {
		return os.Getenv("API_BASE_URL")
	}
	return API_BASE_URL
}
