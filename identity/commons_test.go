package identity

import (
	"net/http"
	"net/http/httptest"

	"github.com/tarent/gomulocity/generic"
)

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

var identity = `{
	"self": "https://t0818.cumulocity.com/identity",
	"externalId": "SomeExternalUrl/urltype/externalId",
	"externalIdsOfGlobalId": "someGlobalIdCollectionUrl/GlobalId/externalIds"
	}`
