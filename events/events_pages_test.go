package gomulocity_event

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEvents_NextPage_Success(t *testing.T) {
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, event)))
	}))
	defer ts.Close()

	api := buildEventsApi(ts.URL)

	collection := createCollection("https://t0818.cumulocity.com/event/events?source=1111111&pageSize=5&currentPage=3", "")
	expectedUrl := ts.URL + "/event/events?source=1111111&pageSize=5&currentPage=3"
	nextCollection, _ := api.NextPage(collection)

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
	api := buildEventsApi("https://fake.net")

	collection := createCollection("", "")
	nextCollection, _ := api.NextPage(collection)

	if nextCollection != nil {
		t.Errorf("NextPage() should return nil. Was: %v", nextCollection)
	}
}

func TestEvents_NextPage_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	api := buildEventsApi(ts.URL)

	collection := createCollection("https://t0818.cumulocity.com/event/events?source=1111111&pageSize=5&currentPage=3", "")
	_, error := api.NextPage(collection)

	if error == nil {
		t.Errorf("NextPage() should return error. Nil was given.")
	}
}

func TestEvents_PreviousPage_Success(t *testing.T) {
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, event)))
	}))
	defer ts.Close()

	api := buildEventsApi(ts.URL)

	collection := createCollection("", "https://t0818.cumulocity.com/event/events?source=1111111&pageSize=5&currentPage=1")
	expectedUrl := ts.URL + "/event/events?source=1111111&pageSize=5&currentPage=1"
	nextCollection, _ := api.PreviousPage(collection)

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
	api := buildEventsApi("https://fake.net")

	collection := createCollection("", "")
	nextCollection, _ := api.PreviousPage(collection)

	if nextCollection != nil {
		t.Errorf("PreviousPage() should return nil. Was: %v", nextCollection)
	}
}

func TestEvents_PreviousPage_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	api := buildEventsApi(ts.URL)

	collection := createCollection("", "https://t0818.cumulocity.com/event/events?source=1111111&pageSize=5&currentPage=1")
	_, error := api.PreviousPage(collection)

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
