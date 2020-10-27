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

func SelfRegistration(client gomulocity.Gomulocity, timer time.Duration) (*device_bootstrap.DeviceCredentials, *generic.Error) {
	// checks if credentials are already exist
	deviceCredentials, err := readDeviceCredentials()
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("could not read device credentials from file: %s", err), "Agent-SelfRegistration")
	}
	if validateDeviceCredentials(deviceCredentials) {
		return nil, generic.ClientError(fmt.Sprintf("Device credentials are already saved: %#v", deviceCredentials), "Agent-SelfRegistration")
	}

	var deviceRegistration *device_bootstrap.DeviceRegistration
	var genericErr *generic.Error

	// reads registration credentials from file
	deviceRegistration, err = readDeviceRegistrationCredentials()
	if err != nil {
		return nil, generic.ClientError(err.Error(), "Agent-SelfRegistration")
	}

	if !validateDeviceRegistrationData(deviceRegistration) {
		// generates a new random id
		deviceID := generateRandomDeviceID(100000, 999999999)

		for {
			// Requests device credentials from cumulocity by using the provided deviceID
			// If genericErr is nil and deviceCredentials is not nil, the device registration is completed
			deviceCredentials, genericErr = client.DeviceCredentials.Create(deviceID)
			if genericErr != nil {
				log.Println(genericErr.Error())
			}
			if deviceCredentials != nil {
				break
			}
			time.Sleep(timer * time.Second)
		}
	}

	defer func() {
		// Writes device registration data to file at the end of the function
		if genericErr == nil {
			if err := storeDeviceRegistrationCredentials(deviceRegistration); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// Writes device credentials to file (device_creds.json)
	if err = storeDeviceCredentials(deviceCredentials); err != nil {
		return nil, generic.ClientError(err.Error(), "Agent-SelfRegistration")
	}
	return deviceCredentials, nil
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
