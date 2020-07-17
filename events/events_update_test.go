package events

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// given: A create event
var eventUpdate = &UpdateEvent{
	Text:             "This is my new test event",
	AdditionalFields: map[string]interface{}{},
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

func TestEvents_Update_Event_CustomFields(t *testing.T) {
	// given: A test server
	ts := updateEventHttpServer(200)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	// and: The update event
	updateEvent := &UpdateEvent{
		Text: "This is my test event",
		AdditionalFields: map[string]interface{}{
			"Custom1": 4711,
			"Custom2": "Hello World",
		},
	}

	// when: We send the create event
	_, err := api.UpdateEvent(eventId, updateEvent)

	// then: No error is returned
	if err != nil {
		t.Fatalf("UpdateEvent() got an unexpected error: %s", err.Error())
	}

	// and: A body was captured
	if bodyCapture == nil {
		t.Fatalf("UpdateEvent() Captured request is nil.")
	}

	// and: The body is a json structure
	var bodyMap map[string]interface{}
	jErr := json.Unmarshal(*bodyCapture, &bodyMap)

	if jErr != nil {
		t.Fatalf("UpdateEvent() request body can not be parsed %v", err)
	}

	// and: The "Custom1" and "Custom2" field is flattened
	custom1, _ := bodyMap["Custom1"].(float64)
	custom2, _ := bodyMap["Custom2"].(string)
	if custom1 != 4711 || custom2 != "Hello World" {
		t.Errorf("UpdateEvent() additional fields - \nevent = %v \nwant [Custom1, Custom2]", jErr)
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
