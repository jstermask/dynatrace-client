package dynatrace_client

import (
	"fmt"
	"net/http"
)

type DynatraceExtensionDeleteRequest struct {
	Id string
}



func (c *DynatraceClient) DeleteExtension(request *DynatraceExtensionDeleteRequest) (error) {
	req, err := http.NewRequest("DELETE", c.getConfigurationApiRequest("/extensions/"+request.Id), nil)
	c.addRequestHeaders(req)
	req.Header.Set("accept", "*/*")
	if err != nil {
		return err
	}

	response, err := c.Client.Do(req)
	if err != nil {
		return err
	}


	if response.StatusCode != 204 {
		return fmt.Errorf("unable to delete extension %s. Status code : %d", request.Id, response.StatusCode)
	}

	return nil

}
