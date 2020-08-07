package devicecontrol

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeviceControl_DeleteOperationCollection_Happy_single_query_values(t *testing.T) {
	var capturedUrl string
	url := "/devicecontrol/operations?operationsByAgentId=agentID&operationsByDeviceId=deviceID&operationsByStatus=status"
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusNoContent)
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	query := OperationQuery{
		DeviceID: "deviceID",
		Status:   "status",
		AgentID:  "agentID",
	}

	err := api.DeleteOperationCollection(query)
	if err != nil {t.Errorf("invalid url. Expected: %v, actual: %v", url, capturedUrl)
		t.Errorf("received an unexpected error: %s", err)
	}

	if capturedUrl != url {
		t.Errorf("invalid url. Expected: %v, actual: %v", url, capturedUrl)
	}
}

func TestDeviceControl_DeleteOperationCollection_Happy_combined_query_values(t *testing.T) {
	var capturedUrl string
	url := "/devicecontrol/operations?operationsByDeviceIdAndStatus=deviceIDAndStatus"
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusNoContent)
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	query := OperationQuery{
		DeviceIDAndStatus: "deviceIDAndStatus",
	}

	err := api.DeleteOperationCollection(query)
	if err != nil {
		t.Errorf("received an unexpected error: %s", err)
	}

	if capturedUrl != url {
		t.Errorf("invalid url. Expected: %v, actual: %v", url, capturedUrl)
	}
}

func TestDeviceControl_DeleteOperationCollection_Unhappy_invalid_status(t *testing.T) {
	var capturedUrl string
	url := "/devicecontrol/operations?operationsByAgentId=agentID&operationsByDeviceId=deviceID&operationsByStatus=status"

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	query := OperationQuery{
		DeviceID: "deviceID",
		Status:   "status",
		AgentID:  "agentID",
	}

	err := api.DeleteOperationCollection(query)
	if err != nil {
		if err.Error() != generic.ClientError("given response body is empty", "CreateErrorFromResponse").Error() {
			t.Errorf("received an unexpected error: %s", err)
		}
	}

	if capturedUrl != url {
		t.Errorf("invalid url. Expected: %v, actual: %v", url, capturedUrl)
	}
}
