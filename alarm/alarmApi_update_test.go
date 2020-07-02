package alarm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func updateAlarmHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		var alarm UpdateAlarm
		_ = json.Unmarshal(body, &alarm)
		alarmUpdateCapture = &alarm
		updateUrlCapture = r.URL.Path
		requestCapture = r

		w.WriteHeader(status)
		response, _ := json.Marshal(responseAlarm)
		_, _ = w.Write(response)
	}))
}

var alarmUpdateCapture *UpdateAlarm
var updateUrlCapture string

// given: An update alarm
var alarmUpdate = &UpdateAlarm{
	Text:     "This is my updated test alarm",
	Status:   ACKNOWLEDGED,
	Severity: WARNING,
}

func TestAlarmApi_Update_Alarm_Success_SendsData(t *testing.T) {
	// given: A test server
	ts := updateAlarmHttpServer(200)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	_, err := api.Update(alarmId, alarmUpdate)

	if err != nil {
		t.Fatalf("UpdateAlarm() got an unexpected error: %s", err.Error())
	}

	if strings.Contains(updateUrlCapture, alarmId) == false {
		t.Errorf("UpdateAlarm() The target URL does not contains the alarm Id: url: %s - expected alarmId %s", updateUrlCapture, alarmId)
	}

	if alarmUpdateCapture == nil {
		t.Fatalf("UpdateAlarm() Captured alarm is nil.")
	}

	if !reflect.DeepEqual(alarmUpdate, alarmUpdateCapture) {
		t.Errorf("UpdateAlarm() alarm = %v, want %v", alarmUpdate, alarmUpdateCapture)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.alarm+json"
	if header != want {
		t.Errorf("UpdateAlarm() accept header = %v, want %v", header, want)
	}

	header = requestCapture.Header.Get("Content-Type")
	want = "application/vnd.com.nsn.cumulocity.alarm+json"
	if header != want {
		t.Errorf("UpdateAlarm() content-type header = %v, want %v", header, want)
	}
}

func TestAlarmApi_Update_Alarm_Success_ReceivesData(t *testing.T) {
	// given: A test server
	ts := updateAlarmHttpServer(200)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	alarm, err := api.Update(alarmId, alarmUpdate)

	if err != nil {
		t.Fatalf("UpdateAlarm() got an unexpected error: %s", err.Error())
	}

	if !reflect.DeepEqual(alarm, responseAlarm) {
		t.Errorf("UpdateAlarm() alarm = %v, want %v", alarm, responseAlarm)
	}
}

func TestAlarmApi_Update_Alarm_BadRequest(t *testing.T) {
	// given: A test server
	ts := createAlarmHttpServer(400)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	_, err := api.Update(alarmId, alarmUpdate)

	if err == nil {
		t.Errorf("UpdateAlarm() expected error on 400 - bad request")
		return
	}
}
