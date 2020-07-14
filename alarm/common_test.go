package alarm

import (
	"encoding/json"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

func createAlarmHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		var alarm NewAlarm
		_ = generic.ObjectFromJson(body, &alarm)
		createAlarmCapture = &alarm
		bodyCapture = &body
		urlCapture = r.URL.Path
		requestCapture = r

		w.WriteHeader(status)
		responseCapture, _ := json.Marshal(responseAlarm)
		_, _ = w.Write(responseCapture)
	}))
}

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

var requestCapture *http.Request
var bodyCapture *[]byte
var urlCapture string
var createAlarmCapture *NewAlarm

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
