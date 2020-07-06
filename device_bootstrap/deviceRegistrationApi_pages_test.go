package device_bootstrap

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var deviceId = "4711"
var deviceRegistration = `{
		"id": "4711", 
		"status": "PENDING_ACCEPTANCE", 
		"self": "https://myFancyCloudInstance.com/devicecontrol/newDeviceRequests/4711",
		"owner": "me@company.com",
		"customProperties": {},
		"creationTime": "2020-07-03T10:16:35.870+02:00",
		"tenantId": "myCloud" 
}`

var deviceRegistrationCollectionTemplate = `{
	"next": "https://t0818.cumulocity.com/devicecontrol/newDeviceRequests?pageSize=5&currentPage=2",
	"self": "https://t0815.cumulocity.com/devicecontrol/newDeviceRequests?pageSize=5&currentPage=1", 
	"newDeviceRequests":[%s], 
	"statistics": {
		"pageSize":5, 
		"currentPage":1
	}
}`

func TestDeviceRegistrationApi_NextPage_Success(t *testing.T) {
	// given: An Http server with a next collection with one deviceRegistration.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(deviceRegistrationCollectionTemplate, deviceRegistration)))
	}))
	defer ts.Close()

	// and: the system under test
	api := buildDeviceRegistrationApi(ts)

	// when: We create an existing collection and call `NextPage`
	expectedUrl := ts.URL + "/devicecontrol/newDeviceRequests?pageSize=5&currentPage=3"
	collection := createCollection(expectedUrl, "")
	nextCollection, _ := api.NextPage(collection)

	// then: We got the next collection with one deviceRegistration.
	if capturedUrl != expectedUrl {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, expectedUrl)
	}

	if nextCollection == nil {
		t.Fatalf("NextPage() nextCollection is nil")
	}

	if len(nextCollection.DeviceRegistrations) != 1 {
		t.Fatalf("NextPage() captured URL = %v, expected %v", capturedUrl, expectedUrl)
	}

	deviceRegistration := nextCollection.DeviceRegistrations[0]
	if deviceRegistration.Id != deviceId {
		t.Errorf("NextPage() next deviceRegistration id = %v, expected %v", deviceRegistration.Id, deviceId)
	}
}

func TestDeviceRegistrationApi_NextPage_NotAvailable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(deviceRegistrationCollectionTemplate, "")))
	}))
	defer ts.Close()

	// given: The system under test
	api := buildDeviceRegistrationApi(ts)

	// when: We call `NextPage` with no URLs
	collection := createCollection("", "")
	nextCollection, _ := api.NextPage(collection)

	// then: No `nextCollection` is available.
	if nextCollection != nil {
		t.Errorf("NextPage() should return nil. Was: %v", nextCollection)
	}
}

func TestDeviceRegistrationApi_NextPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(deviceRegistrationCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildDeviceRegistrationApi(ts)

	// when: We call `NextPage` with a given URL
	collection := createCollection(ts.URL+"/devicecontrol/newDeviceRequests?pageSize=5&currentPage=3", "")
	nextCollection, _ := api.NextPage(collection)

	// then: `nextCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("NextPage() should return an empty collection on empty collection response.")
	}
}

func TestDeviceRegistrationApi_NextPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildDeviceRegistrationApi(ts)

	// when: We call `NextPage` with a given URL
	collection := createCollection(ts.URL+"/devicecontrol/newDeviceRequests?pageSize=5&currentPage=3", "")
	_, err := api.NextPage(collection)

	// then: an error occurred
	if err == nil {
		t.Errorf("NextPage() should return error. Nil was given.")
	}
}

func TestDeviceRegistrationApi_PreviousPage_Success(t *testing.T) {
	// given: A Http server with a previous collection with one deviceRegistration.
	var capturedUrl string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = "http://" + r.Host + r.URL.String()
		_, _ = w.Write([]byte(fmt.Sprintf(deviceRegistrationCollectionTemplate, deviceRegistration)))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildDeviceRegistrationApi(ts)

	// when: We create an existing collection and call `PreviousPage`
	expectedUrl := ts.URL + "/devicecontrol/newDeviceRequests?pageSize=5&currentPage=1"
	collection := createCollection("", expectedUrl)
	nextCollection, _ := api.PreviousPage(collection)

	// then: We got the previous collection with one deviceRegistration.
	if capturedUrl != expectedUrl {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, expectedUrl)
	}

	if nextCollection == nil {
		t.Fatalf("PreviousPage() nextCollection is nil")
	}

	if len(nextCollection.DeviceRegistrations) != 1 {
		t.Fatalf("PreviousPage() captured URL = %v, expected %v", capturedUrl, expectedUrl)
	}

	deviceRegistration := nextCollection.DeviceRegistrations[0]
	if deviceRegistration.Id != deviceId {
		t.Errorf("PreviousPage() next deviceRegistration id = %v, expected %v", deviceRegistration.Id, deviceId)
	}
}

func TestDeviceRegistrationApi_PreviousPage_NotAvailable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(deviceRegistrationCollectionTemplate, "")))
	}))
	defer ts.Close()

	// given: The system under test
	api := buildDeviceRegistrationApi(ts)

	// when: We call `PreviousPage` with no URLs
	collection := createCollection("", "")
	nextCollection, _ := api.PreviousPage(collection)

	// then: No `previousCollection` is available.
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return nil. Was: %v", nextCollection)
	}
}

func TestDeviceRegistrationApi_PreviousPage_Empty(t *testing.T) {
	// given: A Http server with a next, but empty collection.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf(deviceRegistrationCollectionTemplate, "")))
	}))
	defer ts.Close()

	// and: The system under test
	api := buildDeviceRegistrationApi(ts)

	// when: We call `PreviousPage` with a given URL
	collection := createCollection(ts.URL+"/devicecontrol/newDeviceRequests?pageSize=5&currentPage=3", "")
	nextCollection, _ := api.NextPage(collection)

	// then: `previousCollection` ist `nil`
	if nextCollection != nil {
		t.Errorf("PreviousPage() should return an empty collection on empty collection response.")
	}
}

func TestDeviceRegistrationApi_PreviousPage_Error(t *testing.T) {
	// given: A Http server and an internal server error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	// and: The system under test
	api := buildDeviceRegistrationApi(ts)

	// when: We call `PreviousPage` with a given URL
	collection := createCollection("", ts.URL+"/devicecontrol/newDeviceRequests?pageSize=5&currentPage=1")
	_, error := api.PreviousPage(collection)

	// then: an error occurred
	if error == nil {
		t.Errorf("PreviousPage() should return error. Nil was given.")
	}
}

func createCollection(next string, prev string) *DeviceRegistrationCollection {
	return &DeviceRegistrationCollection{
		Next:                next,
		Self:                "https://t0818.cumulocity.com/devicecontrol/newDeviceRequests?pageSize=5&currentPage=2",
		Prev:                prev,
		DeviceRegistrations: []DeviceRegistration{},
		Statistics:          nil,
	}
}
