package dynatrace_client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/jstermask/dynatrace_client/extension"
	"github.com/stretchr/testify/assert"
)

func TestGetExtensionBinarySuccess(t *testing.T) {
	getExtensionWitResponseCode(t, false)
}

func TestGetExtensionBinaryFail(t *testing.T) {
	getExtensionWitResponseCode(t, true)
}

func getExtensionWitResponseCode(t *testing.T, shouldFail bool) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedToken := r.Header.Get("Authorization")
		assert.Equal(t, "Api-Token SomeToken", receivedToken, "wrong token received")
		assert.Equal(t, "/api/config/v1/extensions/custom.jmx.test.jvm/binary", r.URL.Path, "wrong server path")
		assert.Equal(t, "GET", r.Method)

		packagedExtension, err := extension.CreatePackagedExtension(`
		{
			"name": "custom.jmx.test.jvm",
			"description" : "my extension"
		}
		`)
		if err != nil {
			t.Fatalf("Fail %v", err)
		}

		

		file, err := os.Open(packagedExtension.FilePath)
		if err != nil {
			t.Fatalf("Fail %v", err)
		}
		respBytes, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("Fail %v", err)
		}

		if !shouldFail {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(400)
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

	req := DynatraceExtensionGetBinaryRequest{
		Id: "custom.jmx.test.jvm",
	}
	
	resp, err := client.GetExtensionBinary(&req)
	if err != nil && !shouldFail {
		t.Fatalf("Fail %v", err)
	}

	if !shouldFail {
		assert.Equal(t, "custom.jmx.test.jvm", resp.Id, "Wrong id")
		var values map[string]string
		json.Unmarshal([]byte(resp.Payload), &values)
		assert.Equal(t, "custom.jmx.test.jvm", values["name"], "Wrong name")
		assert.Equal(t, "my extension", values["description"], "Wrong name")
	}
}
