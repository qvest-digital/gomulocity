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
	DeviceCredentials  device_bootstrap.DeviceCredentialsApi
	DeviceRegistration device_bootstrap.DeviceRegistrationApi
	DeviceInformation  deviceinformation.Client
	AlarmApi           alarm.AlarmApi
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

	bootstrapClient := &generic.Client{
		HTTPClient: &hc,
		BaseURL:    baseURL,
		Username:   "devicebootstrap",
		Password:   "Fhdt1bb1f",
	}

	return Client{
		DeviceCredentials:  device_bootstrap.NewDeviceCredentialsApi(bootstrapClient),
		DeviceRegistration: device_bootstrap.NewDeviceRegistrationApi(client),
		DeviceInformation:  deviceinformation.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
		AlarmApi:           alarm.NewAlarmApi(client),
	}
}
