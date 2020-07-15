package inventory
//
//import (
//	"fmt"
//	"github.com/tarent/gomulocity/generic"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestInventoryApi_ReferenceCollection(t *testing.T) {
//	tests := []struct {
//		name        string
//		deviceID    string
//		reference   string
//		respCode    int
//		respBody    string
//		expectedErr *generic.Error
//	}{
//		{
//			name:      "Happy",
//			deviceID:  "104940",
//			reference: "childAdditions",
//			respCode:  http.StatusOK,
//			respBody:  ReferenceCollectionJson,
//		},
//		{
//			name:        "Happy - no collection found",
//			deviceID:    "104940",
//			reference:   "childAdditions",
//			respCode:    http.StatusNotFound,
//			expectedErr: clientError("no reference collection found for reference: childAdditions", "ReferenceCollection"),
//		},
//		{
//			name:      "Unhappy - status is not StatusOK or StatusNotFound",
//			deviceID:  "104940",
//			reference: "childAdditions",
//			respCode:  http.StatusBadRequest,
//			respBody:  `{"error": "inventory/Not Found","message": "Finding device data from database failed : No managedObject for id '354365346346'!","info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"}`,
//			expectedErr: &generic.Error{
//				ErrorType: "inventory/Not Found",
//				Message:   "Finding device data from database failed : No managedObject for id '354365346346'!",
//				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
//			},
//		},
//		{
//			name: "Unhappy - no deviceID and reference given",
//			expectedErr: clientError("given deviceID or reference is empty", "ReferenceCollection"),
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var reqURL string
//
//			testserver := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
//				reqURL = fmt.Sprintf("%v", request.URL.String())
//
//				response.WriteHeader(tt.respCode)
//				_, err := response.Write([]byte(tt.respBody))
//				if err != nil {
//					t.Fatalf("failed to write resp body: %s", err)
//				}
//			}))
//			defer testserver.Close()
//
//			c := inventoryApi{
//				Client: &generic.Client{
//					HTTPClient: testserver.Client(),
//					BaseURL:    testserver.URL,
//				},
//				ManagedObjectsPath: managedObjectPath,
//			}
//
//			referenceCollection, err := c.ReferenceCollection(tt.deviceID, tt.reference)
//			if err != nil && err.Error() != tt.expectedErr.Error() {
//				t.Errorf("received an unexpected error: expected: %v, actual: %v", tt.expectedErr, err)
//			}
//
//			if err == nil && tt.respCode != http.StatusNotFound && len(referenceCollection.References) == 0 {
//				t.Error("no references found")
//			}
//
//			if len(tt.deviceID) > 0 && len(tt.reference) > 0 {
//				expectedReqURL := fmt.Sprintf("%v/%v/%v", managedObjectPath, tt.deviceID, tt.reference)
//				if reqURL != expectedReqURL {
//					t.Fatalf("unexpected request url: %v expected: %v", reqURL, expectedReqURL)
//				}
//			}
//		})
//	}
//}
