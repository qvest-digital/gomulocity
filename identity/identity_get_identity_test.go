package identity

import (
	"testing"
)

func TestIdentity_Get_Existing_Identity(t *testing.T) {
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

	if Identity.Self == "" {
		t.Fatalf("GetIdentity() returns empty Identity")
	}
}

func TestIdentity_Get_Nonexisting_Identity(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, ``)
	defer ts.Close()

	// and: the api as system under test
	api := buildIdentityAPI(ts.URL)
	Identity, err := api.GetIdentity()

	if err == nil {
		t.Fatalf("GetIdentity() returned no Error")
	}

	if Identity != nil {
		t.Fatalf("GetIdentity() returned a wrong Identity")
	}
}
