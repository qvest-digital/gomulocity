package measurement

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMeasurementApi_NextPage_Success(t *testing.T) {
	// given: An Http server with a next collection with one measurement.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(measurementCollectionTemplate, measurement)))
	}))
	defer ts.Close()

	// and: the system under test
	api := buildMeasurementApi(ts.URL)

	// when: We create an existing collection and call `NextPage`
	nextPageUrl := ts.URL + "/measurement/measurements?source=1111111&pageSize=5&currentPage=3"
	collection := createCollection(nextPageUrl, "")
	nextCollection, _ := api.NextPage(collection)

	// then: We got the next collection with one measurement.
	if capturedUrl != nextPageUrl {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, nextPageUrl)
	}

	if nextCollection == nil {
		t.Fatalf("NextPage() nextCollection is nil")
	}

	if len(nextCollection.Measurements) != 1 {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, nextPageUrl)
	}

	measurement := nextCollection.Measurements[0]
	if measurement.Id != measurementId {
		t.Errorf("NextPage() next measurement id = %v, expected %v", measurement.Id, measurementId)
	}
}

func TestMeasurementApi_NextPage_NotAvailable(t *testing.T) {
	// given: The system under test
	api := buildMeasurementApi("https://does.not.exist")

	// when: We call `NextPage` with no URLs
	collection := createCollection("", "")
	nextCollection, err := api.NextPage(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: No `nextCollection` is available.
	if nextCollection != nil {
		t.Errorf("NextPage() should return nil. Was: %v", nextCollection)
	}
}

func TestMeasurementApi_NextPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(measurementCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildMeasurementApi(ts.URL)

	// when: We call `NextPage` with a given URL
	collection := createCollection(ts.URL+"/measurement/measurements?source=1111111&pageSize=5&currentPage=3", "")
	nextCollection, err := api.NextPage(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: `nextCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("NextPage() should return an empty collection on empty collection response.")
	}
}

func TestMeasurementApi_NextPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildMeasurementApi(ts.URL)

	// when: We call `NextPage` with a given URL
	collection := createCollection(ts.URL+"/measurement/measurements?source=1111111&pageSize=5&currentPage=3", "")
	_, err := api.NextPage(collection)

	// then: an error occurred
	if err == nil {
		t.Errorf("NextPage() should return error. Nil was given.")
	}
}

func TestMeasurementApi_PreviousPage_Success(t *testing.T) {
	// given: A Http server with a previous collection with one measurement.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(measurementCollectionTemplate, measurement)))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildMeasurementApi(ts.URL)

	// when: We create an existing collection and call `PreviousPage`
	previousPageUrl := ts.URL + "/measurement/measurements?source=1111111&pageSize=5&currentPage=1"
	collection := createCollection("", previousPageUrl)
	nextCollection, _ := api.PreviousPage(collection)

	// then: We got the previous collection with one measurement.
	if capturedUrl != previousPageUrl {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, previousPageUrl)
	}

	if nextCollection == nil {
		t.Fatalf("PreviousPage() nextCollection is nil")
	}

	if len(nextCollection.Measurements) != 1 {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, previousPageUrl)
	}

	measurement := nextCollection.Measurements[0]
	if measurement.Id != measurementId {
		t.Errorf("PreviousPage() next measurement id = %v, expected %v", measurement.Id, measurementId)
	}
}

func TestMeasurementApi_PreviousPage_NotAvailable(t *testing.T) {
	// given: The system under test
	api := buildMeasurementApi("https://does.not.exist")

	// when: We call `PreviousPage` with no URLs
	collection := createCollection("", "")
	nextCollection, err := api.PreviousPage(collection)

	if err != nil {
		t.Errorf("NextPage() should not return an error. Was: %v", err)
	}

	// then: No `previousCollection` is available.
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return nil. Was: %v", nextCollection)
	}
}

func TestMeasurementApi_PreviousPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(measurementCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildMeasurementApi(ts.URL)

	// when: We call `PreviousPage` with a given URL
	collection := createCollection("", ts.URL+"/measurement/measurements?source=1111111&pageSize=5&currentPage=1")
	nextCollection, err := api.PreviousPage(collection)

	if err != nil {
		t.Errorf("PreviousPage() should not return an error. Was: %v", err)
	}

	// then: `previousCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return an empty collection on empty collection response.")
	}
}

func TestMeasurementApi_PreviousPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildMeasurementApi(ts.URL)

	// when: We call `PreviousPage` with a given URL
	collection := createCollection("", ts.URL+"/measurement/measurements?source=1111111&pageSize=5&currentPage=1")
	_, error := api.PreviousPage(collection)

	// then: an error occurred
	if error == nil {
		t.Errorf("PreviousPage() should return error. Nil was given.")
	}
}

func createCollection(next string, prev string) *MeasurementCollection {
	return &MeasurementCollection{
		Next:       next,
		Self:       "https://t0818.cumulocity.com/measurement/measurements?source=1111111&pageSize=5&currentPage=2",
		Prev:       prev,
		Measurements:     []Measurement{},
		Statistics: nil,
	}
}
