package audit

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuditApi_CreateAuditRecord(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(testAuditRecordJSON))
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// when
	auditRecord, err := api.CreateAuditRecord(testCreateAuditRecord)

	// then
	if err != nil {
		t.Fatalf("received an unexpected error: %s", err)
	}

	if auditRecord != nil {
		if auditRecord.Activity != testCreateAuditRecord.Activity ||
			auditRecord.Type != testCreateAuditRecord.Type ||
			auditRecord.Time != testCreateAuditRecord.Time ||
			auditRecord.Text != testCreateAuditRecord.Text {
			t.Errorf("respose audit record values do not match with the requested record values")
		}
	} else {
		t.Error("audit record must not be nil")
	}
}

func TestAuditApi_CreateAuditRecord_error_unmarshalling(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("<>"))
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// when
	_, err := api.CreateAuditRecord(testCreateAuditRecord)

	// then
	expectedErr := generic.ClientError("Error while unmarshalling response body: invalid character '<' looking for beginning of value", "CreateAuditRecord")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}

func TestAuditApi_CreateAuditRecord_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseCreateAuditJSON))
	}))

	defer ts.Close()

	// and: the api as system under test
	api := buildAuditApi(ts.URL)

	// when
	_, err := api.CreateAuditRecord(testCreateAuditRecord)

	// then
	expectedErr := createErroneousResponseCreateAudit(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("received an unexpected error: %s", err)
	}
}
