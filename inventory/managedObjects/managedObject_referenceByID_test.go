package managedObjects

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestManagedObjectApi_ReferenceByID(t *testing.T) {
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
			respCode:    http.StatusOK,
			respBody:    ReferenceByID,
		},
		{
			name:        "Unhappy - statuscode is not StatusOK",
			deviceID:    "104940",
			referenceID: "23270",
			reference:   "childAssets",
			respCode:    http.StatusBadRequest,
			respBody:    `{"error":"inventory/Not Found", "message":"Finding device data from database failed : No managedObject for id '23270'!", "info":"https://www.cumulocity.com/guides/reference-guide/#error_reporting"}`,
			expectedErr: &generic.Error{
				ErrorType: "inventory/Not Found",
				Message:   "Finding device data from database failed : No managedObject for id '23270'!",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			},
		},
		{
			name:        "Unhappy - no ids and references given",
			expectedErr: clientError("given deviceID, reference or referenceID is empty", "ReferenceByID"),
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

			c := managedObjectApi{
				Client: &generic.Client{
					HTTPClient: testserver.Client(),
					BaseURL:    testserver.URL,
				},
				ManagedObjectsPath: managedObjectPath,
			}

			reference, err := c.ReferenceByID(tt.deviceID, tt.reference, tt.referenceID)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("received an unexpected error: expected: %v, actual: %v", tt.expectedErr, err)
			}

			if len(tt.deviceID) != 0 || len(tt.referenceID) != 0 || len(tt.reference) != 0 {
				expectedReqURL := fmt.Sprintf("%v/%v/%v/%v", managedObjectPath, tt.deviceID, tt.reference, tt.referenceID)
				if reqURL != expectedReqURL {
					t.Fatalf("unexpected request url: %v expected: %v", reqURL, expectedReqURL)
				}
			}

			if err == nil {
				if len(reference.ManagedObject.ID) == 0 {
					t.Error("value of key 'ID' is empty")
				}
				if len(reference.ManagedObject.Name) == 0 {
					t.Error("value of key 'name' is empty")
				}
				if len(reference.ManagedObject.ChildAdditions.References) == 0 {
					t.Error("no references found for childAdditions")
				}
			}
		})
	}
}
