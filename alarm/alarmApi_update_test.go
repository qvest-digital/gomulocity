package alarm

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// given: An update alarm
var alarmUpdate = &UpdateAlarm{
	Text:             "This is my updated test alarm",
	Status:           ACKNOWLEDGED,
	Severity:         WARNING,
	AdditionalFields: map[string]interface{}{},
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

	if strings.Contains(urlCapture, alarmId) == false {
		t.Errorf("UpdateAlarm() The target URL does not contains the alarm Id: url: %s - expected alarmId %s", urlCapture, alarmId)
	}

	if updateAlarmCapture == nil {
		t.Fatalf("UpdateAlarm() Captured alarm is nil.")
	}

	if !reflect.DeepEqual(alarmUpdate, updateAlarmCapture) {
		t.Errorf("UpdateAlarm() \n alarm = %v\n want %v", alarmUpdate, updateAlarmCapture)
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

func TestAlarmApi_Update_Alarm_CustomFields(t *testing.T) {
	// given: A test server
	ts := updateAlarmHttpServer(200)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	alarmUpdate = &UpdateAlarm{
		Text:     "This is my test alarm",
		Severity: MAJOR,
		Status:   ACTIVE,
		AdditionalFields: map[string]interface{}{
			"Custom1": 4711,
			"Custom2": "Hello World",
		},
	}

	_, err := api.Update(alarmId, alarmUpdate)

	if err != nil {
		t.Fatalf("UpdateAlarm() got an unexpected error: %s", err.Error())
	}

	// and: A body was captured
	if bodyCapture == nil {
		t.Fatalf("UpdateAlarm() Captured request is nil.")
	}

	// and: The body is a json structure
	var bodyMap map[string]interface{}
	jErr := json.Unmarshal(*bodyCapture, &bodyMap)

	if jErr != nil {
		t.Fatalf("UpdateAlarm() request body can not be parsed %v", err)
	}

	// and: The "Custom1" and "Custom2" field is flattened
	custom1, _ := bodyMap["Custom1"].(float64)
	custom2, _ := bodyMap["Custom2"].(string)
	if custom1 != 4711 || custom2 != "Hello World" {
		t.Errorf("UpdateAlarm() additional fields - \ncustom fields = [%.2f, %s] \nwant [4711, Hello World]", custom1, custom2)
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
