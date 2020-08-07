package devicecontrol

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"strconv"
	"testing"
)

func TestDeviceControl_UpdateBulkOperation(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusOK, bulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	bulkOperation := &UpdateBulkOperation{
		CreationRamp: 15,
	}
	operation, err := api.UpdateBulkOperation(bulkOperationID, bulkOperation)
	if err != nil {
		t.Errorf("Received an unexpected error: %s", err)
	}

	if err == nil {
		if operation == nil {
			t.Error("bulkOperation is unexpectedly nil")
		}
	}

	if strconv.FormatInt(int64(operation.ID), 10) != bulkOperationID {
		t.Errorf("Received an unexpected operation ID. Expected: %v, actual: %v", bulkOperationID, operation.ID)
	}

	if operation.CreationRamp != bulkOperation.CreationRamp {
		t.Errorf("Received an unexpected creation ramp. Expected: %v, actual: %v", bulkOperation.CreationRamp, operation.CreationRamp)
	}
}

func TestDeviceControl_UpdateBulkOperation_empty_id(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusOK, bulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	bulkOperation := &UpdateBulkOperation{
		CreationRamp: 15,
	}
	_, err := api.UpdateBulkOperation("", bulkOperation)
	if err != nil {
		if err.Error() != generic.ClientError("Updating bulkOperation without a bulkOperationID is not allowed", "UpdateBulkOperation").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}

func TestDeviceControl_UpdateBulkOperation_invalid_status(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusInternalServerError, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	bulkOperation := &UpdateBulkOperation{
		CreationRamp: 15,
	}
	_, err := api.UpdateBulkOperation(bulkOperationID, bulkOperation)
	if err != nil {
		if err.Error() != generic.ClientError("given response body is empty", "CreateErrorFromResponse").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}

func TestDeviceControl_UpdateBulkOperation_invalid_response(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusOK, "<invalid json>")
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	bulkOperation := &UpdateBulkOperation{
		CreationRamp: 15,
	}
	_, err := api.UpdateBulkOperation(bulkOperationID, bulkOperation)
	if err != nil {
		if err.Error() != generic.ClientError("Error while marshalling response: Error while unmarshalling json: invalid character '<' looking for beginning of value", "UpdateBulkOperation").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}
