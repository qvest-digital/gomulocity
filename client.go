package gomulocity

import (
	"github.com/tarent/gomulocity/alarm"
	"github.com/tarent/gomulocity/devicecontrol"
	"github.com/tarent/gomulocity/deviceinformation"
	"net/http"
	"time"
)

type Client struct {
	DeviceControl     devicecontrol.Client
	DeviceInformation deviceinformation.Client
	Alarm             alarm.Client
}

func NewClient(baseURL, username, password string) Client {
	hc := http.Client{
		Timeout: 2 * time.Second,
	}

	return Client{
		DeviceControl:     devicecontrol.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
		DeviceInformation: deviceinformation.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
		Alarm:             alarm.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
	}
}
