package measurement

import (
	"strings"
	"testing"
)

func TestMeasurementApi_Get_ExistingId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, measurement)
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)

	measurement, err := api.Get(measurementId)

	if err != nil {
		t.Fatalf("Get() got an unexpected error: %s", err.Error())
	}

	if measurement == nil {
		t.Fatalf("Get() returns an unexpected nil measurement.")
	}

	if measurement.Id != measurementId {
		t.Errorf("Get() measurement id = %v, want %v", measurement.Id, measurementId)
	}

	custom1, ok1 := measurement.Metrics["Custom1"].(string)
	custom2, ok2 := measurement.Metrics["Custom2"].(interface{})

	if !(ok1 && custom1 == "Hello world") {
		t.Errorf("GetForDevice() custom1 = %v, want %v", custom1, "Hello world")
	}
	if !(ok2 && custom2.(float64) == 1234) {
		t.Errorf("GetForDevice() custom2 = %v, want %v", custom2, 1234)
	}
}

func TestMeasurementApi_Get_NotExistingId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(404, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)

	measurement, err := api.Get(measurementId)

	if err != nil {
		t.Fatalf("Get() got an unexpected error: %s", err.Error())
		return
	}

	if measurement != nil {
		t.Fatalf("Get() got an unexpected measurement. Should be nil.")
		return
	}
}

func TestMeasurementApi_Get_Measurement_BadRequest(t *testing.T) {
	badRequest := `{
				"error": "bad request",
				"message": "Invalid request parameter!",
				"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
			}`
	// given: A test server
	ts := buildHttpServer(400, badRequest)
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)
	_, err := api.Get(measurementId)

	if err == nil {
		t.Errorf("Get() expected error on 400 - bad request")
		return
	}

	if !strings.Contains(err.ErrorType, "400") {
		t.Errorf("Get() expected error on 400 - bad request. Got: %s", err.ErrorType)
		return
	}
}

func TestMeasurementApi_Get_MalformedJson(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, "{ foo ...")
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)

	_, err := api.Get(measurementId)

	if err == nil {
		t.Error("Get() Error expected but nil was given")
		return
	}

	if !strings.Contains(err.Message, "Error while parsing") {
		t.Errorf("Get() expected parsing error. Got: %s", err.Message)
		return
	}
}

func TestMeasurementApi_Get_MalformedErrorJson(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(400, "{ foo ...")
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)

	_, err := api.Get(measurementId)

	if err == nil {
		t.Error("Get() Error expected but nil was given")
		return
	}

	if !strings.Contains(err.Message, "Error while parsing") {
		t.Errorf("Get() expected parsing error. Got: %s", err.Message)
		return
	}
}
