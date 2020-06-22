package gomulocity_event

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var DeviceId = "1111111"
var EventId = "2222222"
var event = `{
	"creationTime": "2020-01-01T01:00:10.000Z",
	"source": {
		"name": "test-device",
		"self": "https://t0818.cumulocity.com/inventory/managedObjects/1111111",
		"id": "1111111"
	},
	"type": "threshold",
	"self": "https://t0818.cumulocity.com/event/events/2222222",
	"time": "2020-01-01T01:00:00.000Z",
	"text": "over 21°C",
	"id": "2222222"
}`
var eventCollectionTemplate = `{
    "next": "https://t0818.cumulocity.com/event/events?source=1111111&pageSize=5&currentPage=2",
    "self": "https://t0815.cumulocity.com/event/events?source=1111111&pageSize=5&currentPage=1",
    "events": [%s],
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
    }
}`

func TestEvents_GetForDevice(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()

		if queries["source"][0] == DeviceId {
			_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, event)))
		} else {
			_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, "")))
		}
	}))
	defer ts.Close()

	// and: a configured http client
	httpClient := http.DefaultClient
	client := Client{
		HTTPClient: httpClient,
		BaseURL:    ts.URL,
		Username:   "foo",
		Password:   "bar",
	}

	// and: the api as system under test
	api := NewEventsApi(client)

	tests := []struct {
		name                string
		givenDeviceId       string // aka source
		expectedEventsCount int
		expectedEventId     string
	}{
		{"existing device id", DeviceId, 1, EventId},
		{"non existing device id", "4711", 0, "-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection, err := api.GetForDevice(tt.givenDeviceId)

			if err != nil {
				t.Errorf("GetForDevice() got an unexpected error: %s", err.Error())
				return
			}

			if collection == nil {
				t.Error("GetForDevice() got no explict error but the collection was nil.")
				return
			}

			if len(collection.Events) != tt.expectedEventsCount {
				t.Errorf("GetForDevice() events count = %v, want %v", len(collection.Events), tt.expectedEventsCount)
			}

			if tt.expectedEventsCount == 1 {
				event := collection.Events[0]
				if event.Id != tt.expectedEventId {
					t.Errorf("GetForDevice() event id = %v, want %v", len(event.Id), tt.expectedEventId)
				}
			}
		})
	}
}
