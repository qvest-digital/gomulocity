package inventory

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestManagedObjectApi_DeleteReference(t *testing.T) {
	tests := []struct {
		name        string
		deviceID    string
		referenceID string
		reference   string
		respCode    int
		respBody    string
		expectedErr *generic.Error
	}{
		{
			name:        "Happy",
			deviceID:    "104940",
			referenceID: "232704",
			reference:   "childAdditions",
			respCode:    http.StatusNoContent,
		},
		{
			name:        "Unhappy - no ids and references given",
			expectedErr: clientError("given deviceID, reference or referenceID is empty", "DeleteReference"),
		},
		{
			name:        "Unhappy - status is not StatusNoContent",
			deviceID:    "104940",
			referenceID: "232704",
			reference:   "childAdditions",
			respCode:    http.StatusBadRequest,
			respBody:    `{"error":"inventory/Not Found", "message":"Finding device data from database failed : No managedObject for id '23270414'!", "info":"https://www.cumulocity.com/guides/reference-guide/#error_reporting"}`,
			expectedErr: &generic.Error{
				ErrorType: "inventory/Not Found",
				Message:   "Finding device data from database failed : No managedObject for id '23270414'!",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqURL string

			testserver := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				reqURL = fmt.Sprintf("%v", request.URL.String())

				if len(tt.respBody) > 0 {
					response.Write([]byte(tt.respBody))
				}

				response.WriteHeader(tt.respCode)
			}))
			defer testserver.Close()

			c := managedObjectApi{
				Client: &generic.Client{
					HTTPClient: testserver.Client(),
					BaseURL:    testserver.URL,
				},
				ManagedObjectsPath: managedObjectPath,
			}

			err := c.DeleteReference(tt.deviceID, tt.reference, tt.referenceID)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("received an unexpected error: expected: %v, actual: %v", tt.expectedErr, err)
			}

			if len(tt.deviceID) != 0 || len(tt.referenceID) != 0 || len(tt.reference) != 0 {
				expectedReqURL := fmt.Sprintf("%v/%v/%v/%v", managedObjectPath, tt.deviceID, tt.reference, tt.referenceID)
				if reqURL != expectedReqURL {
					t.Fatalf("unexpected request url: %v expected: %v", reqURL, expectedReqURL)
				}
			}
		})
	}
}
