package gomulocity_event

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	url "net/url"
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
