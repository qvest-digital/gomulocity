package events

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	url "net/url"
	"strings"
	"testing"
	"time"
)

func TestEvents_FindWithFilter(t *testing.T) {
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		capturedUrl = r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, event)))
	}))
	defer ts.Close()

	dateFrom, _ := time.Parse(time.RFC3339, "2020-06-01T01:00:00.00Z")
	dateTo, _ := time.Parse(time.RFC3339, "2020-06-30T01:00:00.00Z")

	tests := []struct {
		name          string
		query         EventQuery
		expectedQuery string
	}{
		{
			"All",
			EventQuery{},
			"",
		},
		{
			"ForDateAndFragmentType",
			EventQuery{DateFrom: &dateFrom, DateTo: &dateTo, FragmentType: "FragmentType_1"},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z&fragmentType=FragmentType_1",
		},
		{
			"ForFragmentTypeAndType",
			EventQuery{FragmentType: "FragmentType_1", Type: "Type_1"},
			"fragmentType=FragmentType_1&type=Type_1",
		},
		{
			"ForSourceAndType",
			EventQuery{Source: "4711", Type: "Type_1"},
			"source=4711&type=Type_1",
		},
		{
			"ForTimeAndType",
			EventQuery{DateFrom: &dateFrom, DateTo: &dateTo, Type: "Type_1"},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z&type=Type_1",
		},
		{
			"ForDateAndFragmentTypeAndType",
			EventQuery{DateFrom: &dateFrom, DateTo: &dateTo, FragmentType: "FragmentType_1", Type: "Type_1"},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z&fragmentType=FragmentType_1&type=Type_1",
		},
		{
			"ForFragmentType",
			EventQuery{FragmentType: "FragmentType_1"},
			"fragmentType=FragmentType_1",
		},
		{
			"ForSource",
			EventQuery{Source: "4711"},
			"source=4711",
		},
		{
			"ForSourceAndTimeAndType",
			EventQuery{Source: "4711", DateFrom: &dateFrom, DateTo: &dateTo, Type: "Type_1"},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z&source=4711&type=Type_1",
		},
		{
			"ForTime",
			EventQuery{DateFrom: &dateFrom, DateTo: &dateTo},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z",
		},
		{
			"ForSourceAndDateAndFragmentTypeAndType",
			EventQuery{Source: "4711", DateFrom: &dateFrom, DateTo: &dateTo, FragmentType: "FragmentType_1", Type: "Type_1"},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z&fragmentType=FragmentType_1&source=4711&type=Type_1",
		},
		{
			"ForSourceAndFragmentTypeAndType",
			EventQuery{Source: "4711", FragmentType: "FragmentType_1", Type: "Type_1"},
			"fragmentType=FragmentType_1&source=4711&type=Type_1",
		},
		{
			"ForType",
			EventQuery{Type: "Type_1"},
			"type=Type_1",
		},
		{
			"ForSourceAndFragmentType",
			EventQuery{Source: "4711", FragmentType: "FragmentType_1"},
			"fragmentType=FragmentType_1&source=4711",
		},
		{
			"ForSourceAndTime",
			EventQuery{Source: "4711", DateFrom: &dateFrom, DateTo: &dateTo},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z&source=4711",
		},
	}

	api := buildEventsApi(ts.URL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api.Find(tt.query)
			cUrl, err := url.Parse(capturedUrl)

			if err != nil {
				t.Fatalf("Find() - The captured URL is invalid - URL: %s, error: %s", capturedUrl, err.Error())
			}

			if cUrl.RawQuery != tt.expectedQuery {
				t.Errorf("Find() = %v, want %v", cUrl.RawQuery, tt.expectedQuery)
			}
		})
	}
}

func TestEvents_Find_HandlesPageSize(t *testing.T) {
	tests := []struct {
		name        string
		pageSize    int
		errExpected bool
	}{
		{"Negative", -1, true},
		{"Zero", 0, false},
		{"Max", 2000, false},
		{"too large", 2001, true},
		{"in range", 10, false},
	}

	// given: A test server
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, event)))
	}))
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := EventQuery{
				Source:   deviceId,
				PageSize: tt.pageSize,
			}
			_, err := api.Find(query)

			if tt.errExpected {
				if err == nil {
					t.Error("GetForDevice() error expected but was nil")
				}
			}

			if !tt.errExpected {
				contains := strings.Contains(capturedUrl, fmt.Sprintf("pageSize=%d", tt.pageSize))

				if tt.pageSize != 0 && !contains {
					t.Errorf("Find() expected pageSize '%d' in url. '%s' given", tt.pageSize, capturedUrl)
				}

				if tt.pageSize == 0 && contains {
					t.Error("Find() expected no pageSize in url on value 0")
				}
			}
		})
	}
}

func TestEvents_ReturnsCollection(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, event)))
	}))
	defer ts.Close()

	api := buildEventsApi(ts.URL)

	collection, err := api.Find(EventQuery{})

	if err != nil {
		t.Fatalf("Find() - Error given but no expected")
	}

	if len(collection.Events) != 1 {
		t.Fatalf("Find() = Collection size = %v, want %v", len(collection.Events), 1)
	}

	event := collection.Events[0]
	if event.Id != eventId {
		t.Fatalf("Find() = Collection event id = %v, want %v", event.Id, eventId)
	}
}

func TestEvents_ReturnsCollection_CustomElements(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(eventCollectionTemplate, event)))
	}))
	defer ts.Close()

	api := buildEventsApi(ts.URL)

	collection, err := api.Find(EventQuery{})

	if err != nil {
		t.Fatalf("Find() - Error given but no expected")
	}

	if len(collection.Events) != 1 {
		t.Fatalf("Find() = Collection size = %v, want %v", len(collection.Events), 1)
	}

	event := collection.Events[0]
	if event.Id != eventId {
		t.Fatalf("Find() = Collection event id = %v, want %v", event.Id, eventId)
	}

	if len(event.AdditionalFields) != 2 {
		t.Fatalf("Find() AdditionalFields length = %d, want %d", len(event.AdditionalFields), 2)
	}

	custom1, ok1 := event.AdditionalFields["custom1"].(string)
	custom2, ok2 := event.AdditionalFields["custom2"].([]interface{})

	if !(ok1 && custom1 == "Hello") {
		t.Errorf("Find() custom1 = %v, want %v", custom1, "Hello")
	}
	if !(ok2 && custom2[0] == "Foo" && custom2[1] == "Bar") {
		t.Errorf("Find() custom2 = [%v, %v], want [%v, %v]", custom2[0], custom2[1], "Foo", "Bar")
	}
}

func TestEvents_FindReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		error := `{
			"error": "undefined/validationError",
			"message": "My fancy error",
			"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
		}`

		w.WriteHeader(400)
		_, _ = w.Write([]byte(error))
	}))
	defer ts.Close()

	api := buildEventsApi(ts.URL)

	_, err := api.Find(EventQuery{})

	if err == nil {
		t.Fatalf("Find() - Error expected")
	}

	if err.Message != "My fancy error" {
		t.Errorf("Find() = '%v', want '%v'", err.Message, "My fancy error")
	}
}
