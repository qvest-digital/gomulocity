package inventory

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestManagedObjectApi_DeleteManagedObject(t *testing.T) {
	tests := []struct {
		name        string
		deviceID    string
		respCode    int
		respBody    string
		expectedErr *generic.Error
	}{
		{
			name:     "Happy",
			deviceID: "104940",
			respCode: http.StatusNoContent,
		},
		{
			name:        "Unhappy - no deviceID given",
			expectedErr: clientError("given deviceID is empty", "DeleteManagedObject"),
		},
		{
			name:     "Unhappy - status is not StatusNoContent",
			deviceID: "104940",
			respCode: http.StatusBadRequest,
			respBody: `{"error": "inventory/Not Found","message": "Finding device data from database failed : No managedObject for id '213213213213213'!","info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"}`,
			expectedErr: &generic.Error{
				ErrorType: "inventory/Not Found",
				Message:   "Finding device data from database failed : No managedObject for id '213213213213213'!",
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
					t.Fatalf("failed to write resp body: %s", err)
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

			err := c.DeleteManagedObject(tt.deviceID)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("received an unexpected error: expected: %v, actual: %v", tt.expectedErr, err)
			}

			if len(tt.deviceID) > 0 {
				expectedReqURL := fmt.Sprintf("%v/%v", managedObjectPath, tt.deviceID)
				if reqURL != expectedReqURL {
					t.Fatalf("unexpected request url: %v expected: %v", reqURL, expectedReqURL)
				}
			}
		})
	}
}
