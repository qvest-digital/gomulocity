package measurement

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// given: A new measurement
var newMeasurement = &NewMeasurement{
	MeasurementType: "TestMeasurement",
	Time:            &measurementTime,
	Source:          Source{Id: "4711"},
	Metrics: map[string]interface{}{
		"AirPressure": ValueFragment{Value: 1011.2, Unit: "hPa"},
		"Humidity":    ValueFragment{Value: 51, Unit: "%RH"},
		"Temperature": ValueFragment{Value: 23.45, Unit: "C"},
	},
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

	if newMeasurement.MeasurementType != createMeasurementCapture.MeasurementType ||
		!newMeasurement.Time.Equal(*createMeasurementCapture.Time) ||
		!reflect.DeepEqual(newMeasurement.Source, createMeasurementCapture.Source) {

		t.Errorf("CreateMeasurement()\n measurement = %v\n want %v", measurement, responseMeasurement)
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

func TestMeasurementApi_Create_Measurement_Flats_Metrics(t *testing.T) {
	// given: A test server
	ts := createMeasurementHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)
	_, err := api.Create(newMeasurement)

	if err != nil {
		t.Fatalf("CreateMeasurement() got an unexpected error: %s", err.Error())
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

	// and: The "Custom1" and "Custom2" field is flattened
	temperature, _ := bodyMap["Temperature"].(map[string]interface{})
	if temperature["value"] != 23.45 || temperature["unit"] != "C" {
		t.Errorf("CreateEvent() metrics\n temperature = %.2f, %s \nwant {23.45, C}", temperature["value"], temperature["unit"])
	}

	humidity, _ := bodyMap["Humidity"].(map[string]interface{})
	if humidity["value"] != 51.00 || humidity["unit"] != "%RH" {
		t.Errorf("CreateEvent() metrics\n humidity = %.2f, %s \nwant {51.00, %%RH}", humidity["value"], humidity["unit"])
	}

	airPressure, _ := bodyMap["AirPressure"].(map[string]interface{})
	if airPressure["value"] != 1011.2 || airPressure["unit"] != "hPa" {
		t.Errorf("CreateEvent() metrics\n air pressure = %.2f, %s \nwant {1011.2, hPa}", airPressure["value"], airPressure["unit"])
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

	if measurement.Id != responseMeasurement.Id ||
		measurement.MeasurementType != responseMeasurement.MeasurementType ||
		!measurement.Time.Equal(*responseMeasurement.Time) ||
		!reflect.DeepEqual(measurement.Source, responseMeasurement.Source) {

		t.Errorf("CreateMeasurement()\n measurement = %v\n want %v", measurement, responseMeasurement)
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

	if !strings.Contains(err.ErrorType, "400") {
		t.Errorf("CreateMeasurement() expected error on 400 - bad request. Got: %s", err.ErrorType)
		return
	}
}
