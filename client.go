package gomulocity

import (
	"github.com/tarent/gomulocity/alarm"
	"github.com/tarent/gomulocity/inventory/managedObjects"
	"github.com/tarent/gomulocity/device_bootstrap"
	"github.com/tarent/gomulocity/generic"
	"github.com/tarent/gomulocity/measurement"
	"net/http"
	"time"
)

type Gomulocity struct {
	DeviceCredentials  device_bootstrap.DeviceCredentialsApi
	DeviceRegistration device_bootstrap.DeviceRegistrationApi
	AlarmApi           alarm.AlarmApi
	MeasurementApi    measurement.MeasurementApi
	Inventory         managedObjects.ManagedObjectApi
}

func NewGomulocity(baseURL, username, password string, bootstrapUsername, bootstrapPassword string) Gomulocity {
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
		Username:   bootstrapUsername,
		Password:   bootstrapPassword,
	}

		return Gomulocity{
		DeviceCredentials:  device_bootstrap.NewDeviceCredentialsApi(bootstrapClient),
		DeviceRegistration: device_bootstrap.NewDeviceRegistrationApi(client),
		AlarmApi:           alarm.NewAlarmApi(client),
		MeasurementApi:    measurement.NewMeasurementApi(client),
		Inventory:         managedObjects.NewManagedObjectApi(client),
	}
}
