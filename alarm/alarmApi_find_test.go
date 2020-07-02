package alarm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	url "net/url"
	"strings"
	"testing"
	"time"
)

func TestAlarmApi_FindWithFilter(t *testing.T) {
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(alarmCollectionTemplate, alarm)))
	}))
	defer ts.Close()

	dateFrom, _ := time.Parse(time.RFC3339, "2020-06-01T01:00:00.00Z")
	dateTo, _ := time.Parse(time.RFC3339, "2020-06-30T01:00:00.00Z")

	tests := []struct {
		name          string
		query         AlarmFilter
		expectedQuery string
	}{
		{
			"EmptyFilter",
			AlarmFilter{},
			"pageSize=1",
		},
		{
			"Date",
			AlarmFilter{DateFrom: &dateFrom, DateTo: &dateTo},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z&pageSize=1",
		},
		{
			"Status",
			AlarmFilter{Status: []Status{ACTIVE, CLEARED, ACKNOWLEDGED}},
			"pageSize=1&status=ACTIVE%2CCLEARED%2CACKNOWLEDGED",
		},
		{
			"Severity",
			AlarmFilter{Severity: CRITICAL},
			"pageSize=1&severity=CRITICAL",
		},
		{
			"SourceId",
			AlarmFilter{SourceId: "123"},
			"pageSize=1&source=123",
		},
		{
			"Type",
			AlarmFilter{Type: "testAlarm"},
			"pageSize=1&type=testAlarm",
		},
		{
			"Resolved",
			AlarmFilter{Resolved: "true"},
			"pageSize=1&resolved=true",
		},
		{
			"WithSourceAssets",
			AlarmFilter{WithSourceAssets: true, SourceId: "123"},
			"pageSize=1&source=123&withSourceAssets=true",
		},
		{
			"WithSourceDevices",
			AlarmFilter{WithSourceDevices: true, SourceId: "123"},
			"pageSize=1&source=123&withSourceDevices=true",
		},
		{
			"All",
			AlarmFilter{
				Status:            []Status{ACKNOWLEDGED},
				SourceId:          "123",
				WithSourceAssets:  true,
				WithSourceDevices: true,
				Resolved:          "false",
				Severity:          MAJOR,
				DateFrom:          &dateFrom,
				DateTo:            &dateTo,
				Type:              "testAlarm",
			},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z&pageSize=1&resolved=false&severity=MAJOR&source=123&status=ACKNOWLEDGED&type=testAlarm&withSourceAssets=true&withSourceDevices=true",
		},
	}

	api := buildAlarmApi(ts.URL)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api.Find(&tt.query, 1)
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

func TestAlarmApi_Find_WithInvalidFilter(t *testing.T) {
	tests := []struct {
		name          string
		query         AlarmFilter
		expectedError string
	}{
		{
			"Resolved",
			AlarmFilter{Resolved: "CLEARED"},
			"if 'Resolved' parameter is set, only 'true' and 'false' values are accepted",
		},
		{
			"WithSourceAssets",
			AlarmFilter{WithSourceAssets: true},
			"when 'WithSourceAssets' parameter is defined also SourceID must be set",
		},
		{
			"WithSourceDevices",
			AlarmFilter{WithSourceDevices: true},
			"when 'WithSourceDevices' parameter is defined also SourceID must be set",
		},
	}

	api := buildAlarmApi("test.url")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := api.Find(&tt.query, 1)

			if strings.Contains(tt.expectedError, err.Message) {
				t.Errorf("Error in Find(): [%v], expected: [%v]", err.Message, tt.expectedError)
			}
		})
	}
}

func TestAlarmApi_Find_HandlesPageSize(t *testing.T) {
	tests := []struct {
		name        string
		pageSize    int
		errExpected bool
	}{
		{"Negative", -1, true},
		{"Zero", 0, true},
		{"Min", 1, false},
		{"Max", 2000, false},
		{"too large", 2001, true},
		{"in range", 10, false},
	}

	// given: A test server
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(alarmCollectionTemplate, alarm)))
	}))
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := AlarmFilter{
				SourceId:   deviceId,
			}
			_, err := api.Find(&query, tt.pageSize)

			if tt.errExpected {
				if err == nil {
					t.Error("GetForDevice() error expected but was nil")
				}
			}

			if !tt.errExpected {
				contains := strings.Contains(capturedUrl, fmt.Sprintf("pageSize=%d&source=%s", tt.pageSize, deviceId))

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

func TestAlarmApi_ReturnsCollection(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(alarmCollectionTemplate, alarm)))
	}))
	defer ts.Close()

	api := buildAlarmApi(ts.URL)

	collection, err := api.Find(&AlarmFilter{}, 1)

	if err != nil {
		t.Fatalf("Find() - Error given but no expected")
	}

	if len(collection.Alarms) != 1 {
		t.Fatalf("Find() = Collection size = %v, want %v", len(collection.Alarms), 1)
	}

	alarm := collection.Alarms[0]
	if alarm.Id != alarmId {
		t.Fatalf("Find() = Collection alarm id = %v, want %v", alarm.Id, alarmId)
	}
}

func TestAlarmApi_FindReturnsError(t *testing.T) {
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

	api := buildAlarmApi(ts.URL)

	_, err := api.Find(&AlarmFilter{}, 1)

	if err == nil {
		t.Fatalf("Find() - Error expected")
	}

	if err.ErrorType != "undefined/validationError" {
		t.Errorf("Find() = '%v', want '%v'", err.ErrorType, "undefined/validationError")
	}

	if err.Message != "My fancy error" {
		t.Errorf("Find() = '%v', want '%v'", err.Message, "My fancy error")
	}

	if err.Info != "https://www.cumulocity.com/guides/reference-guide/#error_reporting" {
		t.Errorf("Find() = '%v', want '%v'", err.Info, "https://www.cumulocity.com/guides/reference-guide/#error_reporting")
	}
}
