package dynatrace_client

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteExtensionSuccess(t *testing.T) {
	deleteExtensionWitResponseCode(t, false)
}

func TestDeleteExtensionFail(t *testing.T) {
	deleteExtensionWitResponseCode(t, true)
}

func deleteExtensionWitResponseCode(t *testing.T, shouldFail bool) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedToken := r.Header.Get("Authorization")
		assert.Equal(t, "Api-Token SomeToken", receivedToken, "wrong token received")
		assert.Equal(t, "/api/config/v1/extensions/custom.jmx.testext", r.URL.Path, "wrong server path")
		assert.Equal(t, "DELETE", r.Method)

		if !shouldFail {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(400)
		}
	}))
	defer svr.Close()

	token := "SomeToken"

	client := DynatraceClient{
		ApiToken: token,
		EnvUrl:   svr.URL,
		Client:   &http.Client{},
	}

	req := DynatraceExtensionDeleteRequest{
		Id: "custom.jmx.testext",
	}

	err := client.DeleteExtension(&req)
	if err != nil && !shouldFail {
		t.Fatalf("Fail %v", err)
	}

}
