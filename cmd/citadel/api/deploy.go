package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"citadel/cmd/citadel/util"
)

func DeployFromTarball(tarball io.ReadCloser, orgId string, appSlug string, releaseCmd string) (bool, error) {
	// Retrieve the token from the config file
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return false, err
	}

	// Create a new HTTP request
	url := RetrieveApiBaseUrl() + "/orgs/" + orgId + "/apps/" + appSlug + "/deployments"

	// Create the request, sending the gzipBuf as a form data field named "tarball"
	form := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(form)
	part, err := writer.CreateFormFile("tarball", "tarball")
	if err != nil {
		return false, err
	}
	_, err = io.Copy(part, tarball)
	if err != nil {
		return false, err
	}
	err = writer.Close()
	if err != nil {
		return false, err
	}

	err = writer.WriteField("releaseCommand", releaseCmd)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", url, form)
	if err != nil {
		return false, err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the HTTP request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	var response struct {
		Healthcheck bool `json:"healthcheck"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, err
	}

	return response.Healthcheck, nil
}

func RedeployApplication(
	appSlug string,
) error {
	// Retrieve the token from the config file
	token, err := util.RetrieveTokenFromConfig()
	if err != nil {
		return err
	}

	// Create a new HTTP request
	url := RetrieveApiBaseUrl() + "/api/apps/" + appSlug + "/redeploy"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the HTTP request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	return nil
}
