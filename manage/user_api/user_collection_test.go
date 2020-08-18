package user_api

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserApi_UserCollection(t *testing.T) {
	var expectedUrl = "/user?pageSize=2&username=mmark&groups=group1,group3"
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userCollectionJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// and: query filter
	filter := &QueryFilter{
		Username: "mmark",
		Groups: []Group{
			{
				ID:   "group1",
				Name: "group1",
			},
			{
				ID:   "group3",
				Name: "group3",
			},
		},
	}

	// when
	collection, err := api.UserCollection(filter, 2)

	// then
	if err != nil {
		t.Errorf("Unexpected error while getting user collection: %s", err)
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}

	if len(collection.Users) == 0 && err == nil {
		t.Error("user collection does not contain any user. It should contain two user!")
	}

	if !reflect.DeepEqual(testUserCollection, collection) {
		t.Errorf("UserCollection() want: %v, actual: %v", testUserCollection, collection)
	}
}

func TestUserApi_UserCollection_invalid_pageSize(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusOK, userCollectionJSON)
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.UserCollection(&QueryFilter{}, -1)

	// then
	expectedErr := generic.ClientError("Error while building pageSize parameter to fetch measurements: The page size must be between 1 and 2000. Was -1","FindUserCollection")
	if err.Error() != expectedErr.Error() {
		t.Errorf("Received an unexpected error. Expected: %s, actual: %s", expectedErr, err)
	}
}

func TestUserApi_UserCollection_invalid_status(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusInternalServerError, erroneousResponseJSON)
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.UserCollection(&QueryFilter{}, 2)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Errorf("Received an unexpected error. Expected: %s, actual: %s", expectedErr, err)
	}
}

func TestUserApi_UserCollection_without_filter(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusInternalServerError, erroneousResponseJSON)
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.UserCollection(nil, 2)

	// then
	expectedErr := generic.ClientError("Given filter is empty", "FindUserCollection")
	if err.Error() != expectedErr.Error() {
		t.Errorf("Received an unexpected error. Expected: %s, actual: %s", expectedErr, err)
	}
}

