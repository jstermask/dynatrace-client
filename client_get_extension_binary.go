package dynatrace_client

import (
	"fmt"
	"github.com/jstermask/dynatrace_client/extension"
	"io"
	"net/http"
)

type DynatraceExtensionGetBinaryRequest struct {
	Id string
}

type DynatraceExtensionGetBinaryResponse struct {
	Id      string
	Payload string
}

func (c *DynatraceClient) GetExtensionBinary(request *DynatraceExtensionGetBinaryRequest) (*DynatraceExtensionGetBinaryResponse, error) {
	req, err := http.NewRequest("GET", c.getConfigurationApiRequest("/extensions/"+request.Id+"/binary"), nil)
	c.addRequestHeaders(req)
	req.Header.Set("accept", "*/*")
	if err != nil {
		return nil, err
	}

	response, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unable to get extension binary %s. Status code : %d, Message : %s", request.Id, response.StatusCode, string(bodyBytes))
	}

	extensionPayload, err := extension.GetExtensionPayloadFromPackage(bodyBytes)
	if err != nil {
		return nil, err
	}

	return &DynatraceExtensionGetBinaryResponse{
		Id:      request.Id,
		Payload: extensionPayload.Payload,
	}, nil

}
