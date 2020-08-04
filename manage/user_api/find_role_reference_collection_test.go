package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserApi_FindRoleReferenceCollection_for_user(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/%v/users/%v/roles?pageSize=5", tenantID, username)
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
	referenceCollection, err := api.FindRoleReferenceCollection(tenantID, username, "", 5)

	// then
	if err != nil {
		t.Fatalf("received an unexpected error: %s", err)
	}

	if len(referenceCollection.References) == 0 {
		t.Error("role reference collection does not contain any references. It should contain two references!")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}

	if !reflect.DeepEqual(roleReferenceCollection, referenceCollection) {
		t.Errorf("FindRoleReferenceCollection() want: %v, actual: %v", roleReferenceCollection, referenceCollection)
	}
}

func TestUserApi_FindRoleReferenceCollection_for_group(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/%v/groups/%v/roles?pageSize=5", tenantID, groupID)
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
	referenceCollection, err := api.FindRoleReferenceCollection(tenantID, "", groupID, 5)

	// then
	if err != nil {
		t.Fatalf("received an unexpected error: %s", err)
	}

	if len(referenceCollection.References) == 0 {
		t.Error("role reference collection does not contain any references. It should contain two references!")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}

	if !reflect.DeepEqual(roleReferenceCollection, referenceCollection) {
		t.Errorf("FindRoleReferenceCollection() want: %v, actual: %v", roleReferenceCollection, referenceCollection)
	}
}

func TestUserApi_FindRoleReferenceCollection_WithoutUsernameOrGroupID(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(roleReferenceCollectionJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.FindRoleReferenceCollection(tenantID, "", "", 5)

	// then
	expectedErr := generic.ClientError("Getting role reference collection without username or groupID is not allowed", "FindRoleReferenceCollection")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}

func TestUserApi_FindRoleReferenceCollection_invalid_pageSize(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(roleReferenceCollectionJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.FindRoleReferenceCollection(tenantID, "", groupID, -1)

	// then
	expectedErr := generic.ClientError("Error while building pageSize parameter to fetch reference collection: The page size must be between 1 and 2000. Was -1", "FindRoleReferenceCollection")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}

func TestUserApi_FindRoleReferenceCollection_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.FindRoleReferenceCollection(tenantID, username, "", 5)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}


