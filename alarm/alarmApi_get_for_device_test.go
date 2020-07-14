package alarm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAlarmApi_GetForDevice_ExistingId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, fmt.Sprintf(alarmCollectionTemplate, alarm))
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	collection, err := api.GetForDevice(deviceId, 5)

	if err != nil {
		t.Fatalf("GetForDevice() got an unexpected error: %s", err.Error())
	}

	if collection == nil {
		t.Fatalf("GetForDevice() got no explict error but the collection was nil.")
	}

	if len(collection.Alarms) != 1 {
		t.Fatalf("GetForDevice() alarms count = %v, want %v", len(collection.Alarms), 1)
	}

	alarm := collection.Alarms[0]
	if alarm.Id != alarmId {
		t.Errorf("GetForDevice() alarm id = %v, want %v", alarm.Id, alarmId)
	}
}

func TestEvents_GetForDevice_CustomElements(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, fmt.Sprintf(alarmCollectionTemplate, alarm))
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	collection, err := api.GetForDevice(deviceId, 5)

	if err != nil {
		t.Fatalf("GetForDevice() got an unexpected error: %s", err.Error())
	}

	if collection == nil {
		t.Fatalf("GetForDevice() got no explict error but the collection was nil.")
	}

	if len(collection.Alarms) != 1 {
		t.Fatalf("GetForDevice() alarms count = %v, want %v", len(collection.Alarms), 1)
	}

	alarm := collection.Alarms[0]
	if len(alarm.AdditionalFields) != 2 {
		t.Fatalf("GetForDevice() AdditionalFields length = %d, want %d", len(alarm.AdditionalFields), 2)
	}

	custom1, ok1 := alarm.AdditionalFields["custom1"].(string)
	custom2, ok2 := alarm.AdditionalFields["custom2"].([]interface{})

	if !(ok1 && custom1 == "Hello") {
		t.Errorf("GetForDevice() custom1 = %v, want %v", custom1, "Hello")
	}
	if !(ok2 && custom2[0] == "Foo" && custom2[1] == "Bar") {
		t.Errorf("GetForDevice() custom2 = [%v, %v], want [%v, %v]", custom2[0], custom2[1], "Foo", "Bar")
	}
}

func TestAlarmApi_GetForDevice_HandlesPageSize(t *testing.T) {
	tests := []struct {
		name        string
		pageSize    int
		errExpected bool
	}{
		{"Negative", -1, true},
		{"Zero", 0, true},
		{"Min", 1, false},
		{"Max", 2000, false},
		{"too large", 2001, true},
		{"in range", 10, false},
	}

	// given: A test server
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(alarmCollectionTemplate, alarm)))
	}))
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := api.GetForDevice(deviceId, tt.pageSize)

			if tt.errExpected {
				if err == nil {
					t.Error("GetForDevice() error expected but was nil")
				}
			}

			if !tt.errExpected {
				contains := strings.Contains(capturedUrl, fmt.Sprintf("pageSize=%d", tt.pageSize))

				if !contains {
					t.Errorf("GetForDevice() expected pageSize '%d' in url. '%s' given", tt.pageSize, capturedUrl)
				}
			}
		})
	}
}

func TestAlarmApi_GetForDevice_NotExistingId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, fmt.Sprintf(alarmCollectionTemplate, ""))
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	collection, err := api.GetForDevice(deviceId, 5)

	if err != nil {
		t.Fatalf("GetForDevice() got an unexpected error: %s", err.Error())
		return
	}

	if collection == nil {
		t.Fatalf("GetForDevice() got no explict error but the collection was nil.")
		return
	}

	if len(collection.Alarms) != 0 {
		t.Fatalf("GetForDevice() alarms count = %v, want %v", len(collection.Alarms), 0)
	}
}

func TestAlarmApi_GetForDevice_MalformedResponse(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, "{ foo ...")
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	_, err := api.GetForDevice(deviceId, 5)

	if err == nil {
		t.Errorf("GetForDevice() Expected error - non given")
		return
	}
}
