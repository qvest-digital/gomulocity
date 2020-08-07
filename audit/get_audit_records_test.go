package audit

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAuditApi_GetAuditRecords(t *testing.T) {
	expectedUrl := "/audit/auditRecords?application=Omniscape&pageSize=2&revert=false&type=Alarm&user=gfa-agent"
	var capturedUrl string
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testAuditRecordsJSON))
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// and: query params
	query := &AuditQuery{
		Type:        "Alarm",
		User:        "gfa-agent",
		Application: "Omniscape",
	}

	// when
	auditRecords, err := api.GetAuditRecords(query, 2)

	// then
	if err != nil {
		t.Fatalf("received an unexpected error: %s", err)
	}

	if auditRecords != nil {
		if len(auditRecords.AuditRecords) == 0 {
			t.Error("audit record collection does not contain any audit records")
		}
		if !reflect.DeepEqual(testAuditRecords, auditRecords) {
			t.Errorf("GetAuditRecords() auditRecord: %v\nwant: %v", auditRecords, testAuditRecord)
		}
	} else {
		t.Error("audit record must not be nil")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}
}

func TestAuditApi_GetAuditRecords_error_marshalling(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<>"))
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// and: query params
	query := &AuditQuery{
		Type:        "Alarm",
		User:        "gfa-agent",
		Application: "Omniscape",
	}

	// when
	_, err := api.GetAuditRecords(query, 2)

	// then
	expectedErr := generic.ClientError("Error while parsing response JSON: Error while unmarshalling json: invalid character '<' looking for beginning of value", "CollectionResponseParser")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}

func TestAuditApi_GetAuditRecords_empty_response(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// and: query params
	query := &AuditQuery{
		Type:        "Alarm",
		User:        "gfa-agent",
		Application: "Omniscape",
	}

	// when
	_, err := api.GetAuditRecords(query, 2)

	// then
	expectedErr := generic.ClientError("Response body was empty", "CollectionResponseParser")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}

func TestAuditApi_GetAuditRecords_invalid_pageSize(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// and: query params
	query := &AuditQuery{
		Type:        "Alarm",
		User:        "gfa-agent",
		Application: "Omniscape",
	}

	// when
	_, err := api.GetAuditRecords(query, -1)

	// then
	expectedErr := generic.ClientError("Error while building pageSize parameter to fetch audit records: The page size must be between 1 and 2000. Was -1", "FindAuditRecords")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}
