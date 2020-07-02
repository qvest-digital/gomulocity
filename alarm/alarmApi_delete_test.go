package alarm

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


var alarmFilter = AlarmFilter{
	Status:            []Status{ACTIVE,ACKNOWLEDGED},
	SourceId:          "123",
	WithSourceAssets:  true,
	WithSourceDevices: true,
	Resolved:          "false",
	Severity:          MAJOR,
	DateFrom:          &dateFrom,
	DateTo:            &dateTo,
	Type:              "testAlarm",
}

var expectedUrlParameters = "alarm/alarms?dateFrom=2020-06-29T10%3A11%3A12Z&dateTo=2020-06-30T13%3A14%3A15Z&resolved=false&severity=MAJOR&source=123&status=ACTIVE%2CACKNOWLEDGED&type=testAlarm&withSourceAssets=true&withSourceDevices=true"

func TestAlarmApi_Delete_Alarm_Success(t *testing.T) {
	var capturedUrl string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusNoContent)
	}))

	// given: A test server
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	err := api.Delete(&alarmFilter)

	if err != nil {
		t.Fatalf("DeleteAlarm() got an unexpected error: %s", err.Error())
	}

	if strings.Contains(capturedUrl, expectedUrlParameters) == false {
		t.Errorf("Delete() The target URL does not contains the alarmFilter: url: %s - expected %s", capturedUrl, expectedUrlParameters)
	}
}

func TestAlarmApi_Delete_Alarm_NotFound(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusNotFound, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	err := api.Delete(&alarmFilter)

	if err == nil {
		t.Errorf("Delete() expected error on 404 - not found")
		return
	}
}
