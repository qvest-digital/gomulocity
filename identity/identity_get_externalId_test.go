package identity

import (
	"testing"
)

func TestIdentity_Get_Existing_ExternalId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, externalID)
	defer ts.Close()

	// and: the api as system under test
	api := buildIdentityAPI(ts.URL)

	receivedID, err := api.GetExternalID("someType", "someId")

	if err != nil {
		t.Fatalf("GetExternalId() got an unexpected error: %s", err.Error())
	}

	if receivedID == nil {
		t.Fatalf("GetExternalId() returns nil.")
	}

	if receivedID.Self == "" {
		t.Fatalf("GetExternalId() returns empty Identity")
	}
}

func TestIdentity_Get_Nonexisting_ExternalId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, ``)
	defer ts.Close()

	// and: the api as system under test
	api := buildIdentityAPI(ts.URL)
	receivedId, err := api.GetExternalID("someType", "someNonextistentId")

	if err == nil {
		t.Fatalf("GetExternalId() returned no Error")
	}

	if receivedId != nil {
		t.Fatalf("GetExternalId() returned a wrong Identity")
	}
}
