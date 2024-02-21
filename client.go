package dynatrace_client

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const ConfigurationApiPath string = "/api/config/v1"

type DynatraceClient struct {
	ApiToken string
	EnvUrl   string
	Client   *http.Client
}

type DynatraceError struct {
	Code    string
	Message string
}


type DynatraceExtensionInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewClient(envUrl *string, apiToken *string) (*DynatraceClient, error) {

	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	client := DynatraceClient{
		Client: &httpClient,
	}

	if envUrl == nil {
		return nil, errors.New("environment URL is missing")
	}

	if apiToken == nil {
		return nil, errors.New("API Token is missing")
	}

	client.ApiToken = *apiToken
	client.EnvUrl = *envUrl

	err := client.validateConnection()
	if err != nil {
		return nil, err
	}

	return &client, nil
}


func (c *DynatraceClient) getConfigurationApiRequest(path string) string {
	return fmt.Sprintf("%s%s%s", c.EnvUrl, ConfigurationApiPath, path)
}

func (c *DynatraceClient) addRequestHeaders(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Api-Token %s", c.ApiToken))
	req.Header.Set("accept", "application/json; charset=utf-8")
}
