package identity

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIdentity_Delete_ExternalId_Success(t *testing.T) {
	var capturedUrl string
	externalId := "someExternalId"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusNoContent)
	}))

	// given: A test server
	defer ts.Close()

	// and: the api as system under test
	api := buildIdentityAPI(ts.URL)

	err := api.DeleteExternalID("someType", externalId)

	if err != nil {
		t.Fatalf("DeleteExternalID() got an unexpected error: %s", err.Error())
	}

	if strings.Contains(capturedUrl, externalId) == false {
		t.Errorf("DeleteExternalID() The target URL does not contains the Id: url: %s - expected externalID %s", capturedUrl, externalId)
	}
}

func TestIdentity_Delete_ExternalId_NotFound(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusNotFound, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildIdentityAPI(ts.URL)

	err := api.DeleteExternalID("someType", "nonexisting")

	if err == nil {
		t.Errorf("DeleteExternalID() expected error on 404 - not found")
		return
	}
}
