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

func updateManyAlarmsHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		var statusUpdate StatusUpdate
		_ = json.Unmarshal(body, &statusUpdate)
		statusUpdateCapture = &statusUpdate
		updateUrlCapture = r.URL.Path
		requestCapture = r

		w.WriteHeader(status)
		response, _ := json.Marshal("")
		_, _ = w.Write(response)
	}))
}

type StatusUpdate struct {
	Status	Status `json:"status"`
}

var statusUpdateCapture *StatusUpdate

// given: A new alarm status and an updateAlarmsFilter
var newAlarmStatus = ACKNOWLEDGED
var expectedStatusUpdate = StatusUpdate {
	Status:   newAlarmStatus,
}

var updateAlarmsFilter = &UpdateAlarmsFilter{
	Resolved: "false",
	Status:   ACTIVE,
	SourceId: "123",
	Severity: MINOR,
	DateFrom: &dateFrom,
	DateTo:   &dateTo,
}

var expectedUpdateFilter = "alarm/alarms?dateFrom=2020-06-29T10%3A11%3A12Z&dateTo=2020-06-30T13%3A14%3A15Z&resolved=false&severity=MINOR&source=123&status=ACTIVE"


func TestAlarmApi_UpdateMany_Alarm_Success_SendsData(t *testing.T) {
	// given: A test server
	ts := updateManyAlarmsHttpServer(200)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	err := api.UpdateMany(updateAlarmsFilter, newAlarmStatus)

	if err != nil {
		t.Fatalf("UpdateManyAlarms() got an unexpected error: %s", err.Error())
	}

	if strings.Contains(updateUrlCapture, expectedUpdateFilter) == false {
		t.Errorf("UpdateManyAlarms() The target URL does not contains the request parameters: url: [%s] - expected [%s]", updateUrlCapture, expectedUpdateFilter)
	}

	if alarmUpdateCapture == nil {
		t.Fatalf("UpdateManyAlarms() Captured alarm is nil.")
	}

	if !reflect.DeepEqual(expectedStatusUpdate, statusUpdateCapture) {
		t.Errorf("UpdateManyAlarms() alarm = %v, want %v", statusUpdateCapture, expectedStatusUpdate)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.alarm+json"
	if header != want {
		t.Errorf("UpdateManyAlarms() accept header = %v, want %v", header, want)
	}
}

func TestAlarmApi_UpdateMany_Alarm_Success_Background(t *testing.T) {
	// given: A test server
	ts := updateManyAlarmsHttpServer(202)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	err := api.UpdateMany(updateAlarmsFilter, newAlarmStatus)

	if err != nil {
		t.Fatalf("UpdateManyAlarms() got an unexpected error: %s", err.Error())
	}
}

func TestAlarmApi_UpdateMany_Alarm_BadRequest(t *testing.T) {
	// given: A test server
	ts := updateManyAlarmsHttpServer(400)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	err := api.UpdateMany(updateAlarmsFilter, newAlarmStatus)

	if err == nil {
		t.Errorf("UpdateManyAlarms() expected error on 400 - bad request")
		return
	}
}
