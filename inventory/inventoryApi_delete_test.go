package inventory

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestInventoryApi_CommonPropertiesOnDelete(t *testing.T) {
	var reqBasicAuthUsername, reqBasicAuthPassword, reqURL, reqAccept, reqContentType, reqBody string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqURL = req.URL.String()
		reqBodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("failed to read c8y request body: %s", err)
		}
		defer req.Body.Close()
		reqBody = string(reqBodyBytes)

		reqBasicAuthUsername, reqBasicAuthPassword, _ = req.BasicAuth()
		reqAccept = req.Header.Get("Accept")
		reqContentType = req.Header.Get("Content-Type")
		res.WriteHeader(http.StatusNoContent)
		_, err = res.Write([]byte(""))
		if err != nil {
			t.Fatalf("failed to write resp body: %s", err)
		}
	}))
	defer testServer.Close()

	inventoryApi := buildInventoryApi(testServer)
	err := inventoryApi.Delete(managedObjectId)
	if err != nil {
		t.Errorf("received an unexpected error: %v", err)
	}

	if len(reqBody) > 0 {
		t.Errorf("processed an unexpected c8y request body %q", reqBody)
	}

	if len(reqAccept) > 0 {
		t.Errorf("unexpected request accept header:%q", reqAccept)
	}

	if len(reqContentType) > 0 {
		t.Errorf("unexpected request content-type header:%q", reqContentType)
	}

	if reqBasicAuthUsername != USER {
		t.Errorf("unexpected c8y request basic-auth username. Expected %q. Given: %q", USER, reqBasicAuthUsername)
	}
	if reqBasicAuthPassword != PASSWORD {
		t.Errorf("unexpected c8y request basic-auth password. Expected %q. Given: %q", PASSWORD, reqBasicAuthPassword)
	}

	var expectedC8YRequestURL = fmt.Sprintf("/inventory/managedObjects/%s", managedObjectId)
	if reqURL != expectedC8YRequestURL {
		t.Errorf("unexpected c8y request url. Expected %q. Given: %q", expectedC8YRequestURL, reqURL)
	}
}

func TestInventoryApi_DeleteManagedObjectById(t *testing.T) {
	tests := []struct {
		name        string
		requestedManagedObjectId string
		c8yRespCode int
		c8yRespBody string

		expectedErr *generic.Error
	}{
		{
			name:        "requested Id is empty",
			requestedManagedObjectId: "",
			c8yRespCode: http.StatusNoContent,
			c8yRespBody: "",
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Deleting managedObject without an id is not allowed",
				Info:      "DeleteManagedObject",
			},
		}, {
			name:        "bad credentials",
			requestedManagedObjectId: managedObjectId,
			c8yRespCode: http.StatusUnauthorized,
			c8yRespBody: `{
				"error": "security/Unauthorized",
				"message": "Invalid credentials! : Bad credentials",
				"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
			}`,
			expectedErr: &generic.Error{
				ErrorType: "401: security/Unauthorized",
				Message:   "Invalid credentials! : Bad credentials",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			},
		}, {
			name:        "error without status code on DELETE",
			requestedManagedObjectId: managedObjectId,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while deleting managedObject with id.*",
				Info:      "DeleteManagedObject",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(tt.c8yRespCode)
				_, err := res.Write([]byte(tt.c8yRespBody))
				if err != nil {
					t.Fatalf("failed to write resp body: %s", err)
				}
			}))
			defer testServer.Close()

			inventoryApi := buildInventoryApi(testServer)
			err := inventoryApi.Delete(tt.requestedManagedObjectId)

			if matched, _ := regexp.MatchString(fmt.Sprint(tt.expectedErr), fmt.Sprint(err)); !matched {
				t.Fatalf("received an unexpected error: %s\nExpected: %s", err, tt.expectedErr)
			}
		})
	}
}
