package gomulocity_event

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEvents_NextPage_Success(t *testing.T) {
	// given: A Http server with a next collection with one event.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, event)))
	}))
	defer ts.Close()

	// and: the system under test
	api := buildEventsApi(ts.URL)

	// when: We create an existing collection and call `NextPage`
	collection := createCollection(ts.URL+"/event/events?source=1111111&pageSize=5&currentPage=3", "")
	expectedUrl := ts.URL + "/event/events?source=1111111&pageSize=5&currentPage=3"
	nextCollection, _ := api.NextPage(collection)

	// then: We got the next collection with one event.
	if capturedUrl != expectedUrl {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, expectedUrl)
	}

	if nextCollection == nil {
		t.Fatalf("NextPage() nextCollection is nil")
	}

	if len(nextCollection.Events) != 1 {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, expectedUrl)
	}

	event := nextCollection.Events[0]
	if event.Id != eventId {
		t.Errorf("NextPage() next event id = %v, expected %v", event.Id, eventId)
	}
}

func TestEvents_NextPage_NotAvailable(t *testing.T) {
	// given: The system under test
	api := buildEventsApi("https://does.not.exist")

	// when: We call `NextPage` with no URLs
	collection := createCollection("", "")
	nextCollection, _ := api.NextPage(collection)

	// then: No `nextCollection` is available.
	if nextCollection != nil {
		t.Errorf("NextPage() should return nil. Was: %v", nextCollection)
	}
}

func TestEvents_NextPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildEventsApi(ts.URL)

	// when: We call `NextPage` with a given URL
	collection := createCollection(ts.URL+"/event/events?source=1111111&pageSize=5&currentPage=3", "")
	nextCollection, _ := api.NextPage(collection)

	// then: `nextCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("NextPage() should return an empty collection on empty collection response.")
	}
}

func TestEvents_NextPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildEventsApi(ts.URL)

	// when: We call `NextPage` with a given URL
	collection := createCollection(ts.URL+"/event/events?source=1111111&pageSize=5&currentPage=3", "")
	_, err := api.NextPage(collection)

	// then: an error occurred
	if err == nil {
		t.Errorf("NextPage() should return error. Nil was given.")
	}
}

func TestEvents_PreviousPage_Success(t *testing.T) {
	// given: A Http server with a previous collection with one event.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, event)))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildEventsApi(ts.URL)

	// when: We create an existing collection and call `PreviousPage`
	collection := createCollection("", ts.URL+"/event/events?source=1111111&pageSize=5&currentPage=1")
	expectedUrl := ts.URL + "/event/events?source=1111111&pageSize=5&currentPage=1"
	nextCollection, _ := api.PreviousPage(collection)

	// then: We got the previous collection with one event.
	if capturedUrl != expectedUrl {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, expectedUrl)
	}

	if nextCollection == nil {
		t.Fatalf("PreviousPage() nextCollection is nil")
	}

	if len(nextCollection.Events) != 1 {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, expectedUrl)
	}

	event := nextCollection.Events[0]
	if event.Id != eventId {
		t.Errorf("PreviousPage() next event id = %v, expected %v", event.Id, eventId)
	}
}

func TestEvents_PreviousPage_NotAvailable(t *testing.T) {
	// given: The system under test
	api := buildEventsApi("https://does.not.exist")

	// when: We call `PreviousPage` with no URLs
	collection := createCollection("", "")
	nextCollection, _ := api.PreviousPage(collection)

	// then: No `previousCollection` is available.
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return nil. Was: %v", nextCollection)
	}
}

func TestEvents_PreviousPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildEventsApi(ts.URL)

	// when: We call `PreviousPage` with a given URL
	collection := createCollection(ts.URL+"/event/events?source=1111111&pageSize=5&currentPage=3", "")
	nextCollection, _ := api.NextPage(collection)

	// then: `previousCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return an empty collection on empty collection response.")
	}
}

func TestEvents_PreviousPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildEventsApi(ts.URL)

	// when: We call `PreviousPage` with a given URL
	collection := createCollection("", ts.URL+"/event/events?source=1111111&pageSize=5&currentPage=1")
	_, error := api.PreviousPage(collection)

	// then: an error occurred
	if error == nil {
		t.Errorf("PreviousPage() should return error. Nil was given.")
	}
}

func createCollection(next string, prev string) *EventCollection {
	return &EventCollection{
		Next:       next,
		Self:       "https://t0818.cumulocity.com/event/events?source=1111111&pageSize=5&currentPage=2",
		Prev:       prev,
		Events:     []Event{},
		Statistics: nil,
	}
}
