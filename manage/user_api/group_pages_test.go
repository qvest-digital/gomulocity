package user_api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserApi__Group_NextPage_Success(t *testing.T) {
	// given: An Http server with a next collection with one group.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(groupCollectionTemplate, groupJSON)))
	}))
	defer ts.Close()

	// and: the system under test
	api := buildUserApi(ts.URL)

	// when: We create an existing collection and call `NextPage`
	nextPageUrl := ts.URL + "/user/1111111/users/msmith/groups?pageSize=5&currentPage=3"
	collection := createGroupCollection(nextPageUrl, "")
	nextCollection, _ := api.NextPageGroupReferenceCollection(collection)

	// then: We got the next collection with one group.
	if capturedUrl != nextPageUrl {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, nextPageUrl)
	}

	if nextCollection == nil {
		t.Fatalf("NextPage() nextCollection is nil")
	}

	if len(nextCollection.Groups) != 1 {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, nextPageUrl)
	}

	group := nextCollection.Groups[0]
	if group.ID != groupID {
		t.Errorf("NextPage() next group id = %v, expected %v", group.ID, groupID)
	}
}

func TestUserApi_Group_NextPage_NotAvailable(t *testing.T) {
	// given: The system under test
	api := buildUserApi("https://does.not.exist")

	// when: We call `NextPage` with no URLs
	collection := createGroupCollection("", "")
	nextCollection, err := api.NextPageGroupReferenceCollection(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: No `nextCollection` is available.
	if nextCollection != nil {
		t.Errorf("NextPage() should return nil. Was: %v", nextCollection)
	}
}

func TestUserApi_Group_NextPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(groupCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `NextPage` with a given URL
	collection := createGroupCollection(ts.URL+"/user/1111111/users/msmith/groups?pageSize=5&currentPage=3", "")
	nextCollection, err := api.NextPageGroupReferenceCollection(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: `nextCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("NextPage() should return an empty collection on empty collection response.")
	}
}

func TestUserApi_Group_NextPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `NextPage` with a given URL
	collection := createGroupCollection(ts.URL+"/user/1111111/users/msmith/groups?pageSize=5&currentPage=3", "")
	_, err := api.NextPageGroupReferenceCollection(collection)

	// then: an error occurred
	if err == nil {
		t.Errorf("NextPage() should return error. Nil was given.")
	}
}

func TestUserApi_Group_PreviousPage_Success(t *testing.T) {
	// given: A Http server with a previous collection with one group.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(groupCollectionTemplate, groupJSON)))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We create an existing collection and call `PreviousPage`
	previousPageUrl := ts.URL + "/user/1111111/users/msmith/groups?pageSize=5&currentPage=1"
	collection := createGroupCollection("", previousPageUrl)
	nextCollection, _ := api.PreviousPageGroupCollection(collection)

	// then: We got the previous collection with one group.
	if capturedUrl != previousPageUrl {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, previousPageUrl)
	}

	if nextCollection == nil {
		t.Fatalf("PreviousPage() nextCollection is nil")
	}

	if len(nextCollection.Groups) != 1 {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, previousPageUrl)
	}

	group := nextCollection.Groups[0]
	if group.ID != groupID {
		t.Errorf("PreviousPage() next group id = %v, expected %v", group.ID, groupID)
	}
}

func TestUserApi_Group_PreviousPage_NotAvailable(t *testing.T) {
	// given: The system under test
	api := buildUserApi("https://does.not.exist")

	// when: We call `PreviousPage` with no URLs
	collection := createGroupCollection("", "")
	nextCollection, err := api.PreviousPageGroupCollection(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: No `previousCollection` is available.
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return nil. Was: %v", nextCollection)
	}
}

func TestUserApi_Group_PreviousPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(groupCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `PreviousPage` with a given URL
	collection := createGroupCollection("", ts.URL+"/user/roles?pageSize=5&currentPage=1")
	nextCollection, err := api.PreviousPageGroupCollection(collection)

	if err != nil {
		t.Errorf("PreviousPage() should not return an error. Was: %v", err)
	}

	// then: `previousCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return an empty collection on empty collection response.")
	}
}

func TestUserApi_Group_PreviousPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `PreviousPage` with a given URL
	collection := createGroupCollection("", ts.URL+"/user/roles?pageSize=5&currentPage=1")
	_, error := api.PreviousPageGroupCollection(collection)

	// then: an error occurred
	if error == nil {
		t.Errorf("PreviousPage() should return error. Nil was given.")
	}
}

func createGroupCollection(next string, prev string) *GroupReferenceCollection {
	return &GroupReferenceCollection{
		Next:       next,
		Self:       "https://t200588189.cumulocity.com/user/1111111/users/msmith/groups?pageSize=5&currentPage=2",
		Prev:       prev,
		Groups:     []Group{},
		Statistics: nil,
	}
}
