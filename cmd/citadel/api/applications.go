package api

import (
	"bytes"
	"citadel/app/models"
	"citadel/cmd/citadel/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func RetrieveApplications() ([]models.Application, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return nil, err
	}

	url := RetrieveApiBaseUrl() + "/apps"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("[1]", err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[2]", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("[3]", err)
		return nil, err
	}

	var applications []models.Application
	err = json.NewDecoder(resp.Body).Decode(&applications)
	if err != nil {
		fmt.Println("[4]", err)
		return nil, err
	}

	return applications, nil
}

func CreateApplication(
	name string,
	cpu string,
	memory string,
) (models.Application, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return models.Application{}, err
	}

	url := RetrieveApiBaseUrl() + "/apps"
	payload := `{"name": "` + name + `",`
	payload += `"cpu": "` + cpu + `",`
	payload += `"ram": "` + memory + `"}`

	body := bytes.NewBuffer([]byte(payload))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return models.Application{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Application{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		var output map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&output)
		if err != nil {
			return models.Application{}, err
		}
		if output["message"] != nil {
			return models.Application{}, errors.New(output["message"].(string))
		}
	}

	if resp.StatusCode != 200 {
		return models.Application{}, errors.New("an error occurred while creating the application")
	}

	var application models.Application
	err = json.NewDecoder(resp.Body).Decode(&application)
	if err != nil {
		return models.Application{}, err
	}

	return application, nil
}
