package events

import (
	"reflect"
	"strings"
	"testing"
)

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

	if strings.Contains(requestCapture.URL.Path, eventId) == false {
		t.Errorf("UpdateEvent() The target URL does not contains the event Id: url: %s - expected eventId %s", updateUrlCapture, eventId)
	}

	if updateEventCapture == nil {
		t.Fatalf("UpdateEvent() Captured event is nil.")
	}

	if !reflect.DeepEqual(eventUpdate, updateEventCapture) {
		t.Errorf("UpdateEvent() event = %v, want %v", eventUpdate, updateEventCapture)
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
