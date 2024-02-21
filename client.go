package dynatrace_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/jstermask/dynatrace_client/extension"
	"github.com/jstermask/dynatrace_client/model"
)

const ConfigurationApiPath string = "/api/config/v1"

type DynatraceClient struct {
	ApiToken string
	EnvUrl   string
	Client   *http.Client
}

type GetExtensionsResponse struct {
	Extensions   []DynatraceExtensionInfo `json:"extensions"`
	TotalResults int                      `json:"totalResults"`
	NextPageKey  string                   `json:"nextPageKey"`
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

func (c *DynatraceClient) CreateExtension(request *model.DynatraceExtensionRequest) (*model.DynatraceExtensionResponse, error) {
	packagedExtension, err := extension.CreatePackagedExtension(request.Payload) 
	if(err != nil) {
		return nil, err
	}
	defer packagedExtension.Dispose()

	file, _ := os.Open(packagedExtension.FilePath)
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", path.Base(packagedExtension.FilePath))
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", c.getConfigurationApiRequest("/extensions"), body)
	c.addRequestHeaders(req)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	response, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 201 {
		return nil, fmt.Errorf("unable to create extension %s. Status code : %d", request.Name, response.StatusCode)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var dynaResp model.DynatraceExtensionResponse
	err = json.Unmarshal(bodyBytes, &dynaResp)
	if err != nil {
		return nil, err
	}

	return &dynaResp, nil
}

func (c *DynatraceClient) validateConnection() error {
	_, err := c.getExtensions(10, nil)
	return err
}

func (c *DynatraceClient) getExtensions(pageSize int, nextPageKey *string) (*GetExtensionsResponse, error) {
	req, err := http.NewRequest("GET", c.getConfigurationApiRequest("/extensions"), nil)
	c.addRequestHeaders(req)
	if err != nil {
		return nil, err
	}
	queryParams := req.URL.Query()
	queryParams.Set("pageSize", fmt.Sprintf("%d", pageSize))

	if nextPageKey != nil {
		queryParams.Set("nextPageKey", *nextPageKey)
	}

	req.URL.RawQuery = queryParams.Encode()

	response, err := c.Client.Do(req)

	if err != nil {
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("can't retrieve extensions: reason %d %s", response.StatusCode, string(responseBody))
	}

	var extensions GetExtensionsResponse
	err = json.Unmarshal(responseBody, &extensions)
	if err != nil {
		return nil, err
	}

	return &extensions, nil

}

func (c *DynatraceClient) getConfigurationApiRequest(path string) string {
	return fmt.Sprintf("%s%s%s", c.EnvUrl, ConfigurationApiPath, path)
}

func (c *DynatraceClient) addRequestHeaders(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Api-Token %s", c.ApiToken))
	req.Header.Set("accept", "application/json; charset=utf-8")
}
