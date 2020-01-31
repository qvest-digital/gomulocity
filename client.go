package gomulocity

import (
	"github.com/tarent/gomulocity/devicecontrol"
	"net/http"
	"time"
)

type Client struct {
	DeviceControl devicecontrol.Client
}

func NewClient(baseURL, username, password string) Client {
	hc := http.Client{
		Timeout: 2 * time.Second,
	}

	return Client{
		DeviceControl: devicecontrol.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
	}
}
