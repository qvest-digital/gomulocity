package selfregistration

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/device_bootstrap"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"os"
	"strings"
)

var (
	deviceRegistrationCredentialsFilePath = "examples/selfregistration/registration_credentials.json"
	expectedErrSubString                  = "Device is in state PENDING_ACCEPTANCE, (not ACCEPTED)"
)

func storeDeviceRegistrationCredentials(creds *device_bootstrap.DeviceRegistration) error {
	bytes, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("Error while saving device registration credentials: %s\nGiven credentials were: %#v", err, creds)
	}
	return ioutil.WriteFile(deviceRegistrationCredentialsFilePath, bytes, os.ModePerm)
}

func readDeviceRegistrationCredentials() (*device_bootstrap.DeviceRegistration, error) {
	file, err := ioutil.ReadFile(deviceRegistrationCredentialsFilePath)
	if err != nil {
		return nil, fmt.Errorf("Error while reading stored device credentials: %s", err)
	}
	creds := &device_bootstrap.DeviceRegistration{}
	if err := json.Unmarshal(file, creds); err != nil {
		return nil, fmt.Errorf("Error while unmarshalling device registration credentials: %s", err)
	}
	return creds, nil
}

func validateDeviceRegistrationData(reg *device_bootstrap.DeviceRegistration) bool {
	return reg != nil && len(reg.Id) > 0 || len(reg.Self) > 0 ||
		len(reg.Owner) > 0 || len(reg.TenantId) > 0
}

func isExpectedErr(err *generic.Error) bool {
	return strings.Contains(err.Message, expectedErrSubString)
}
