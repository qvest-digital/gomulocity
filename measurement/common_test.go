package measurement

import (
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

func buildMeasurementApi(url string) MeasurementApi {
	httpClient := http.DefaultClient
	client := generic.Client{
		HTTPClient: httpClient,
		BaseURL:    url,
		Username:   "foo",
		Password:   "bar",
	}
	return NewMeasurementApi(&client)
}

func buildHttpServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func createMeasurementHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		var measurement NewMeasurement
		_ = generic.ObjectFromJson(body, &measurement)
		createMeasurementCapture = &measurement
		requestCapture = r
		bodyCapture = &body

		w.WriteHeader(status)
		responseCapture, _ := generic.JsonFromObject(responseMeasurement)
		_, _ = w.Write(responseCapture)
	}))
}

var createMeasurementCapture *NewMeasurement

var requestCapture *http.Request
var bodyCapture *[]byte
var measurementTime, _ = time.Parse(time.RFC3339, "2020-06-26T10:43:25.130Z")

var dateFrom, _ = time.Parse(time.RFC3339, "2020-06-29T10:11:12.000Z")
var dateTo, _ = time.Parse(time.RFC3339, "2020-06-30T13:14:15.000Z")

var deviceId = "1111111"
var measurementId = "2222222"

var responseMeasurement = &Measurement{
	Id:              "1337",
	MeasurementType: "TestMeasurement",
	Time:            &measurementTime,
	Source:          Source{Id: "4711"},
	Self:            "https://t0815.cumulocity.com/measurement/measurements/1337",
	Metrics: map[string]interface{}{
		"AirPressure": ValueFragment{Value: 1011.2, Unit: "hPa"},
		"Humidity":    ValueFragment{Value: 51, Unit: "%RH"},
		"Temperature": ValueFragment{Value: 23.45, Unit: "C"},
	},
}

var measurement = `{
	"id": "2222222",
	"self": "https://t0815.cumulocity.com/measurement/measurements/2222222",
	"type": "test-gomulocity-Measurement",
	"time": "2020-06-30T08:32:04.261Z",
	"source": {
		"id": "1111111"
	},
	"AirPressure":{"value":1011.2,"unit":"hPa"},
	"Humidity":{"value":51,"unit":"%RH"},
	"Temperature":{"value":23.45,"unit":"C"}
}`

var measurementCollectionTemplate = `{
    "next": "https://t0818.cumulocity.com/measurement/measurements?source=1111111&pageSize=5&currentPage=2",
    "self": "https://t0815.cumulocity.com/measurement/measurements?source=1111111&pageSize=5&currentPage=1",
    "measurements": [%s],
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
    }
}`
