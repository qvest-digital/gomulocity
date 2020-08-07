package user_api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserApi_NextPage_Success(t *testing.T) {
	// given: An Http server with a next collection with one role.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(roleCollectionTemplate, roleJSON)))
	}))
	defer ts.Close()

	// and: the system under test
	api := buildUserApi(ts.URL)

	// when: We create an existing collection and call `NextPage`
	nextPageUrl := ts.URL + "/user/roles?pageSize=5&currentPage=3"
	collection := createRoleCollection(nextPageUrl, "")
	nextCollection, _ := api.NextPageRoleCollection(collection)

	// then: We got the next collection with one role.
	if capturedUrl != nextPageUrl {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, nextPageUrl)
	}

	if nextCollection == nil {
		t.Fatalf("NextPage() nextCollection is nil")
	}

	if len(nextCollection.Roles) != 1 {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, nextPageUrl)
	}

	role := nextCollection.Roles[0]
	if role.ID != roleID {
		t.Errorf("NextPage() next role id = %v, expected %v", role.ID, roleID)
	}
}

func TestUserApi_NextPage_NotAvailable(t *testing.T) {
	// given: The system under test
	api := buildUserApi("https://does.not.exist")

	// when: We call `NextPage` with no URLs
	collection := createRoleCollection("", "")
	nextCollection, err := api.NextPageRoleCollection(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: No `nextCollection` is available.
	if nextCollection != nil {
		t.Errorf("NextPage() should return nil. Was: %v", nextCollection)
	}
}

func TestUserApi_Role_NextPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(userCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `NextPage` with a given URL
	collection := createRoleCollection(ts.URL+"/user/roles?pageSize=5&currentPage=3", "")
	nextCollection, err := api.NextPageRoleCollection(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: `nextCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("NextPage() should return an empty collection on empty collection response.")
	}
}

func TestUserApi_Role_NextPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `NextPage` with a given URL
	collection := createRoleCollection(ts.URL+"/user/roles?pageSize=5&currentPage=3", "")
	_, err := api.NextPageRoleCollection(collection)

	// then: an error occurred
	if err == nil {
		t.Errorf("NextPage() should return error. Nil was given.")
	}
}

func TestUserApi_Role_PreviousPage_Success(t *testing.T) {
	// given: A Http server with a previous collection with one role.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(roleCollectionTemplate, roleJSON)))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We create an existing collection and call `PreviousPage`
	previousPageUrl := ts.URL + "/user/roles?pageSize=5&currentPage=1"
	collection := createRoleCollection("", previousPageUrl)
	nextCollection, _ := api.PreviousPageRoleCollection(collection)

	// then: We got the previous collection with one role.
	if capturedUrl != previousPageUrl {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, previousPageUrl)
	}

	if nextCollection == nil {
		t.Fatalf("PreviousPage() nextCollection is nil")
	}

	if len(nextCollection.Roles) != 1 {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, previousPageUrl)
	}

	role := nextCollection.Roles[0]
	if role.ID != roleID {
		t.Errorf("PreviousPage() next role id = %v, expected %v", role.ID, roleID)
	}
}

func TestUserApi_Role_PreviousPage_NotAvailable(t *testing.T) {
	// given: The system under test
	api := buildUserApi("https://does.not.exist")

	// when: We call `PreviousPage` with no URLs
	collection := createRoleCollection("", "")
	nextCollection, err := api.PreviousPageRoleCollection(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: No `previousCollection` is available.
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return nil. Was: %v", nextCollection)
	}
}

func TestUserApi_Role_PreviousPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(userCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `PreviousPage` with a given URL
	collection := createRoleCollection("", ts.URL+"/user/roles?pageSize=5&currentPage=1")
	nextCollection, err := api.PreviousPageRoleCollection(collection)

	if err != nil {
		t.Errorf("PreviousPage() should not return an error. Was: %v", err)
	}

	// then: `previousCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return an empty collection on empty collection response.")
	}
}

func TestUserApi_Role_PreviousPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `PreviousPage` with a given URL
	collection := createRoleCollection("", ts.URL+"/user/roles?pageSize=5&currentPage=1")
	_, error := api.PreviousPageRoleCollection(collection)

	// then: an error occurred
	if error == nil {
		t.Errorf("PreviousPage() should return error. Nil was given.")
	}
}

func createRoleCollection(next string, prev string) *RoleCollection {
	return &RoleCollection{
		Next:       next,
		Self:       "https://t200588189.cumulocity.com/user/roles?pageSize=5&currentPage=2",
		Prev:       prev,
		Roles:      []Role{},
		Statistics: nil,
	}
}
