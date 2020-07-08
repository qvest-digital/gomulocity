package alarm

import (
	"testing"
)

func TestAlarmApi_Get_ExistingId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, alarm)
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	alarm, err := api.Get(alarmId)

	if err != nil {
		t.Fatalf("Get() got an unexpected error: %s", err.Error())
	}

	if alarm == nil {
		t.Fatalf("Get() returns an unexpected nil alarm.")
	}

	if alarm.Id != alarmId {
		t.Errorf("Get() alarm id = %v, want %v", alarm.Id, alarmId)
	}
}

func TestAlarmApi_Get_NotExistingId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(404, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	alarm, err := api.Get(alarmId)

	if err != nil {
		t.Fatalf("Get() got an unexpected error: %s", err.Error())
		return
	}

	if alarm != nil {
		t.Fatalf("Get() got an unexpected alarm. Should be nil.")
		return
	}
}

func TestAlarmApi_Get_MalformedJson(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, "{ foo ...")
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	_, err := api.Get(alarmId)

	if err == nil {
		t.Error("Get() Error expected but nil was given")
		return
	}
}
