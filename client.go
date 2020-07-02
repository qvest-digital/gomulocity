package gomulocity

import (
	"github.com/tarent/gomulocity/alarm"
	"github.com/tarent/gomulocity/devicecontrol"
	"github.com/tarent/gomulocity/generic"
	"github.com/tarent/gomulocity/inventory/managedObjects"
	"net/http"
	"time"
)

type Client struct {
	DeviceControl     devicecontrol.Client
	AlarmApi          alarm.AlarmApi
	Inventory         managedObjects.ManagedObjectApi
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
		DeviceControl:     devicecontrol.Client{HTTPClient: &hc, BaseURL: baseURL, Username: username, Password: password},
		AlarmApi:          alarm.NewAlarmApi(client),
		Inventory:         managedObjects.NewManagedObjectApi(client),
	}
}
