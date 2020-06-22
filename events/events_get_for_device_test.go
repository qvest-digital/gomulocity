package gomulocity_event

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func buildEventsApi(url string) Events {
	httpClient := http.DefaultClient
	client := Client{
		HTTPClient: httpClient,
		BaseURL:    url,
		Username:   "foo",
		Password:   "bar",
	}
	return NewEventsApi(client)
}

func buildHttpServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

var deviceId = "1111111"
var eventId = "2222222"
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

func TestEvents_GetForDevice_ExistingId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, fmt.Sprintf(eventCollectionTemplate, event))
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	t.Run("existing device id", func(t *testing.T) {
		collection, err := api.GetForDevice(deviceId)

		if err != nil {
			t.Fatalf("GetForDevice() got an unexpected error: %s", err.Error())
		}

		if collection == nil {
			t.Fatalf("GetForDevice() got no explict error but the collection was nil.")
		}

		if len(collection.Events) != 1 {
			t.Fatalf("GetForDevice() events count = %v, want %v", len(collection.Events), 1)
		}

		event := collection.Events[0]
		if event.Id != eventId {
			t.Errorf("GetForDevice() event id = %v, want %v", event.Id, eventId)
		}
	})
}

func TestEvents_GetForDevice_NotExistingId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, fmt.Sprintf(eventCollectionTemplate, ""))
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	t.Run("non existing device id", func(t *testing.T) {
		collection, err := api.GetForDevice(deviceId)

		if err != nil {
			t.Fatalf("GetForDevice() got an unexpected error: %s", err.Error())
			return
		}

		if collection == nil {
			t.Fatalf("GetForDevice() got no explict error but the collection was nil.")
			return
		}

		if len(collection.Events) != 0 {
			t.Fatalf("GetForDevice() events count = %v, want %v", len(collection.Events), 0)
		}
	})
}

func TestEvents_GetForDevice_MalformedResponse(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, "{ foo ...")
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	t.Run("malformed json", func(t *testing.T) {
		_, err := api.GetForDevice(deviceId)

		if err == nil {
			t.Errorf("GetForDevice() Expected error - non given")
			return
		}
	})
}
