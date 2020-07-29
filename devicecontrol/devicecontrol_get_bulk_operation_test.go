package devicecontrol

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"strconv"
	"testing"
)

func TestDeviceControl_GetBulkOperation_Happy(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusOK, bulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	bulkOperation, err := api.GetBulkOperation(bulkOperationID)
	if err != nil {
		t.Errorf("Received an unexpected error: %s", err)
	}

	if bulkOperation == nil {
		t.Error("bulkOperation is unexpectedly nil")
	}

	if strconv.FormatInt(int64(bulkOperation.ID), 10) != bulkOperationID {
		t.Errorf("Invalid BulkoperationID. Expected: %v, actual: %v", bulkOperationID, bulkOperation.ID)
	}
}

func TestDeviceControl_GetBulkOperation_invalid_status(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusInternalServerError, erroneousResponseBulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	_, err := api.GetBulkOperation(bulkOperationID)
	if err != nil {
		if err.Error() != generic.CreateErrorFromResponse([]byte(erroneousResponseBulkOperation), http.StatusInternalServerError).Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}

func TestDeviceControl_GetBulkOperation_empty_id(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusInternalServerError, erroneousResponseBulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	_, err := api.GetBulkOperation("")
	if err != nil {
		if err.Error() != generic.ClientError("Getting bulk operation without a bulkOperationID is not allowed", "GetBulkOperation").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}
