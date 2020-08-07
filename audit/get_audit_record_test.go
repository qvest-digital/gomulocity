package audit

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAuditApi_GetAuditRecord(t *testing.T) {
	expectedUrl := fmt.Sprintf("/audit/auditRecords/%v", auditID)
	var capturedUrl string
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testAuditRecordJSON))
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// when
	auditRecord, err := api.GetAuditRecord(auditID)

	// then
	if err != nil {
		t.Fatalf("received an unexpected error: %s", err)
	}

	if auditRecord != nil {
		if !reflect.DeepEqual(testAuditRecord, auditRecord) {
			t.Errorf("GetAuditRecord() auditRecord: %v\nwant: %v", auditRecord, testAuditRecord)
		}
	} else {
		t.Error("audit record must not be nil")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}
}

func TestAuditApi_GetAuditRecord_without_recordID(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// when
	_, err := api.GetAuditRecord("")

	// then
	expectedErr := generic.ClientError("Getting an audit record without recordID is not allowed", "GetAuditRecord")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}

func TestAuditApi_GetAuditRecord_error_marshalling(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<>"))
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// when
	_, err := api.GetAuditRecord(auditID)

	// then
	expectedErr := generic.ClientError("Error while unmarshalling response: invalid character '<' looking for beginning of value", "GetAuditRecord")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}

func TestAuditApi_GetAuditRecord_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// when
	_, err := api.GetAuditRecord(auditID)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError, "")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}
