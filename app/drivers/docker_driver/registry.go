package dockerDriver

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/docker/docker/api/types/registry"
)

func getRegistryAuth() string {
	authConfig := registry.AuthConfig{
		RegistryToken: os.Getenv("REGISTRY_TOKEN"),
		ServerAddress: os.Getenv("REGISTRY_HOST"),
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	return authStr
}
