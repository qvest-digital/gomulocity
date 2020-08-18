package identity

import (
	"reflect"
	"testing"
)

var createExternalId = &NewExternalID{
	ExternalId: "someId",
	Type:       "someType",
}

func TestIdentity_Create_ExternalId_Success_SendsData(t *testing.T) {
	// given: A test server
	ts := buildCreateIdHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildIdentityAPI(ts.URL)
	_, err := api.CreateExternalID(*createExternalId, "123")

	if err != nil {
		t.Fatalf("CreateExternalId() got an unexpected error: %s", err.Error())
	}

	if createExternalIdCapture == nil {
		t.Fatalf("CreateExternalId() Captured ID is nil.")
	}

	if !reflect.DeepEqual(createExternalId, createExternalIdCapture) {
		t.Errorf("CreateExternalId() ID = %v, want %v", createExternalId, *createExternalIdCapture)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.externalId+json"
	if header != want {
		t.Errorf("CreateExternalId() accent header = %v, want %v", header, want)
	}
}

func TestIdentity_Create_ExternalId_Success_ReceivesData(t *testing.T) {
	// given: A test server
	ts := buildCreateIdHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildIdentityAPI(ts.URL)
	externalId, err := api.CreateExternalID(*createExternalId, "123")

	if err != nil {
		t.Fatalf("CreateExternalId() got an unexpected error: %s", err.Error())
	}

	if !reflect.DeepEqual(externalId, responseExternalId) {
		t.Errorf("CreateExternalId() ID = %v, want %v", createExternalId, createExternalIdCapture)
	}
}

func Test_Create_ExternalId_BadRequest(t *testing.T) {
	// given: A test server
	ts := buildCreateIdHttpServer(400)
	defer ts.Close()

	// and: the api as system under test
	api := buildIdentityAPI(ts.URL)
	_, err := api.CreateExternalID(*createExternalId, "123")

	if err == nil {
		t.Errorf("CreateExternalID() expected error on 400 - bad request")
		return
	}
}
