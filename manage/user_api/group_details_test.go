package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserApi_GroupDetails(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/management/groups/%v", groupID)
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(groupJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	group, err := api.GroupDetails(groupID)

	// then
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err)
	}

	if group != nil {
		if len(group.Roles) == 0 {
			t.Error("Group does not contain any roles. It should contain three roles!")
		}

		if group.ID != groupID {
			t.Errorf("GroupIDs do not match! Expected ID: %v, actual: %v", groupID, group.ID)
		}
	} else {
		t.Error("GroupDetails() response group must not be nil")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}

	if !reflect.DeepEqual(testGroup, group) {
		t.Errorf("GroupDetails() want: %v, actual: %v", testGroup, group)
	}
}

func TestUserApi_GroupDetails_error_unmarshalling(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<>"))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GroupDetails(groupID)

	// then
	expectedErr := generic.ClientError("Error while unmarshalling response: invalid character '<' looking for beginning of value", "GroupDetails")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_GroupDetails_without_groupID(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(groupJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GroupDetails("")

	// then
	expectedErr := generic.ClientError("Getting group details without groupID is not allowed", "GroupDetails")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_GroupDetails_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GroupDetails(groupID)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}
