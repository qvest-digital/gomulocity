package identity

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/tarent/gomulocity/deviceinformation"
	"github.com/tarent/gomulocity/generic"
)

var identity = `{
	"self": "https://t0818.cumulocity.com/identity",
	"externalId": "SomeExternalUrl/urltype/externalId",
	"externalIdsOfGlobalId": "someGlobalIdCollectionUrl/GlobalId/externalIds"
	}`

var externalID = `{
	"self": "selfUrl",
	"externalId": "someExternalId",
	"type": "someType",
	"managedObject":{

	}
	}`

var requestCapture *http.Request
var createExternalIdCapture *ExternalID
var responseExternalId = ExternalID{
	Self:          "someSelfAssignedURL",
	ExternalId:    "someId",
	Type:          "someType",
	ManagedObject: deviceinformation.ManagedObject{},
}

func buildIdentityAPI(url string) IdentityAPI {
	httpClient := http.DefaultClient
	client := generic.Client{
		HTTPClient: httpClient,
		BaseURL:    url,
		Username:   "foo",
		Password:   "bar",
	}
	return NewIdentityAPI(client)
}

func buildHttpServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func buildCreateIdHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		var createExternalId ExternalID
		_ = json.Unmarshal(body, &createExternalId)
		createExternalIdCapture = &createExternalId
		requestCapture = r

		w.WriteHeader(status)
		response, _ := json.Marshal(responseExternalId)
		_, _ = w.Write(response)
	}))
}
