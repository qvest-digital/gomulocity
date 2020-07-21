package identity

import "testing"

func TestEvents_Get_Existing_Identity(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, identity)
	defer ts.Close()

	// and: the api as system under test
	api := buildIdentityAPI(ts.URL)

	Identity, err := api.GetIdentity()

	if err != nil {
		t.Fatalf("GetIdentity() got an unexpected error: %s", err.Error())
	}

	if Identity == nil {
		t.Fatalf("GetIdentity() returns nil.")
	}

	if Identity.self == "" {
		t.Fatalf("GetIdentity() returns empty Identity")
	}
}
