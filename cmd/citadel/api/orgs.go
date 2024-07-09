package api

import (
	"bytes"
	"citadel/cmd/citadel/util"
	"citadel/internal/models"
	"encoding/json"
	"errors"
	"net/http"
)

func RetrieveOrgs() ([]models.Organization, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return nil, err
	}

	url := RetrieveApiBaseUrl() + "/orgs"
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

	var orgs []models.Organization
	err = json.NewDecoder(resp.Body).Decode(&orgs)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func CreateOrganization(
	name string,
	cpu string,
	memory string,
) (models.Organization, error) {
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return models.Organization{}, err
	}

	url := RetrieveApiBaseUrl() + "/orgs"
	payload := `{"name": "` + name + `"}`

	body := bytes.NewBuffer([]byte(payload))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return models.Organization{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Organization{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		var output map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&output)
		if err != nil {
			return models.Organization{}, err
		}
		if output["message"] != nil {
			return models.Organization{}, errors.New(output["message"].(string))
		}
	}

	if resp.StatusCode != 200 {
		return models.Organization{}, errors.New("an error occurred while creating the application")
	}

	var org models.Organization
	err = json.NewDecoder(resp.Body).Decode(&org)
	if err != nil {
		return models.Organization{}, err
	}

	return org, nil
}
