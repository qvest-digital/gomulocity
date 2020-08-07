package devicecontrol

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeviceControl_GetOperationCollection(t *testing.T) {
	var capturedUrl string
	url := "/devicecontrol/operations?operationsByAgentId=agentID&operationsByDeviceId=deviceID&operationsByStatus=status&pageSize=10"
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(operationCollection))
		if err != nil {
			t.Errorf("Error while writing response: %s", err)
		}
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	query := OperationQuery{
		DeviceID: "deviceID",
		Status:   "status",
		AgentID:  "agentID",
	}
	pageSize := 10

	operationCollection, err := api.GetOperationCollection(query, pageSize)
	if err != nil {
		t.Errorf("Received an unexpected error: %s", err)
	}

	if capturedUrl != url {
		t.Errorf("invalid url. Expected: %v, actual: %v", url, capturedUrl)
	}

	if err == nil {
		if len(operationCollection.Operations) != 3 {
			t.Errorf("Unexpected amount of operations. Expected: %v, actual: %v", 3, len(operationCollection.Operations))
		}
	}
}

func TestDeviceControl_GetOperationCollection_invalid_status(t *testing.T) {
	var capturedUrl string
	url := "/devicecontrol/operations?operationsByAgentId=agentID&operationsByDeviceId=deviceID&operationsByStatus=status&pageSize=10"
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
	pageSize := 10

	operationCollection, err := api.GetOperationCollection(query, pageSize)
	if err != nil {
		if err.Error() != generic.ClientError("given response body is empty", "CreateErrorFromResponse").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}

	if capturedUrl != url {
		t.Errorf("invalid url. Expected: %v, actual: %v", url, capturedUrl)
	}

	if err == nil {
		if len(operationCollection.Operations) != 3 {
			t.Errorf("Unexpected amount of operations. Expected: %v, actual: %v", 3, len(operationCollection.Operations))
		}
	}
}

func TestDeviceControl_GetOperationCollection_empty_query(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	pageSize := 10

	_, err := api.GetOperationCollection(OperationQuery{}, pageSize)
	if err != nil {
		if err.Error() != generic.ClientError("No filter set", "FindOperationCollection").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}
