package auth

import (
	"encoding/json"
	"net/http"

	"citadel/cmd/citadel/api"
	"citadel/cmd/citadel/util"
)

type CheckAuthResponse struct {
	Authenticated bool `json:"authenticated"`
}

func IsLoggedIn() bool {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return false
	}

	return checkAuthenticationTokenAgainstAPI(token)
}

func checkAuthenticationTokenAgainstAPI(token string) bool {
	url := api.RetrieveApiBaseUrl() + "/auth/cli/check"
	bearer := "Bearer " + token

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false
	}

	var checkAuthResponse CheckAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&checkAuthResponse)
	if err != nil {
		return false
	}

	if !checkAuthResponse.Authenticated {
		return false
	}

	return true
}
