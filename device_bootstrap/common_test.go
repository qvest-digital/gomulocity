package device_bootstrap

import (
	"github.com/tarent/gomulocity/generic"
	"net/http/httptest"
	"time"
)

const (
	USER = "foo"
	PASSWORD = "bar"
)

var deviceRegistrationTime, _ = time.Parse(time.RFC3339, "2020-07-03T10:16:35.870+02:00")

func buildDeviceRegistrationApi (testServer *httptest.Server) DeviceRegistrationApi {
	c := &generic.Client{
		HTTPClient: testServer.Client(),
		BaseURL:    testServer.URL,
		Username:   USER,
		Password:   PASSWORD,
	}

	return NewDeviceRegistrationApi(c)
}
