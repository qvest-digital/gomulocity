package inventory

import (
	"fmt"
	jsoncompare "github.com/orasik/gocomparejson"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
)

func TestInventoryApi_CommonPropertiesOnCreate(t *testing.T) {
	var expectedContentType = "application/vnd.com.nsn.cumulocity.managedObject+json"
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
		res.WriteHeader(http.StatusCreated)
		_, err = res.Write([]byte(givenResponseBody))
		if err != nil {
			t.Fatalf("failed to write resp body: %s", err)
		}
	}))
	defer testServer.Close()

	inventoryApi := buildInventoryApi(testServer)
	managedObject, err := inventoryApi.Create(newManagedObject)
	if err != nil {
		t.Errorf("received an unexpected error: %v", err)
	}

	if !reflect.DeepEqual(managedObject, expectedManagedObject) {
		t.Errorf("received an unexpected ManagedObject: %#v. \nExpected: %#v", managedObject, expectedManagedObject)
	}

	if equal, _ := jsoncompare.CompareJSON(reqBody, expectedRequestBody); !equal {
		t.Errorf("processed an unexpected c8y request body %q\nExpected: %q", reqBody, expectedRequestBody)
	}

	if reqAccept != expectedContentType {
		t.Errorf("unexpected request accept header. Expected %q. Given: %q", expectedContentType, reqAccept)
	}
	if reqContentType != expectedContentType {
		t.Errorf("unexpected request content-type header. Expected %q. Given: %q", expectedContentType, reqContentType)
	}

	if reqBasicAuthUsername != USER {
		t.Errorf("unexpected c8y request basic-auth username. Expected %q. Given: %q", USER, reqBasicAuthUsername)
	}
	if reqBasicAuthPassword != PASSWORD {
		t.Errorf("unexpected c8y request basic-auth password. Expected %q. Given: %q", PASSWORD, reqBasicAuthPassword)
	}

	var expectedC8YRequestURL = "/inventory/managedObjects"
	if reqURL != expectedC8YRequestURL {
		t.Errorf("unexpected c8y request url. Expected %q. Given: %q", expectedC8YRequestURL, reqURL)
	}

	custom, ok := managedObject.AdditionalFields["custom"].(string)
	if !ok {
		t.Error("additional fields do not contain 'custom'")
	}

	if custom != "hello" {
		t.Errorf("Received an unexpected value from additionalFields map. Expected: %v, actual: %v", "hello", custom)
	}
}

func TestInventoryApi_CreateManagedObject(t *testing.T) {
	tests := []struct {
		name        string
		requestData *NewManagedObject
		c8yRespCode int
		c8yRespBody string

		expectedErr *generic.Error
	}{
		{
			name:        "request without properties",
			requestData: &NewManagedObject{},
			c8yRespCode: http.StatusCreated,
			c8yRespBody: givenResponseBody,
			expectedErr: nil,
		}, {
			name:        "request is nil",
			requestData: nil,
			c8yRespCode: http.StatusCreated,
			c8yRespBody: givenResponseBody,
			expectedErr: nil,
		}, {
			name:        "bad credentials",
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
			name:        "error without status code on POST",
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while posting a new managedObject: Post.*",
				Info:      "CreateManagedObject",
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
			_, err := inventoryApi.Create(newManagedObject)

			if matched, _ := regexp.MatchString(fmt.Sprint(tt.expectedErr), fmt.Sprint(err)); !matched {
				t.Fatalf("received an unexpected error: %s\nExpected: %s", err, tt.expectedErr)
			}
		})
	}
}
