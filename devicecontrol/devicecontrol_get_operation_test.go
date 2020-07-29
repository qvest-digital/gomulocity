package devicecontrol

import (
	"github.com/tarent/gomulocity/generic"
	"testing"
)

func TestDeviceControl_GetOperation_With_OperationID(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, operation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	operation, err := api.GetOperation(operationID)

	if err != nil {
		t.Fatalf("Get() got an unexpected error: %s", err.Error())
	}

	if operation == nil {
		t.Fatalf("Get() returns an unexpected nil operation.")
	}

	if operation.OperationID != operationID {
		t.Errorf("Get() operation id id = %v, want %v", operation.OperationID, operationID)
	}
}

func TestDeviceControl_GetOperation_Without_OperationID(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, operation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	_, err := api.GetOperation("")

	expectedErr := generic.ClientError("Getting operation without an id is not allowed", "GetOperation")
	if err == nil {
		t.Fatal("Error is unexpectedly nil.")
	} else if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: expected: %v, actual: %v", expectedErr.Error(), err.Error())
	}
}

func TestDeviceControl_GetOperation_Invalid_StatusCode(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(500, erroneousResponse)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	_, err := api.GetOperation("operationID")

	expectedError := generic.CreateErrorFromResponse([]byte(erroneousResponse), 500)
	if err == nil {
		t.Fatal("Error is unexpectedly nil.")
	} else if err.Error() != expectedError.Error() {
		t.Fatalf("Received an unexpected error: expected: %v, actual: %v", expectedError.Error(), err.Error())
	}
}
