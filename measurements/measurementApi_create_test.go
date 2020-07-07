package measurements

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func createMeasurementHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		var measurement Measurement
		_ = json.Unmarshal(body, &measurement)
		createMeasurementCapture = &measurement
		requestCapture = r

		w.WriteHeader(status)
		responseCapture, _ := json.Marshal(responseMeasurement)
		_, _ = w.Write(responseCapture)
	}))
}

var requestCapture *http.Request
var createMeasurementCapture *Measurement

var measurementTime, _ = time.Parse(time.RFC3339, "2020-06-26T10:43:25.130Z")
var responseMeasurement = &Measurement{
	Id:                  "1337",
	MeasurementType:                "TestMeasurement",
	Time:                &measurementTime,
	Source:              Source{Id: "4711"},
	Self:                "https://t0815.cumulocity.com/measurement/measurements/1337",
}

// given: A new measurement
var newMeasurement = &Measurement{
	MeasurementType:   "TestMeasurement",
	Time:   &time.Time{},
	Source: Source{Id: "4711"},
}

func TestMeasurementApi_Create_Measurement_Success_SendsData(t *testing.T) {
	// given: A test server
	ts := createMeasurementHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)
	_, err := api.Create(newMeasurement)

	if err != nil {
		t.Fatalf("CreateMeasurement() got an unexpected error: %s", err.Error())
	}

	if createMeasurementCapture == nil {
		t.Fatalf("CreateMeasurement() Captured measurement is nil.")
	}

	if !reflect.DeepEqual(newMeasurement, createMeasurementCapture) {
		t.Errorf("CreateMeasurement() measurement = %v, want %v", newMeasurement, createMeasurementCapture)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.measurement+json;charset=UTF-8;ver=0.9"
	if header != want {
		t.Errorf("CreateMeasurement() accept header = %v, want %v", header, want)
	}

	header = requestCapture.Header.Get("Content-Type")
	want = "application/vnd.com.nsn.cumulocity.measurement+json;charset=UTF-8;ver=0.9"
	if header != want {
		t.Errorf("CreateMeasurement() Content-Type header = %v, want %v", header, want)
	}
}

func TestMeasurementApi_Create_Measurement_Success_ReceivesData(t *testing.T) {
	// given: A test server
	ts := createMeasurementHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)
	measurement, err := api.Create(newMeasurement)

	if err != nil {
		t.Fatalf("CreateMeasurement() got an unexpected error: %s", err.Error())
	}

	if !reflect.DeepEqual(measurement, responseMeasurement) {
		t.Errorf("CreateMeasurement() measurement = %v, want %v", measurement, responseMeasurement)
	}
}

func TestMeasurementApi_Create_Measurement_BadRequest(t *testing.T) {
	// given: A test server
	ts := createMeasurementHttpServer(400)
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)
	_, err := api.Create(newMeasurement)

	if err == nil {
		t.Errorf("CreateMeasurement() expected error on 400 - bad request")
		return
	}
}
