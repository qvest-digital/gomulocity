package gomulocity

import (
	"github.com/tarent/gomulocity/devicecontrol"
	"github.com/tarent/gomulocity/deviceinformation"
	"github.com/tarent/gomulocity/inventory"
	"net/http"
	"time"
)

type Client struct {
	DeviceControl     devicecontrol.Client
	DeviceInformation deviceinformation.Client
	Inventory         inventory.Client
}

func NewClient(baseURL, username, password string) Client {
	hc := http.Client{
		Timeout: 2 * time.Second,
	}

	return Client{
		DeviceControl:     devicecontrol.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
		DeviceInformation: deviceinformation.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
		Inventory:         inventory.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
	}
}
