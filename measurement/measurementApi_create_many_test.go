package measurement

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

var responseMeasurementCollection = &MeasurementCollection{
	Measurements: []Measurement{
		{
			Id:              "1337",
			MeasurementType: "TestMeasurement1",
			Time:            &measurementTime,
			Source:          Source{Id: "4711"},
			Self:            "https://t0815.cumulocity.com/measurement/measurements/1337",
			Metrics: map[string]interface{}{
				"AirPressure": ValueFragment{Value: 1011.2, Unit: "hPa"},
				"Humidity":    ValueFragment{Value: 51, Unit: "%RH"},
				"Temperature": ValueFragment{Value: 23.45, Unit: "C"},
			},
		},
		{
			Id:              "1007",
			MeasurementType: "TestMeasurement2",
			Time:            &measurementTime,
			Source:          Source{Id: "4711"},
			Self:            "https://t0815.cumulocity.com/measurement/measurements/1007",
			Metrics: map[string]interface{}{
				"AirPressure": ValueFragment{Value: 1011.2, Unit: "hPa"},
				"Humidity":    ValueFragment{Value: 51, Unit: "%RH"},
				"Temperature": ValueFragment{Value: 23.45, Unit: "C"},
			},
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
			Metrics: map[string]interface{}{
				"AirPressure": ValueFragment{Value: 1011.2, Unit: "hPa"},
				"Humidity":    ValueFragment{Value: 51, Unit: "%RH"},
				"Temperature": ValueFragment{Value: 23.45, Unit: "C"},
			},
		},
		{
			MeasurementType: "TestMeasurement2",
			Time:            &measurementTime,
			Source:          Source{Id: "4711"},
			Metrics: map[string]interface{}{
				"AirPressure": ValueFragment{Value: 1011.2, Unit: "hPa"},
				"Humidity":    ValueFragment{Value: 51, Unit: "%RH"},
				"Temperature": ValueFragment{Value: 23.45, Unit: "C"},
			},
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

	assertNewMeasurementCollection(newMeasurements, measurementCollection, t)

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

	assertMeasurementCollection(measurements, responseMeasurementCollection, t)
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

func TestMeasurementApi_CreateMany_Measurement_Flats_Metrics(t *testing.T) {
	// given: A test server
	ts := createManyMeasurementsHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)
	_, err := api.CreateMany(newMeasurements)

	if err != nil {
		t.Fatalf("CreateMany() got an unexpected error: %s", err.Error())
	}

	if bodyCapture == nil {
		t.Fatalf("CreateMeasurement() Captured body is nil.")
	}

	// and: The body is a json structure
	var bodyMap map[string]interface{}
	jErr := json.Unmarshal(*bodyCapture, &bodyMap)

	if jErr != nil {
		t.Fatalf("CreateMeasurement() request body can not be parsed %v", err)
	}

	m := bodyMap["measurements"]
	measurements := m.([]interface{})
	assertMetricsOfMeasurement(measurements[0].(map[string]interface{}), t)
	assertMetricsOfMeasurement(measurements[1].(map[string]interface{}), t)
}

func assertMeasurementCollection(given *MeasurementCollection, want *MeasurementCollection, t *testing.T) {
	if given.Next != want.Next ||
		given.Prev != want.Prev ||
		given.Self != want.Self ||
		!reflect.DeepEqual(given.Statistics, want.Statistics) {

		t.Errorf("CreateMany()\n measurement = %v\n want %v", given, want)
	}

	for i, g := range given.Measurements {
		assertCommonMeasurement(&g, &want.Measurements[i], t)
	}
}

func assertNewMeasurementCollection(given *NewMeasurements, want *NewMeasurements, t *testing.T) {
	for i, g := range given.Measurements {
		assertCommonNewMeasurement(&g, &want.Measurements[i], t)
	}
}
