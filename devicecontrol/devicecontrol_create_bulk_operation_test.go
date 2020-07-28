package devicecontrol

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"testing"
	"time"
)

func TestDeviceControl_CreateBulkOperation(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusCreated, bulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	startDate, _ := time.Parse(time.RFC3339, "2020-01-23T12:29:35.387Z")
	newBulkOperation := &NewBulkOperation{
		StartDate:    startDate,
		CreationRamp: 15,
		OperationPrototype: map[string]interface{}{
			"DeliveryType": "SMS",
			"Command": struct {
				Text string
			}{
				Text: "test",
			},
			"Description": "Execute shell command",
		},
	}

	bulkOperation, err := api.CreateBulkOperation(newBulkOperation)
	if err != nil {
		t.Errorf("Received an unexpected error: %s", err)
	}

	if bulkOperation == nil {
		t.Error("bulk operation is unexpectedly nil")
	}

	if newBulkOperation.StartDate != bulkOperation.StartDate {
		t.Errorf("Received an erroneous startDate. Expected: %v, actual: %v", newBulkOperation.StartDate, bulkOperation.StartDate)
	}
	if newBulkOperation.CreationRamp != bulkOperation.CreationRamp {
		t.Errorf("Received an erroneous creationRamp. Expected: %v, actual: %v", newBulkOperation.CreationRamp, bulkOperation.CreationRamp)
	}

	operationPrototype, ok := bulkOperation.OperationPrototype["operationPrototype"].(map[string]interface{})
	if !ok {
		t.Error("Error while casting OperationPrototype to map[string]interface")
	}
	var deliveryType string
	var description string
	var commandText string

	for key, value := range operationPrototype {
		switch key {
		case "deliveryType":
			deliveryType = value.(string)
		case "description":
			description = value.(string)
		default:
			commandMap, ok := value.(map[string]interface{})
			if !ok {
				t.Error("could not cast interface into map[string]interface")
			}
			commandText, ok = commandMap["text"].(string)
			if !ok {
				t.Error("Error while casting commandMap value to string")
			}
		}
	}

	if deliveryType != "SMS" {
		t.Errorf("Received an unexpected deliveryType. Expected: %v, actual: %v", "SMS", deliveryType)
	}
	if description != "Execute shell command" {
		t.Errorf("Received an unexpected description. Expected: %v, actual: %v", "Execute shell command", description)
	}
	if commandText != "test" {
		t.Errorf("Received an unexpected command-text. Expected: %v, actual: %v", "test", commandText)
	}
}

func TestDeviceControl_CreateBulkOperation_invalid_status(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusInternalServerError, erroneousResponseBulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	_, err := api.CreateBulkOperation(&NewBulkOperation{})
	if err != nil {
		if err.Error() != generic.CreateErrorFromResponse([]byte(erroneousResponseBulkOperation), http.StatusInternalServerError).Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}

func TestDeviceControl_CreateBulkOperation_No_Pointer_Operation(t *testing.T) {
	// given: A test server
	ts := createOperationHttpServer(http.StatusCreated, erroneousResponseBulkOperation)
	defer ts.Close()

	// and: the api as system under test
	api := buildOperationApi(ts.URL)

	_, err := api.CreateBulkOperation(nil)
	if err != nil {
		if err.Error() != generic.ClientError("failed to marshal bulkOperation: input is not a pointer of struct", "CreateBulkOperation").Error() {
			t.Errorf("Received an unexpected error: %s", err)
		}
	}
}
