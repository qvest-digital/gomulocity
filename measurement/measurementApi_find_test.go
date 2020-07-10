package measurement

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	url "net/url"
	"strings"
	"testing"
	"time"
)

func TestMeasurementApi_FindWithFilter(t *testing.T) {
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(measurementCollectionTemplate, measurement)))
	}))
	defer ts.Close()

	dateFrom, _ := time.Parse(time.RFC3339, "2020-06-01T01:00:00.00Z")
	dateTo, _ := time.Parse(time.RFC3339, "2020-06-30T01:00:00.00Z")

	tests := []struct {
		name          string
		query         MeasurementQuery
		expectedQuery string
	}{
		{
			"EmptyFilter",
			MeasurementQuery{},
			"pageSize=1",
		},
		{
			"DateAndRevertFlag",
			MeasurementQuery{DateFrom: &dateFrom, DateTo: &dateTo, Revert: true},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z&pageSize=1&revert=true",
		},
		{
			"SourceId",
			MeasurementQuery{SourceId: "123"},
			"pageSize=1&source=123",
		},
		{
			"Type",
			MeasurementQuery{Type: "testMeasurement"},
			"pageSize=1&type=testMeasurement",
		},
		{
			"ValueFragmentType",
			MeasurementQuery{ValueFragmentType: "fragmentName"},
			"pageSize=1&valueFragmentType=fragmentName",
		},
		{
			"ValueFragmentSeries",
			MeasurementQuery{ValueFragmentSeries: "serieName"},
			"pageSize=1&valueFragmentSeries=serieName",
		},
		{
			"All",
			MeasurementQuery{
				DateFrom:            &dateFrom,
				DateTo:              &dateTo,
				Type:                "testMeasurement",
				ValueFragmentType:   "fragmentName",
				ValueFragmentSeries: "seriesName",
				SourceId:            "123",
				Revert:              true,
			},
			"dateFrom=2020-06-01T01%3A00%3A00Z&dateTo=2020-06-30T01%3A00%3A00Z&pageSize=1&revert=true&source=123&type=testMeasurement&valueFragmentSeries=seriesName&valueFragmentType=fragmentName",
		},
	}

	api := buildMeasurementApi(ts.URL)

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

func TestMeasurementApi_Find_WithInvalidFilter(t *testing.T) {
	api := buildMeasurementApi("test.url")

	_, err := api.Find(&MeasurementQuery{Revert: true}, 1)

	expectedError := "if 'Revert' parameter is set to true, 'DateFrom' and 'DateTo' should be set as well"
	if !strings.Contains(err.Message, expectedError) {
		t.Errorf("Error in Find(): [%v], expected: [%v]", err.Message, expectedError)
	}
}

func TestMeasurementApi_Find_HandlesPageSize(t *testing.T) {
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
		_, _ = w.Write([]byte(fmt.Sprintf(measurementCollectionTemplate, measurement)))
	}))
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := MeasurementQuery{
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

func TestMeasurementApi_ReturnsCollection(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(measurementCollectionTemplate, measurement)))
	}))
	defer ts.Close()

	api := buildMeasurementApi(ts.URL)

	collection, err := api.Find(&MeasurementQuery{}, 1)

	if err != nil {
		t.Fatalf("Find() - Error given but no expected")
	}

	if len(collection.Measurements) != 1 {
		t.Fatalf("Find() = Collection size = %v, want %v", len(collection.Measurements), 1)
	}

	measurement := collection.Measurements[0]
	if measurement.Id != measurementId {
		t.Fatalf("Find() = Collection measurement id = %v, want %v", measurement.Id, measurementId)
	}
}

func TestMeasurementApi_FindReturnsError(t *testing.T) {
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

	api := buildMeasurementApi(ts.URL)

	_, err := api.Find(&MeasurementQuery{}, 1)

	if err == nil {
		t.Fatalf("Find() - Error expected")
	}

	if err.ErrorType != "400: undefined/validationError" {
		t.Errorf("Find() = '%v', want '%v'", err.ErrorType, "undefined/validationError")
	}

	if err.Message != "My fancy error" {
		t.Errorf("Find() = '%v', want '%v'", err.Message, "My fancy error")
	}

	if err.Info != "https://www.cumulocity.com/guides/reference-guide/#error_reporting" {
		t.Errorf("Find() = '%v', want '%v'", err.Info, "https://www.cumulocity.com/guides/reference-guide/#error_reporting")
	}
}
