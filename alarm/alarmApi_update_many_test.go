package alarm

import (
	"strings"
	"testing"
)

// given: A new alarm status and an updateAlarmsFilter
var newAlarmStatus = ACKNOWLEDGED
var expectedStatusUpdate = &UpdateAlarm{
	Status: newAlarmStatus,
}

var updateAlarmsFilter = &UpdateAlarmsFilter{
	Resolved: "false",
	Status:   ACTIVE,
	SourceId: "123",
	Severity: MINOR,
	DateFrom: &dateFrom,
	DateTo:   &dateTo,
}

var expectedUpdateFilter = "/alarm/alarms?dateFrom=2020-06-29T10%3A11%3A12Z&dateTo=2020-06-30T13%3A14%3A15Z&resolved=false&severity=MINOR&source=123&status=ACTIVE"

func TestAlarmApi_BulkStatusUpdate_Alarm_Success_SendsData(t *testing.T) {
	// given: A test server
	ts := updateAlarmHttpServer(200)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	err := api.BulkStatusUpdate(updateAlarmsFilter, newAlarmStatus)

	if err != nil {
		t.Fatalf("BulkStatusUpdate() got an unexpected error: %s", err.Error())
	}

	if strings.Contains(urlCapture, expectedUpdateFilter) == false {
		t.Errorf("BulkStatusUpdate() The target URL does not contains the request parameters: url: [%s] - expected [%s]",
			urlCapture, expectedUpdateFilter)
	}

	if updateAlarmCapture == nil {
		t.Fatalf("BulkStatusUpdate() Captured alarm is nil.")
	}

	if expectedStatusUpdate.Status != updateAlarmCapture.Status {
		t.Errorf("BulkStatusUpdate()\n alarm = %v\n want %v", updateAlarmCapture, expectedStatusUpdate)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.alarm+json"
	if header != want {
		t.Errorf("BulkStatusUpdate() accept header = %v, want %v", header, want)
	}
}

func TestAlarmApi_BulkStatusUpdate_Alarm_Success_Background(t *testing.T) {
	// given: A test server
	ts := updateAlarmHttpServer(202)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	err := api.BulkStatusUpdate(updateAlarmsFilter, newAlarmStatus)

	if err != nil {
		t.Fatalf("BulkStatusUpdate() got an unexpected error: %s", err.Error())
	}
}

func TestAlarmApi_BulkStatusUpdate_Alarm_BadRequest(t *testing.T) {
	// given: A test server
	ts := updateAlarmHttpServer(400)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	err := api.BulkStatusUpdate(updateAlarmsFilter, newAlarmStatus)

	if err == nil {
		t.Errorf("BulkStatusUpdate() expected error on 400 - bad request")
		return
	}
}
