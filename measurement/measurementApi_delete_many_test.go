package measurement

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


var measurementFilter = MeasurementQuery{
	DateFrom:            &dateFrom,
	DateTo:              &dateTo,
	Type:                "testMeasurement",
	ValueFragmentType:   "Temperature",
	ValueFragmentSeries: "Energy",
	SourceId:            "123",
}


func TestMeasurementApi_DeleteMany_Success(t *testing.T) {
	var expectedUrlParameter = "measurement/measurements?dateFrom=2020-06-29T10%3A11%3A12Z&dateTo=2020-06-30T13%3A14%3A15Z&source=123&type=testMeasurement&valueFragmentSeries=Energy&valueFragmentType=Temperature"
	var capturedUrl string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusNoContent)
	}))

	// given: A test server
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)

	err := api.DeleteMany(&measurementFilter)

	if err != nil {
		t.Fatalf("DeleteMeasurement() got an unexpected error: %s", err.Error())
	}

	if !strings.HasSuffix(capturedUrl, expectedUrlParameter) {
		t.Errorf("DeleteMany() The target URL does not contains the measurementFilter: url: %s - expected %s", capturedUrl, expectedUrlParameter)
	}
}

func TestMeasurementApi_DeleteMany_BadRequest(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusBadRequest, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)

	err := api.DeleteMany(&measurementFilter)

	if err == nil {
		t.Errorf("DeleteMany() expected error on 400 - bad request")
		return
	}

	if !strings.Contains(err.ErrorType, "400") {
		t.Errorf("DeleteMany() expected error on 400 - bad request. Got: %s", err.ErrorType)
		return
	}
}

func TestMeasurementApi_Delete_Many_WithoutFilter(t *testing.T) {
	// given: the api as system under test
	api := buildMeasurementApi("")

	err := api.DeleteMany(&MeasurementQuery{})

	if err == nil {
		t.Errorf("DeleteMany() expected error.")
		return
	}

	if !strings.Contains(err.Message, "No filter set") {
		t.Errorf("DeleteMany() expected error with appropriate message. Got: %s", err.Message)
		return
	}
}
