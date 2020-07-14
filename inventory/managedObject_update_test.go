package inventory

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type UpdateResponse struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Self             string       `json:"self"`
	Type             string       `json:"type"`
	LastUpdated      time.Time    `json:"lastUpdated"`
	StrongTypedClass struct{}     `json:"com_othercompany_StrongTypedClass"`
	ChildDevices     ChildDevices `json:"childDevices"`
}

func TestManagedObjectApi_UpdateManagedObject(t *testing.T) {
	lastUpdated, _ := time.Parse(time.RFC3339, "2019-08-23T15:10:00.653Z")

	tests := []struct {
		name        string
		deviceID    string
		c8yRespCode int
		c8yRespBody string
		expectedErr *generic.Error

		Update                 *Update
		expectedUpdateResponse UpdateResponse
	}{
		{
			name:        "Happy",
			deviceID:    "104940",
			c8yRespCode: http.StatusOK,
			c8yRespBody: UpdateManagedObject("FlowerCare1"),
			Update: &Update{
				Type: "c8y_DeviceGroup",
				Name: "FlowerCare",
			},
			expectedUpdateResponse: UpdateResponse{
				ID:               "104940",
				Name:             "FlowerCare1",
				Self:             "https://t200588189.cumulocity.com/inventory/managedObjects/104940",
				Type:             "c8y_DeviceGroup",
				LastUpdated:      lastUpdated,
				StrongTypedClass: struct{}{},
				ChildDevices: ChildDevices{
					Self: "https://t200588189.cumulocity.com/inventory/managedObjects/104940/childDevices",
				},
			},
		},
		{
			name:        "Unhappy - no ids and references given",
			expectedErr: clientError("given deviceID is empty", "UpdateManagedObject"),
		},
		{
			name:        "Unhappy - status is not StatusOk",
			deviceID:    "104940",
			c8yRespCode: http.StatusUnsupportedMediaType,
			c8yRespBody: `{"error": "undefined/validationError","message": "Representation must not be null","info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"}`,
			expectedErr: &generic.Error{
				ErrorType: "undefined/validationError",
				Message:   "Representation must not be null",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqURL string

			testserver := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				reqURL = fmt.Sprintf("%v", request.URL.String())

				response.WriteHeader(tt.c8yRespCode)
				_, err := response.Write([]byte(tt.c8yRespBody))
				if err != nil {
					t.Fatalf("failed to write resp body: %s", err)
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

			updateResponse, err := c.UpdateManagedObject(tt.deviceID, tt.Update)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("received an unexpected error: expected: %v, actual: %v", tt.expectedErr, err)
			}

			if err == nil && len(tt.deviceID) > 0 {
				expectedReqURL := fmt.Sprintf("%v/%v", managedObjectPath, tt.deviceID)
				if reqURL != expectedReqURL {
					t.Fatalf("unexpected request url: %v expected: %v", reqURL, expectedReqURL)
				}

				if updateResponse.ID != tt.expectedUpdateResponse.ID {
					t.Errorf("ID is incorrect: expected: %v, actual: %v", tt.expectedUpdateResponse.ID, updateResponse.ID)
				}
				if updateResponse.Name != tt.expectedUpdateResponse.Name {
					t.Errorf("Name is incorrect: expected: %v, actual: %v", tt.expectedUpdateResponse.Name, updateResponse.Name)
				}
				if updateResponse.Self != tt.expectedUpdateResponse.Self {
					t.Errorf("Self is incorrect: expected: %v, actual: %v", tt.expectedUpdateResponse.Self, updateResponse.Self)
				}
				if updateResponse.Type != tt.expectedUpdateResponse.Type {
					t.Errorf("ID is incorrect: expected: %v, actual: %v", tt.expectedUpdateResponse.ID, updateResponse.ID)
				}
				if updateResponse.LastUpdated != tt.expectedUpdateResponse.LastUpdated {
					t.Errorf("ID is incorrect: expected: %v, actual: %v", tt.expectedUpdateResponse.ID, updateResponse.ID)
				}
				if updateResponse.ChildDevices.Self != tt.expectedUpdateResponse.ChildDevices.Self {
					t.Errorf("ID is incorrect: expected: %v, actual: %v", tt.expectedUpdateResponse.ID, updateResponse.ID)
				}
			}
		})
	}
}
