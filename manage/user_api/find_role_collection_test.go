package user_api

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserApi_FindRoleCollection(t *testing.T) {
	var expectedUrl = "/user/roles?pageSize=5"
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(roleCollectionJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	collection, err := api.RoleCollection(5)
	if err != nil {
		t.Fatalf("received an unexpected error: %s", err)
	}

	if len(collection.Roles) == 0 {
		t.Error("role collection does not contain any roles. It should contain three roles!")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}

	if !reflect.DeepEqual(roleCollection, collection) {
		t.Errorf("FindRoleCollection() want: %v, actual: %v", roleCollection, collection)
	}
}

func TestUserApi_FindRoleCollection_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.RoleCollection(5)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}

func TestUserApi_FindRoleCollection_invalid_pageSize(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.RoleCollection(-1)

	// then
	expectedErr := generic.ClientError("Error while building pageSize parameter to fetch role collection: The page size must be between 1 and 2000. Was -1", "FindRoleCollection")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}
