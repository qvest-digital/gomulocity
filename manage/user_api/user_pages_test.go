package user_api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserApi_User_NextPage_Success(t *testing.T) {
	// given: An Http server with a next collection with one user.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(userCollectionTemplate, testUserJSON)))
	}))
	defer ts.Close()

	// and: the system under test
	api := buildUserApi(ts.URL)

	// when: We create an existing collection and call `NextPage`
	nextPageUrl := ts.URL + "/user/1111111/users?username=mmark&pageSize=2&currentPage=3"
	collection := createUserCollection(nextPageUrl, "")
	nextCollection, _ := api.NextPageUserCollection(collection)

	// then: We got the next collection with one user.
	if capturedUrl != nextPageUrl {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, nextPageUrl)
	}

	if nextCollection == nil {
		t.Fatalf("NextPage() nextCollection is nil")
	}

	if len(nextCollection.Users) != 1 {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, nextPageUrl)
	}

	user := nextCollection.Users[0]
	if user.ID != userID {
		t.Errorf("NextPage() next user id = %v, expected %v", user.ID, userID)
	}
}

func TestUserApi_User_NextPage_NotAvailable(t *testing.T) {
	// given: The system under test
	api := buildUserApi("https://does.not.exist")

	// when: We call `NextPage` with no URLs
	collection := createUserCollection("", "")
	nextCollection, err := api.NextPageUserCollection(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: No `nextCollection` is available.
	if nextCollection != nil {
		t.Errorf("NextPage() should return nil. Was: %v", nextCollection)
	}
}

func TestUserApi_User_NextPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(userCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `NextPage` with a given URL
	collection := createUserCollection(ts.URL+"/user/1111111/users?username=mmark&pageSize=2&currentPage=3", "")
	nextCollection, err := api.NextPageUserCollection(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: `nextCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("NextPage() should return an empty collection on empty collection response.")
	}
}

func TestUserApi_User_NextPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `NextPage` with a given URL
	collection := createUserCollection(ts.URL+"/user/1111111/users?username=mmark&pageSize=2&currentPage=3", "")
	_, err := api.NextPageUserCollection(collection)

	// then: an error occurred
	if err == nil {
		t.Errorf("NextPage() should return error. Nil was given.")
	}
}

func TestUserApi_User_PreviousPage_Success(t *testing.T) {
	// given: A Http server with a previous collection with one user.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(userCollectionTemplate, testUserJSON)))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We create an existing collection and call `PreviousPage`
	previousPageUrl := ts.URL + "/user/1111111/users?username=mmark&pageSize=2&currentPage=1"
	collection := createUserCollection("", previousPageUrl)
	nextCollection, _ := api.PreviousPageUserCollection(collection)

	// then: We got the previous collection with one user.
	if capturedUrl != previousPageUrl {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, previousPageUrl)
	}

	if nextCollection == nil {
		t.Fatalf("PreviousPage() nextCollection is nil")
	}

	if len(nextCollection.Users) != 1 {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, previousPageUrl)
	}

	user := nextCollection.Users[0]
	if user.ID != userID {
		t.Errorf("PreviousPage() next user id = %v, expected %v", user.ID, userID)
	}
}

func TestUserApi_User_PreviousPage_NotAvailable(t *testing.T) {
	// given: The system under test
	api := buildUserApi("https://does.not.exist")

	// when: We call `PreviousPage` with no URLs
	collection := createUserCollection("", "")
	nextCollection, err := api.PreviousPageUserCollection(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: No `previousCollection` is available.
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return nil. Was: %v", nextCollection)
	}
}

func TestUserApi_User_PreviousPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(userCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `PreviousPage` with a given URL
	collection := createUserCollection("", ts.URL+"/user/1111111/users?username=mmark&pageSize=2&currentPage=1")
	nextCollection, err := api.PreviousPageUserCollection(collection)

	if err != nil {
		t.Errorf("PreviousPage() should not return an error. Was: %v", err)
	}

	// then: `previousCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return an empty collection on empty collection response.")
	}
}

func TestUserApi_User_PreviousPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildUserApi(ts.URL)

	// when: We call `PreviousPage` with a given URL
	collection := createUserCollection("", ts.URL+"/user/1111111/users?username=mmark&pageSize=2&currentPage=1")
	_, error := api.PreviousPageUserCollection(collection)

	// then: an error occurred
	if error == nil {
		t.Errorf("PreviousPage() should return error. Nil was given.")
	}
}

func createUserCollection(next string, prev string) *UserCollection {
	return &UserCollection{
		Next:       next,
		Self:       "https://t200588189.cumulocity.com/user/1111111/users?username=mmark&pageSize=2&currentPage=2",
		Prev:       prev,
		Users:      []User{},
		Statistics: nil,
	}
}
