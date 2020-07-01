package gomulocity

import (
	"github.com/tarent/gomulocity/alarm"
	"github.com/tarent/gomulocity/device_bootstrap"
	"github.com/tarent/gomulocity/deviceinformation"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"time"
)

type Client struct {
	DeviceControl     device_bootstrap.Client
	DeviceInformation deviceinformation.Client
	AlarmApi          alarm.AlarmApi
}

func NewClient(baseURL, username, password string) Client {
	hc := http.Client{
		Timeout: 2 * time.Second,
	}

	client := &generic.Client{
		HTTPClient: &hc,
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
	}

	return Client{
		DeviceControl:     device_bootstrap.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
		DeviceInformation: deviceinformation.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
		AlarmApi:          alarm.NewAlarmApi(client),
	}
}
