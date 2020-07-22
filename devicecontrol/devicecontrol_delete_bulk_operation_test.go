package devicecontrol

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeviceControl_DeleteBulkOperation_Happy(t *testing.T) {
	var capturedUrl string
	url := fmt.Sprintf("/devicecontrol/bulkoperations/%v", bulkOperationID)
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusNoContent)
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	err := api.DeleteBulkOperation(bulkOperationID)
	if err != nil {
		t.Errorf("received an unexpected error: %s", err)
	}

	if capturedUrl != url {
		t.Errorf("invalid url. Expected: %v, actual: %v", url, capturedUrl)
	}
}

func TestDeviceControl_DeleteBulkOperation_invalid_status(t *testing.T) {
	var capturedUrl string
	url := fmt.Sprintf("/devicecontrol/bulkoperations/%v", bulkOperationID)
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	err := api.DeleteBulkOperation(bulkOperationID)
	if err != nil {
		if err.Error() != generic.ClientError("given response body is empty", "CreateErrorFromResponse").Error() {
			t.Errorf("received an unexpected error: %s", err)
		}
	}

	if capturedUrl != url {
		t.Errorf("invalid url. Expected: %v, actual: %v", url, capturedUrl)
	}
}

func TestDeviceControl_DeleteBulkOperation_empty_id(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	err := api.DeleteBulkOperation("")
	if err != nil {
		if err.Error() != generic.ClientError("Deleting bulk operation without a bulkOperationID is not allowed", "DeleteBulkOperation").Error() {
			t.Errorf("received an unexpected error: %s", err)
		}
	}
}