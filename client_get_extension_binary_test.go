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
		assert.Equal(t, "/api/config/v1/extensions/custom.my.test.ext/binary", r.URL.Path, "wrong server path")
		assert.Equal(t, "GET", r.Method)

		testfile,_ := os.Open(ExtensionTestDataFile)
		payload,_ := io.ReadAll(testfile)

		packagedExtension, err := extension.CreatePackagedExtension(string(payload))
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
		Id: "custom.my.test.ext",
	}
	
	resp, err := client.GetExtensionBinary(&req)
	if err != nil && !shouldFail {
		t.Fatalf("Fail %v", err)
	}

	if !shouldFail {
		assert.Equal(t, "custom.my.test.ext", resp.Id, "Wrong id")
		var values DynatraceExtensionTest
		json.Unmarshal([]byte(resp.Payload), &values)
		assert.Equal(t, "custom.my.test.ext", values.Name, "Wrong name")
		assert.Equal(t, "1.0", values.Version, "Wrong version")
	}
}
