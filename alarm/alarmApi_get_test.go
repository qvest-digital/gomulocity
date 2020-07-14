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

func TestEvents_Get_CustomElements(t *testing.T) {
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

	if len(alarm.AdditionalFields) != 2 {
		t.Fatalf("Get() AdditionalFields length = %d, want %d", len(alarm.AdditionalFields), 2)
	}

	custom1, ok1 := alarm.AdditionalFields["custom1"].(string)
	custom2, ok2 := alarm.AdditionalFields["custom2"].([]interface{})

	if !(ok1 && custom1 == "Hello") {
		t.Errorf("Get() custom1 = %v, want %v", custom1, "Hello")
	}
	if !(ok2 && custom2[0] == "Foo" && custom2[1] == "Bar") {
		t.Errorf("Get() custom2 = [%v, %v], want [%v, %v]", custom2[0], custom2[1], "Foo", "Bar")
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
