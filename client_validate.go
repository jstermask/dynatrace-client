package dynatrace_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetExtensionsResponse struct {
	Extensions   []DynatraceExtensionInfo `json:"extensions"`
	TotalResults int                      `json:"totalResults"`
	NextPageKey  string                   `json:"nextPageKey"`
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