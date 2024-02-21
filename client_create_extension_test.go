package dynatrace_client

import (
	"archive/zip"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateExtensionSuccess(t *testing.T) {
	createExtensionWitResponseCode(t, false)
}

func TestCreateExtensionFails(t *testing.T) {
	createExtensionWitResponseCode(t, true)
}

func createExtensionWitResponseCode(t *testing.T, shouldFail bool) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedToken := r.Header.Get("Authorization")
		assert.Equal(t, "Api-Token SomeToken", receivedToken, "wrong token received")
		assert.Equal(t, "/api/config/v1/extensions", r.URL.Path, "wrong server path")
		assert.Equal(t, "POST", r.Method)

		file, headers, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("Fail %v", err)
		}

		zipReader, err := zip.NewReader(file, headers.Size)
		if err != nil {
			t.Fatalf("Fail %v", err)
		}

		jsonFile, err := zipReader.Open("plugin.json")
		if err != nil {
			t.Fatalf("Fail %v", err)
		}

		jsonFileBytes, err := io.ReadAll(jsonFile)
		if err != nil {
			t.Fatalf("Fail %v", err)
		}

		var values map[string]string

		err = json.Unmarshal(jsonFileBytes, &values)
		if err != nil {
			t.Fatalf("Fail %v", err)
		}

		assert.Equal(t, "custom.jmx.test.jvm", values["name"], "name is wrong")
		assert.Equal(t, "my extension", values["description"], "description is wrong")

		response := DynatraceExtensionCreateResponse{
			Id:          "SomeId",
			Name:        "SomeName",
			Description: "SomeDescription",
		}

		respBytes, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("Fail %v", err)
		}
		if !shouldFail {
			w.WriteHeader(201)
		}
		w.Write(respBytes)

	}))
	defer svr.Close()

	token := "SomeToken"

	client := DynatraceClient{
		ApiToken: token,
		EnvUrl:   svr.URL,
		Client:   &http.Client{},
	}

	req := DynatraceExtensionCreateRequest{
		Payload: `
		{
			"name": "custom.jmx.test.jvm",
			"description": "my extension"
		}
		`,
	}

	resp, err := client.CreateExtension(&req)
	if err != nil && !shouldFail {
		t.Fatalf("Fail %v", err)
	}

	if !shouldFail {
		assert.Equal(t, "SomeId", resp.Id, "Wrong id")
		assert.Equal(t, "SomeName", resp.Name, "Wrong name ")
		assert.Equal(t, "SomeDescription", resp.Description, "Wrong description")
	}
}
