package selfregistration

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/device_bootstrap"
	"io/ioutil"
	"os"
)

var deviceCredentialsPath = "examples/selfregistration/device_credentials.json"
var managedObjectIDPath = "examples/selfregistration/managedObject_id.json"

func storeDeviceCredentials(creds *device_bootstrap.DeviceCredentials) error {
	bytes, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("Error while saving device credentials: %s\nGiven credentials were: %#v", err, creds)
	}
	return ioutil.WriteFile(deviceCredentialsPath, bytes, os.ModePerm)
}

func readDeviceCredentials() (*device_bootstrap.DeviceCredentials, error) {
	file, err := ioutil.ReadFile(deviceCredentialsPath)
	if err != nil {
		return nil, fmt.Errorf("Error while reading stored device credentials: %s", err)
	}
	creds := &device_bootstrap.DeviceCredentials{}
	if err := json.Unmarshal(file, creds); err != nil {
		return nil, fmt.Errorf("Error while unmarshalling device creds: %s", err)
	}
	return creds, nil
}

func validateDeviceCredentials(creds *device_bootstrap.DeviceCredentials) bool {
	return creds != nil && len(creds.ID) > 0 || len(creds.TenantID) > 0 ||
		len(creds.Username) > 0 || len(creds.Password) > 0
}

func StoreManagedObjectID(ID string) error {
	temp := struct {
		ID string `json:"id"`
	}{ID: ID}

	bytes, err := json.Marshal(temp)
	if err != nil {
		return fmt.Errorf("Error while saving managed object id: %s\nGiven ID was: %v", err, ID)
	}
	return ioutil.WriteFile(managedObjectIDPath, bytes, os.ModePerm)
}
