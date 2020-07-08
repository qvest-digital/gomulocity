package measurement

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


func TestMeasurementApi_Delete_Success(t *testing.T) {
	var expectedUrl = fmt.Sprintf("measurement/measurements/%s", measurementId)
	var capturedUrl string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusNoContent)
	}))

	// given: A test server
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)

	err := api.Delete(measurementId)

	if err != nil {
		t.Fatalf("DeleteMeasurement() got an unexpected error: %s", err.Error())
	}

	if !strings.HasSuffix(capturedUrl, expectedUrl) {
		t.Errorf("Delete() The target URL does not contains the measurementFilter: url: %s - expected %s", capturedUrl, expectedUrl)
	}
}

func TestMeasurementApi_Delete_NotFound(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusNotFound, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildMeasurementApi(ts.URL)

	err := api.Delete(measurementId)

	if err == nil {
		t.Errorf("Delete() expected error on 404 - not found")
		return
	}

	if !strings.Contains(err.ErrorType, "404") {
		t.Errorf("Delete() expected error on 404 - not found. Got: %s", err.ErrorType)
		return
	}
}

func TestMeasurementApi_Delete_WithoutId(t *testing.T) {
	// given: the api as system under test
	api := buildMeasurementApi("")

	err := api.Delete("")

	if err == nil {
		t.Errorf("Delete() expected error.")
		return
	}

	if !strings.Contains(err.Message, "without an id") {
		t.Errorf("Delete() expected error with appropriate message. Got: %s", err.Message)
		return
	}
}
