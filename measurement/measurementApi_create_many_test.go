package measurement

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func createManyMeasurementsHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		var measurement NewMeasurements
		_ = json.Unmarshal(body, &measurement)
		measurementCollection = &measurement
		requestCapture = r

		w.WriteHeader(status)
		responseCapture, _ := json.Marshal(responseMeasurementCollection)
		_, _ = w.Write(responseCapture)
	}))
}

var measurementCollection *NewMeasurements

var responseMeasurementCollection = &MeasurementCollection{
	Measurements: []Measurement{
		{
			Id:              "1337",
			MeasurementType: "TestMeasurement1",
			Time:            &measurementTime,
			Source:          Source{Id: "4711"},
			Self:            "https://t0815.cumulocity.com/measurement/measurements/1337",
		},
		{
			Id:              "1007",
			MeasurementType: "TestMeasurement2",
			Time:            &measurementTime,
			Source:          Source{Id: "4711"},
			Self:            "https://t0815.cumulocity.com/measurement/measurements/1007",
		},
	},
}

// given: new measurements
var newMeasurements = &NewMeasurements{
	Measurements: []NewMeasurement{
		{
			MeasurementType: "TestMeasurement1",
			Time:            &measurementTime,
			Source:          Source{Id: "4711"},
		},
		{
			MeasurementType: "TestMeasurement2",
			Time:            &measurementTime,
			Source:          Source{Id: "4711"},
		},
	},
}

func TestMeasurementApi_CreateMany_Success_SendsData(t *testing.T) {
	// given: A test server
	ts := createManyMeasurementsHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)
	_, err := api.CreateMany(newMeasurements)

	if err != nil {
		t.Fatalf("CreateMany() got an unexpected error: %s", err.Error())
	}

	if measurementCollection == nil {
		t.Fatalf("CreateMany() Captured measurement is nil.")
	}

	if !reflect.DeepEqual(newMeasurements, measurementCollection) {
		t.Errorf("CreateMany() measurement = %v, want %v", newMeasurements, measurementCollection)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.measurementCollection+json;charset=UTF-8;ver=0.9"
	if header != want {
		t.Errorf("CreateMany() accept header = %v, want %v", header, want)
	}

	header = requestCapture.Header.Get("Content-Type")
	want = "application/vnd.com.nsn.cumulocity.measurementCollection+json;charset=UTF-8;ver=0.9"
	if header != want {
		t.Errorf("CreateMany() Content-Type header = %v, want %v", header, want)
	}
}

func TestMeasurementApi_CreateMany_Success_ReceivesData(t *testing.T) {
	// given: A test server
	ts := createManyMeasurementsHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)
	measurements, err := api.CreateMany(newMeasurements)

	if err != nil {
		t.Fatalf("CreateMany() got an unexpected error: %s", err.Error())
	}

	if !reflect.DeepEqual(measurements, responseMeasurementCollection) {
		t.Errorf("CreateMany() measurements = %v, want %v", measurements, responseMeasurementCollection)
	}
}

func TestMeasurementApi_CreateMany_BadRequest(t *testing.T) {
	// given: A test server
	ts := createManyMeasurementsHttpServer(400)
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)
	_, err := api.CreateMany(newMeasurements)

	if err == nil {
		t.Errorf("CreateMany() expected error on 400 - bad request")
		return
	}

	if !strings.Contains(err.ErrorType, "400") {
		t.Errorf("CreateMany() expected error on 400 - bad request. Got: %s", err.ErrorType)
		return
	}
}
