package managedObjects

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	Next = "https://t200588189.cumulocity.com/inventory/managedObjects?query=has(c8y_IsDevice)&pageSize=3&currentPage=%v"
)

func TestManagedObjectApi_ManagedObjectCollection(t *testing.T) {
	tests := []struct {
		name             string
		filter           ManagedObjectCollectionFilter
		respCode         int
		respBody         string
		expectedRespBody string
		expectedErr      string
	}{
		{
			name: "Happy - no QueryLanguage given",
			filter: ManagedObjectCollectionFilter{
				Type:          "c8y_SensorPhone",
				Owner:         "device_4D8AFED3",
				FragmentType:  "fragmentType",
				QueryLanguage: "",
			},
			respCode: http.StatusOK,
			respBody: NewManagedObjectCollection_ResponseBody(""),
		},
		{
			name: "Happy - QueryLanguage given",
			filter: ManagedObjectCollectionFilter{
				//Filter like Type and Owner should be ignored
				Type:          "c8y_SensorPhone",
				Owner:         "device_4D8AFED3",
				FragmentType:  "fragmentType",
				QueryLanguage: "query=has(c8y_IsDevice)",
			},
			respCode: http.StatusOK,
			respBody: NewManagedObjectCollection_ResponseBody(""),
		},
		{
			name:     "Unhappy - statuscode is not statusOK",
			respCode: http.StatusBadRequest,
			respBody: `{"error":"inventory/Not Found", "message":"Finding device data from database failed : No managedObject for id '1+'!", "info":"https://www.cumulocity.com/guides/reference-guide/#error_reporting"}`,

			expectedErr: generic.Error{
				ErrorType: "inventory/Not Found",
				Message:   "Finding device data from database failed : No managedObject for id '1+'!",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			}.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var username, password, reqURL string

			testserver := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				reqURL = fmt.Sprintf("%v", request.URL.String())

				username, password, _ = request.BasicAuth()
				response.WriteHeader(tt.respCode)
				_, err := response.Write([]byte(tt.respBody))
				if err != nil {
					t.Fatal("failed to write response body")
				}
			}))
			defer testserver.Close()

			u := "<username>"
			p := "<password>"

			c := managedObjectApi{
				Client: &generic.Client{
					HTTPClient: testserver.Client(),
					BaseURL:    testserver.URL,
					Username:   "<username>",
					Password:   "<password>",
				},
				ManagedObjectsPath: managedObjectPath,
			}

			managedObjectCollection, err := c.ManagedObjectCollection(tt.filter)
			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("received an unexpected error: expected: %v, actual: %v", tt.expectedErr, err)
			}

			if username != u {
				t.Fatalf("unexpected c8y auth-username: %v expected: %v", u, username)
			}
			if password != p {
				t.Fatalf("unexpected c8y auth-password: %v expected: %v", p, password)
			}

			expectedReqURL := fmt.Sprintf("%v?%v", managedObjectPath, tt.filter.QueryParams())
			if reqURL != expectedReqURL {
				t.Fatalf("unexpected request url: %v expected: %v", reqURL, expectedReqURL)
			}

			if err == nil {
				if len(managedObjectCollection.ManagedObjects) != 3 {
					t.Errorf("unexpected count of collections. Expected: %v, actual: %v", 3, len(managedObjectCollection.ManagedObjects))
				}
			}
		})
	}
}
