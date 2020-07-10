package measurement

import (
	"github.com/tarent/gomulocity/generic"
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

var requestCapture *http.Request
var measurementTime, _ = time.Parse(time.RFC3339, "2020-06-26T10:43:25.130Z")

var dateFrom, _ = time.Parse(time.RFC3339, "2020-06-29T10:11:12.000Z")
var dateTo, _ = time.Parse(time.RFC3339, "2020-06-30T13:14:15.000Z")

var deviceId = "1111111"
var measurementId = "2222222"
var measurement = `{
            "id": "2222222",
            "self": "https://t0815.cumulocity.com/measurement/measurements/2222222",
            "type": "test-gomulocity-Measurement",
            "time": "2020-06-30T08:32:04.261Z",
            "source": {
                "id": "1111111"
            }
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
