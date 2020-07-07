package alarm

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"time"
)

func buildAlarmApi(url string) AlarmApi {
	httpClient := http.DefaultClient
	client := generic.Client{
		HTTPClient: httpClient,
		BaseURL:    url,
		Username:   "foo",
		Password:   "bar",
	}
	return NewAlarmApi(&client)
}

func buildHttpServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

var dateFrom, _ = time.Parse(time.RFC3339, "2020-06-29T10:11:12.000Z")
var dateTo, _ = time.Parse(time.RFC3339, "2020-06-30T13:14:15.000Z")

var deviceId = "1111111"
var alarmId = "2222222"
var alarm = `{
            "id": "2222222",
            "self": "https://t0815.cumulocity.com/alarm/alarms/2222222",
            "creationTime": "2020-06-30T08:32:04.413Z",
            "type": "test-gomulocity-Alarm",
            "time": "2020-06-30T08:32:04.261Z",
            "text": "Test creation of an alarm",
            "source": {
                "id": "1111111",
                "name": "testGomulocityDevice"
            },
            "status": "ACTIVE",
            "severity": "MINOR",
            "count": 1,
            "firstOccurrenceTime": "2020-06-30T08:32:04Z"
        }`

var alarmCollectionTemplate = `{
    "next": "https://t0818.cumulocity.com/alarm/alarms?source=1111111&pageSize=5&currentPage=2",
    "self": "https://t0815.cumulocity.com/alarm/alarms?source=1111111&pageSize=5&currentPage=1",
    "alarms": [%s],
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
    }
}`
