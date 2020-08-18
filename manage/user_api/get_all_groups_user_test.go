package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserApi_GetAllGroupsOfUser(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/%v/users/%v/groups?pageSize=%v", tenantID, username, 5)
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(groupReferenceCollectionJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	groupReferenceCollection, err := api.GetAllGroupsOfUser(tenantID, username, 5)

	// then
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err)
	}

	if groupReferenceCollection != nil {
		if len(groupReferenceCollection.Groups) == 0 {
			t.Error("GroupReferenceCollection does not contain any groups. It should contain three groups!")
		}
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}

	if !reflect.DeepEqual(testGroupReferenceCollection, groupReferenceCollection) {
		t.Errorf("GetAllGroupsOfUser() want: %v, actual: %v", testGroupReferenceCollection, groupReferenceCollection)
	}
}

func TestUserApi_GetAllGroupsOfUser_without_username(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GetAllGroupsOfUser(tenantID, "", 5)

	// then
	expectedErr := generic.ClientError("Getting a group reference collection without tenantID and username is not allowed", "FindGroupReferenceCollection")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_GetAllGroupsOfUser_invalid_pageSize(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GetAllGroupsOfUser(tenantID, username, -1)

	// then
	expectedErr := generic.ClientError("Error while building pageSize parameter to fetch group references: The page size must be between 1 and 2000. Was -1", "FindGroupReferenceCollection")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_GetAllGroupsOfUser_parse_response_empty_body(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GetAllGroupsOfUser(tenantID, username, 5)

	// then
	expectedErr := generic.ClientError("Response body was empty", "CollectionResponseParser")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_GetAllGroupsOfUser_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GetAllGroupsOfUser(tenantID, username, 5)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}
