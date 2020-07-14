package inventory

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestManagedObjectApi_ManagedObjectByID(t *testing.T) {
	tests := []struct {
		name        string
		deviceID    string
		respCode    int
		respBody    string
		expectedErr *generic.Error
	}{
		{
			name:     "Happy",
			deviceID: "deviceID",
			respCode: http.StatusOK,
			respBody: ManagedObjectByID,
		},
		{
			name:     "Unhappy - no deviceID",
			respCode: http.StatusOK,
			expectedErr: clientError("given deviceID is empty", "ManagedObjectByID"),
		},
		{
			name:     "Unhappy - statuscode is not statusOK",
			deviceID: "deviceID",
			respCode: http.StatusBadRequest,
			respBody: `{"error":"inventory/Not Found", "message":"Finding device data from database failed : No managedObject for id '1'!", "info":"https://www.cumulocity.com/guides/reference-guide/#error_reporting"}`,
			expectedErr: &generic.Error{
				ErrorType: "inventory/Not Found",
				Message:   "Finding device data from database failed : No managedObject for id '1'!",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqURL string

			testserver := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				reqURL = fmt.Sprintf("%v", request.URL.String())

				response.WriteHeader(tt.respCode)
				_, err := response.Write([]byte(tt.respBody))
				if err != nil {
					t.Fatal("failed to write response body")
				}
			}))
			defer testserver.Close()

			c := inventoryApi{
				Client: &generic.Client{
					HTTPClient: testserver.Client(),
					BaseURL:    testserver.URL,
				},
				ManagedObjectsPath: managedObjectPath,
			}

			managedObject, err := c.ManagedObjectByID(tt.deviceID)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("received an unexpected error: expected: %v, actual: %v", tt.expectedErr, err)
			}

			if len(tt.deviceID) != 0 {
				expectedReqURL := fmt.Sprintf("%v/%v", managedObjectPath, tt.deviceID)
				if reqURL != expectedReqURL {
					t.Fatalf("unexpected request url: %v expected: %v", reqURL, expectedReqURL)
				}
			}

			if err == nil {
				if len(managedObject.Name) == 0 {
					t.Error("value of key 'name' is empty")
				}
				if len(managedObject.ID) == 0 {
					t.Error("value of key 'name' is empty")
				}
				if managedObject.C8YStatus.LastUpdated.Date.Date.IsZero() {
					t.Error("lastUpdated values is not set")
				}
			}
		})
	}
}
