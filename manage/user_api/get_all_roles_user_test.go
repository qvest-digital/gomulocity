package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserApi_GetAllRolesOfAUser(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/%v/users/%v/roles?pageSize=%v", tenantID, username, 5)
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(roleReferenceCollectionJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	ref, err := api.GetAllRolesOfAUser(tenantID, username, 5)

	// then
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err)
	}

	if ref != nil {
		if len(ref.References) == 0 {
			t.Error("role reference collection does not contain any references. It should contain two references!")
		}
	} else {
		t.Error("GetAllRolesOfAUser() role reference must not be nil")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}

	if !reflect.DeepEqual(roleReferenceCollection, ref) {
		t.Errorf("GetAllRolesOfAUser() want: %v, actual: %v", roleReferenceCollection, ref)
	}
}

func TestUserApi_GetAllRolesOfAUser_without_username(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GetAllRolesOfAUser(tenantID, "", 5)

	// then
	expectedErr := generic.ClientError("Getting role reference collection without username or groupID is not allowed", "FindRoleReferenceCollection")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_GetAllRolesOfAUser_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GetAllRolesOfAUser(tenantID, username, 5)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}
