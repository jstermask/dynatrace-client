package dynatrace_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jstermask/dynatrace_client/extension"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

type DynatraceExtensionCreateRequest struct {
	Payload string
}

type DynatraceExtensionCreateResponse struct {
	Id          string
	Name        string
	Description string
	Error       *DynatraceError
}

func (c *DynatraceClient) CreateExtension(request *DynatraceExtensionCreateRequest) (*DynatraceExtensionCreateResponse, error) {
	packagedExtension, err := extension.CreatePackagedExtension(request.Payload)
	if err != nil {
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

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 201 {
		return nil, fmt.Errorf("unable to create extension from zip : %s. Status code : %d. Message : %s. Payload %s", path.Base(packagedExtension.FilePath), response.StatusCode, string(bodyBytes), request.Payload)
	}

	var dynaResp DynatraceExtensionCreateResponse
	err = json.Unmarshal(bodyBytes, &dynaResp)
	if err != nil {
		return nil, err
	}

	return &dynaResp, nil
}
