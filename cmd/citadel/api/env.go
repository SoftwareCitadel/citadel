package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"citadel/cmd/citadel/util"
)

func RetrieveEnvironmentVariables(appSlug string) (map[string]string, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return nil, err
	}

	url := RetrieveApiBaseUrl() + "/apps/" + appSlug + "/env"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	var response struct {
		Env map[string]string `json:"env"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Env, nil
}

func SetEnvironmentVariables(
	appSlug string,
	args []string,
) (bool, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return false, err
	}

	data := "{"
	for i, arg := range args {
		// Check arg follows the pattern <key>=<value>
		if !strings.Contains(arg, "=") {
			return false, errors.New("Invalid argument: " + arg)
		}

		splittedArg := strings.Split(arg, "=")
		key := splittedArg[0]
		value := splittedArg[1]

		// Remove single or double quotes to value
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.TrimPrefix(value, "\"")
			value = strings.TrimSuffix(value, "\"")
		} else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
			value = strings.TrimPrefix(value, "'")
			value = strings.TrimSuffix(value, "'")
		}

		data += "\"" + key + "\":\"" + value + "\""
		if i < len(args)-1 {
			data += ","
		}
	}
	data += "}"

	body := bytes.NewBufferString(data)

	url := RetrieveApiBaseUrl() + "/apps/" + appSlug + "/env"
	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return false, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, errors.New("failed to set environment variable")
	}

	// Parse the response body { showRedeployChoice: true }
	var response struct {
		ShowRedeployChoice bool `json:"showRedeployChoice"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, err
	}

	return response.ShowRedeployChoice, nil
}
