package selfregistration

import (
	"fmt"
	"github.com/tarent/gomulocity"
	"github.com/tarent/gomulocity/device_bootstrap"
	"github.com/tarent/gomulocity/generic"
	"github.com/tarent/gomulocity/identity"
	"github.com/tarent/gomulocity/inventory"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func SelfRegistration(client gomulocity.Gomulocity) (string, *generic.Error) {
	// checks if credentials are already exist
	deviceCredentials, err := readDeviceCredentials()
	if err != nil {
		return "", generic.ClientError(fmt.Sprintf("could not read device credentials from file: %s", err), "Agent-SelfRegistration")
	}
	if validateDeviceCredentials(deviceCredentials) {
		return "", generic.ClientError(fmt.Sprintf("Device credentials are already saved: %#v", deviceCredentials), "Agent-SelfRegistration")
	}

	// reads registration credentials from file
	deviceRegistration, err := readDeviceRegistrationCredentials()
	if err != nil {
		return "", generic.ClientError(err.Error(), "Agent-SelfRegistration")
	}

	var genericErr *generic.Error
	defer func() {
		// Writes device registration data to file at the end of the function
		if genericErr == nil {
			if err := storeDeviceRegistrationCredentials(deviceRegistration); err != nil {
				log.Fatal(err)
			}
		}
	}()

	if !validateDeviceRegistrationData(deviceRegistration) {
		// generates a new random id
		deviceID := generateRandomDeviceID(100000, 999999999)

		// Creates a new device registration by given id
		deviceRegistration, genericErr = client.DeviceRegistration.Create(deviceID)
		if genericErr != nil {
			return "", genericErr
		}
	}

	// checks if a registration is available for the given deviceRegistrationID
	if _, err := registrationIsAvailable(client, deviceRegistration.Id); err != nil {
		return "", generic.ClientError(fmt.Sprintf("no registration found for device: %v, %s", deviceRegistration.Id, err), "Agent-SelfRegistration")
	}

	// Informs cumulocity that the device exists - status: Waiting for acceptance
	_, genericErr = client.DeviceCredentials.Create(deviceRegistration.Id)
	if !isExpectedErr(genericErr) {
		return "", genericErr
	}

	// Updates status to ACCEPTED
	deviceRegistration, genericErr = client.DeviceRegistration.Update(deviceRegistration.Id, device_bootstrap.ACCEPTED)
	if genericErr != nil {
		return "", genericErr
	}

	// Requests device credentials from cumulocity by using the provided deviceRegistrationID
	// If genericErr is nil, the device registration is completed
	deviceCredentials, genericErr = client.DeviceCredentials.Create(deviceRegistration.Id)
	if genericErr != nil {
		return "", genericErr
	}

	// Writes device credentials to file (device_credentials.json)
	if err = storeDeviceCredentials(deviceCredentials); err != nil {
		return "", generic.ClientError(err.Error(), "Agent-SelfRegistration")
	}
	return deviceCredentials.ID, nil
}

func registrationIsAvailable(client gomulocity.Gomulocity, ID string) (*device_bootstrap.DeviceRegistration, *generic.Error) {
	return client.DeviceRegistration.Get(ID)
}

func CreateManagedObjectForCredentials(client gomulocity.Gomulocity, ID string, typ string) (*inventory.ManagedObject, *generic.Error) {
	newManagedObject := &inventory.NewManagedObject{
		Type:         typ,
		Name:         fmt.Sprintf("Device%v", ID),
		C8y_IsDevice: nil,
	}
	return client.Inventory.Create(newManagedObject)
}

func CreateExternalID(client gomulocity.Gomulocity, ID string, typ string, objectID string) (identity.ExternalID, *generic.Error) {
	return client.Identity.CreateExternalID(identity.NewExternalID{
		ExternalId: ID,
		Type:       typ,
	}, objectID)
}

func generateRandomDeviceID(min, max int) string {
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatInt(int64(rand.Intn(max-min)+min), 10)
}
