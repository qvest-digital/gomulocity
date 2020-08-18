package user_api

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserApi_InventoryRoleCollection(t *testing.T) {
	var expectedUrl = "/user/inventoryroles"
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testInventoryRoleCollectionJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	inventoryRoleCollection, err := api.InventoryRoleCollection()

	// then
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err)
	}

	if len(inventoryRoleCollection.Roles) == 0 {
		t.Errorf("collection does not contain any roles. It should contain two roles")
	}

	if !reflect.DeepEqual(testInventoryRoleCollection, inventoryRoleCollection) {
		t.Errorf("GetAllRolesOfAUser() want: %v, actual: %v", testInventoryRoleCollection, inventoryRoleCollection)
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}
}

func TestUserApi_InventoryRoleCollection_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.InventoryRoleCollection()

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_InventoryRoleCollection_error_unmarshalling(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<>"))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.InventoryRoleCollection()

	// then
	expectedErr := generic.ClientError("Error while parsing response JSON: Error while unmarshalling json: invalid character '<' looking for beginning of value", "CollectionResponseParser")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}
