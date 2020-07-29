package devicecontrol

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"testing"
)

func TestDeviceControl_CreateOperation(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusCreated, operation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	newOperation := &NewOperation{
		DeviceID: "4788195",
		AdditionalFields: map[string]interface{}{
			"custom1": "hello",
		},
	}
	operation, err := api.CreateOperation(newOperation)
	if err != nil {
		t.Errorf("Received an unexpected error: %s", err)
	}

	if len(operation.AdditionalFields) != 2 {
		t.Fatalf("Get() AdditionalFields length = %d, want %d", len(operation.AdditionalFields), 2)
	}

	custom1, ok1 := operation.AdditionalFields["custom1"].(string)

	if !(ok1 && custom1 == "Hello") {
		t.Errorf("Get() custom1 = %v, want %v", custom1, "Hello")
	}
}

func TestDeviceControl_CreateOperation_invalid_status(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusInternalServerError, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	newOperation := &NewOperation{
		DeviceID: "4788195",
		AdditionalFields: map[string]interface{}{
			"custom1": "hello",
		},
	}
	_, err := api.CreateOperation(newOperation)
	if err != nil {
		if err.Error() != generic.ClientError("given response body is empty", "CreateErrorFromResponse").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}
