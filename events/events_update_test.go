package gomulocity_event

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func updateEventHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		var event UpdateEvent
		_ = json.Unmarshal(body, &event)
		eventUpdateCapture = &event
		updateUrlCapture = r.URL.Path
		requestCapture = r

		w.WriteHeader(status)
		response, _ := json.Marshal(responseEvent)
		_, _ = w.Write(response)
	}))
}

var eventUpdateCapture *UpdateEvent
var updateUrlCapture string

// given: A create event
var eventUpdate = &UpdateEvent{
	Text: "This is my new test event",
}

func TestEvents_Update_Event_Success_SendsData(t *testing.T) {
	// given: A test server
	ts := updateEventHttpServer(200)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	_, err := api.UpdateEvent(eventId, eventUpdate)

	if err != nil {
		t.Fatalf("UpdateEvent() got an unexpected error: %s", err.Error())
	}

	if strings.Contains(updateUrlCapture, eventId) == false {
		t.Errorf("UpdateEvent() The target URL does not contains the event Id: url: %s - expected eventId %s", updateUrlCapture, eventId)
	}

	if eventUpdateCapture == nil {
		t.Fatalf("UpdateEvent() Captured event is nil.")
	}

	if !reflect.DeepEqual(eventUpdate, eventUpdateCapture) {
		t.Errorf("UpdateEvent() event = %v, want %v", eventUpdate, eventUpdateCapture)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.eventApi+json"
	if header != want {
		t.Errorf("UpdateEvent() accent header = %v, want %v", header, want)
	}
}

func TestEvents_Update_Event_Success_ReceivesData(t *testing.T) {
	// given: A test server
	ts := updateEventHttpServer(200)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)
	event, err := api.UpdateEvent(eventId, eventUpdate)

	if err != nil {
		t.Fatalf("UpdateEvent() got an unexpected error: %s", err.Error())
	}

	if !reflect.DeepEqual(event, responseEvent) {
		t.Errorf("UpdateEvent() event = %v, want %v", createEvent, createEventCapture)
	}
}

func TestEvents_Update_Event_BadRequest(t *testing.T) {
	// given: A test server
	ts := createEventHttpServer(400)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)
	_, err := api.UpdateEvent(eventId, eventUpdate)

	if err == nil {
		t.Errorf("UpdateEvent() expected error on 400 - bad request")
		return
	}
}
