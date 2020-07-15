package alarm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func createAlarmHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		var alarm NewAlarm
		_ = json.Unmarshal(body, &alarm)
		createAlarmCapture = &alarm
		requestCapture = r

		w.WriteHeader(status)
		responseCapture, _ := json.Marshal(responseAlarm)
		_, _ = w.Write(responseCapture)
	}))
}

var requestCapture *http.Request
var createAlarmCapture *NewAlarm

var alarmTime, _ = time.Parse(time.RFC3339, "2020-06-26T10:43:25.130Z")
var responseAlarm = &Alarm{
	Id:                  "1337",
	Type:                "TestAlarm",
	Time:                &alarmTime,
	CreationTime:        &alarmTime,
	Text:                "This is my test alarm",
	Source:              Source{Id: "4711"},
	Self:                "https://t0815.cumulocity.com/alarm/alarms/1337",
	Status:              ACTIVE,
	Severity:            MAJOR,
	Count:               1,
	FirstOccurrenceTime: &alarmTime,
}

// given: A new alarm
var newAlarm = &NewAlarm{
	Type:   "TestAlarm",
	Time:   time.Time{},
	Text:   "This is my test alarm",
	Source: Source{Id: "4711"},
	Severity: MAJOR,
	Status: ACTIVE,
}

func TestAlarmApi_Create_Alarm_Success_SendsData(t *testing.T) {
	// given: A test server
	ts := createAlarmHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	_, err := api.Create(newAlarm)

	if err != nil {
		t.Fatalf("CreateAlarm() got an unexpected error: %s", err.Error())
	}

	if createAlarmCapture == nil {
		t.Fatalf("CreateAlarm() Captured alarm is nil.")
	}

	if !reflect.DeepEqual(newAlarm, createAlarmCapture) {
		t.Errorf("CreateAlarm() alarm = %v, want %v", newAlarm, createAlarmCapture)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.alarm+json"
	if header != want {
		t.Errorf("CreateAlarm() accept header = %v, want %v", header, want)
	}

	header = requestCapture.Header.Get("Content-Type")
	want = "application/vnd.com.nsn.cumulocity.alarm+json"
	if header != want {
		t.Errorf("CreateAlarm() Content-Type header = %v, want %v", header, want)
	}
}

func TestAlarmApi_Create_Alarm_Success_ReceivesData(t *testing.T) {
	// given: A test server
	ts := createAlarmHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	alarm, err := api.Create(newAlarm)

	if err != nil {
		t.Fatalf("CreateAlarm() got an unexpected error: %s", err.Error())
	}

	if !reflect.DeepEqual(alarm, responseAlarm) {
		t.Errorf("CreateAlarm() alarm = %v, want %v", alarm, responseAlarm)
	}
}

func TestAlarmApi_Create_Alarm_BadRequest(t *testing.T) {
	// given: A test server
	ts := createAlarmHttpServer(http.StatusBadRequest)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	_, err := api.Create(newAlarm)

	if err == nil {
		t.Errorf("CreateAlarm() expected error on 400 - bad request")
		return
	}

	if !strings.Contains(err.ErrorType, "400") {
		t.Errorf("CreateAlarm() expected error on 400 - bad request. Got: %s", err.ErrorType)
		return
	}
}
