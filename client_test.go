package dynatrace_client

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const ExtensionTestDataFile = "testdata/custom.my.test.ext.json"

type DynatraceExtensionTest struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Version string `json:"version"`
}

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


