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

func bulkStatusUpdateAlarmsHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		var statusUpdate UpdateAlarm
		_ = json.Unmarshal(body, &statusUpdate)
		statusUpdateCapture = &statusUpdate
		updateUrlCapture = r.URL.String()
		requestCapture = r

		w.WriteHeader(status)
		response, _ := json.Marshal("")
		_, _ = w.Write(response)
	}))
}

var statusUpdateCapture *UpdateAlarm

// given: A new alarm status and an updateAlarmsFilter
var newAlarmStatus = ACKNOWLEDGED
var expectedStatusUpdate = &UpdateAlarm {
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


func TestAlarmApi_BulkStatusUpdate_Alarm_Success_SendsData(t *testing.T) {
	// given: A test server
	ts := bulkStatusUpdateAlarmsHttpServer(200)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	err := api.BulkStatusUpdate(updateAlarmsFilter, newAlarmStatus)

	if err != nil {
		t.Fatalf("BulkStatusUpdate() got an unexpected error: %s", err.Error())
	}

	if strings.Contains(updateUrlCapture, expectedUpdateFilter) == false {
		t.Errorf("BulkStatusUpdate() The target URL does not contains the request parameters: url: [%s] - expected [%s]", updateUrlCapture, expectedUpdateFilter)
	}

	if statusUpdateCapture == nil {
		t.Fatalf("BulkStatusUpdate() Captured alarm is nil.")
	}

	if !reflect.DeepEqual(expectedStatusUpdate, statusUpdateCapture) {
		t.Errorf("BulkStatusUpdate() alarm = %v, want %v", statusUpdateCapture, expectedStatusUpdate)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.alarm+json"
	if header != want {
		t.Errorf("BulkStatusUpdate() accept header = %v, want %v", header, want)
	}
}

func TestAlarmApi_BulkStatusUpdate_Alarm_Success_Background(t *testing.T) {
	// given: A test server
	ts := bulkStatusUpdateAlarmsHttpServer(202)
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
	ts := bulkStatusUpdateAlarmsHttpServer(400)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	err := api.BulkStatusUpdate(updateAlarmsFilter, newAlarmStatus)

	if err == nil {
		t.Errorf("BulkStatusUpdate() expected error on 400 - bad request")
		return
	}
}
