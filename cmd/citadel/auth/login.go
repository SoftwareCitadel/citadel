package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"citadel/cmd/citadel/api"
)

type AuthenticationSessionResponse struct {
	SessionId string `json:"sessionId"`
}

type WaitForLoginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

func GetAuthenticationSessionId() (string, error) {
	resp, err := http.Get(api.RetrieveApiBaseUrl() + "/auth/cli")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var auth AuthenticationSessionResponse

	if err := json.Unmarshal(body, &auth); err != nil {
		return "", err
	}

	return auth.SessionId, nil
}

func WaitForLogin(sessionId string) (string, error) {
	url := api.RetrieveApiBaseUrl() + "/auth/cli/wait/" + sessionId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var wait WaitForLoginResponse

	if err := json.Unmarshal(body, &wait); err != nil {
		return "", err
	}

	if wait.Status == "pending" {
		time.Sleep(1 * time.Second)
		return WaitForLogin(sessionId)
	}

	return wait.Token, nil
}
