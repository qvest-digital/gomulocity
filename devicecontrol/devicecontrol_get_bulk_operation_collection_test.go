package devicecontrol

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"testing"
)

func TestDeviceControl_GetCollectionOfBulkOperation(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusOK, bulkOperationCollection)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	query := OperationQuery{
		DeviceID: "deviceID",
		Status:   "status",
		AgentID:  "agentID",
	}
	pageSize := 5

	bulkCollection, err := api.GetCollectionOfBulkOperation(query, pageSize)
	if err != nil {
		t.Errorf("Received an unexpected error: %s", err)
	}

	if len(bulkCollection.BulkOperations) == 0 {
		t.Error("Error: no bulkOperations found")
	}
}

func TestDeviceControl_GetCollectionOfBulkOperation_invalid_status(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusInternalServerError, erroneousResponseBulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	query := OperationQuery{
		DeviceID: "deviceID",
		Status:   "status",
		AgentID:  "agentID",
	}
	pageSize := 5

	_, err := api.GetCollectionOfBulkOperation(query, pageSize)
	if err != nil {
		if err.Error() != generic.CreateErrorFromResponse([]byte(erroneousResponseBulkOperation), http.StatusInternalServerError).Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}

func TestDeviceControl_GetCollectionOfBulkOperation_empty_query(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusOK, erroneousResponseBulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	_, err := api.GetCollectionOfBulkOperation(OperationQuery{}, 5)
	if err != nil {
		if err.Error() != generic.ClientError("No filter set", "FindBulkOperationCollection").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}

func TestDeviceControl_GetCollectionOfBulkOperation_invalid_pageSize(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusOK, erroneousResponseBulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	query := OperationQuery{
		DeviceID: "deviceID",
		Status:   "status",
		AgentID:  "agentID",
	}

	_, err := api.GetCollectionOfBulkOperation(query, 2001)
	if err != nil {
		if err.Error() != generic.ClientError("Error while building pageSize parameter to fetch measurements: The page size must be between 1 and 2000. Was 2001", "FindMeasurements").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}
