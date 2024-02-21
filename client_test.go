package dynatrace_client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jstermask/dynatrace_client/model"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedToken := r.Header.Get("Authorization")
		assert.Equal(t, "Api-Token SomeToken", receivedToken, "wrong token received")
		assert.Equal(t, "/api/config/v1/extensions", r.URL.Path, "wrong server path")

		response := GetExtensionsResponse{
			Extensions: []DynatraceExtensionInfo{
				{
					Id:   "1",
					Name: "My extension",
					Type: "JMX",
				},
			},
			TotalResults: 1,
			NextPageKey:  "",
		}

		respBytes, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("Fail %v", err)
		}
		w.Write(respBytes)

	}))
	defer svr.Close()

	token := "SomeToken"

	_, err := NewClient(&svr.URL, &token)

	if err != nil {
		t.Fatalf("Fail %v", err)
	}

}

func TestCreateExtension(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedToken := r.Header.Get("Authorization")
		assert.Equal(t, "Api-Token SomeToken", receivedToken, "wrong token received")
		assert.Equal(t, "/api/config/v1/extensions", r.URL.Path, "wrong server path")
		assert.Equal(t, "POST", r.Method)

		response := model.DynatraceExtensionResponse{
			Id: "SomeId",
			Name: "SomeName",
			Description: "SomeDescription",
		}

		respBytes, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("Fail %v", err)
		}
		w.Write(respBytes)

	}))
	defer svr.Close()

	token := "SomeToken"

	client := DynatraceClient{
		ApiToken: token,
		EnvUrl: svr.URL,
		Client: &http.Client{},
	}
	
	req := model.DynatraceExtensionRequest{
		Name: "custom.jmx.ipa.test",
		Payload: `
		{
			"name": "custom.jmx.ipa.jvm",
			"description": "my extension"
		}
		`,
	}

	resp, err := client.CreateExtension(&req)
	if(err != nil) {
		t.Fatalf("Fail %v", err)
	}
	assert.Equal(t, "SomeId", resp.Id, "Wrong id")
	assert.Equal(t, "SomeName", resp.Name, "Wrong name ")
	assert.Equal(t, "SomeDescription", resp.Description, "Wrong description")
	
}
