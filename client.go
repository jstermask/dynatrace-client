package dynatrace_client

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
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

func (c *DynatraceClient) CreateExtension(request *DynatraceExtensionRequest) (*DynatraceExtensionResponse, error) {
	var metadata DynatraceExtensionMetadata
	err := json.Unmarshal([]byte(request.Payload), &metadata)
	if err != nil {
		return nil, err
	}

	// create a zip file containing a plugin.json file with payload content
	zipDir, err := os.MkdirTemp(os.TempDir(), "dynatrace_extension")

	zipFilePath := fmt.Sprintf("%s/%s.zip", zipDir, metadata.Name)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return nil, err
	}
	zipWriter := zip.NewWriter(zipFile)

	entry, err := zipWriter.Create("plugin.json")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(entry, strings.NewReader(request.Payload))
	if err != nil {
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}
	zipFile.Close()
	if err != nil {
		return nil, err
	}

	file, _ := os.Open(zipFilePath)
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", path.Base(zipFilePath))
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", c.getConfigurationApiRequest("/extensions"), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	response, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var dynaResp DynatraceExtensionResponse
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
