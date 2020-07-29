package devicecontrol

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"testing"
)

func TestDeviceControl_UpdateOperation(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusOK, `{"status":"SUCCESSFUL"}`)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	operation := &UpdateOperation{
		Status: "status",
		AdditionalFields: map[string]interface{}{
			"custom1": "hello",
		},
	}

	updateStatus, err := api.UpdateOperation(operationID, operation)
	if err != nil {
		t.Errorf("Received an unexpected error: %s", err)
	}

	if updateStatus != "SUCCESSFUL" {
		t.Errorf("Received an unexpected update status. Expected: %v, actual: %v", "SUCCESSFUl", updateStatus)
	}
}

func TestDeviceControl_UpdateOperation_empty_id(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusOK, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	operation := &UpdateOperation{
		Status: "status",
		AdditionalFields: map[string]interface{}{
			"custom1": "hello",
		},
	}

	_, err := api.UpdateOperation("", operation)
	if err != nil {
		if err.Error() != generic.ClientError("Updating operation without an id is not allowed", "UpdateOperation").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}

func TestDeviceControl_UpdateOperation_erroneous_update_status(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusOK, "<invalid json>")
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	operation := &UpdateOperation{
		Status: "status",
		AdditionalFields: map[string]interface{}{
			"custom1": "hello",
		},
	}

	_, err := api.UpdateOperation(operationID, operation)
	if err != nil {
		if err.Error() != generic.ClientError("Error while unmarshalling update status: invalid character '<' looking for beginning of value", "UpdateOperation").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}
