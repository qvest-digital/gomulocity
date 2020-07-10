package deviceinformation

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tarent/gomulocity/generic"
)

type DeviceCredentials struct {
	ManagedObjects   []ManagedObject          `json:"managedObjects"`
	PagingStatistics generic.PagingStatistics `json:"statistics"`
}

type ManagedObject struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	IsDevice interface{} `json:"c8y_IsDevice"`
}

const (
	deviceCredsPath  = "/inventory/managedObjects"
	deviceCredsQuery = "?currentPage=%v"
)

func (c Client) GetDeviceInformation() (DeviceCredentials, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%v%v%v", c.BaseURL, deviceCredsPath, deviceCredsQuery),
		nil,
	)
	if err != nil {
		return DeviceCredentials{}, fmt.Errorf("failed to initialize rest request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return DeviceCredentials{}, fmt.Errorf("failed to execute rest request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DeviceCredentials{}, generic.BadCredentialsErr
	}

	deviceCredentials := DeviceCredentials{}
	if err := json.NewDecoder(resp.Body).Decode(&deviceCredentials); err != nil {
		return DeviceCredentials{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return filterDevices(deviceCredentials), nil
}

func (d DeviceCredentials) PrintToConsole() {
	for _, object := range d.ManagedObjects {
		fmt.Println(fmt.Sprintf("Device ID: %v Device name: %v", object.ID, object.Name))
	}
	fmt.Printf("Amount of devices: %v", len(d.ManagedObjects))
}

func filterDevices(deviceCreds DeviceCredentials) DeviceCredentials {
	credentials := DeviceCredentials{}

	for _, deviceCred := range deviceCreds.ManagedObjects {
		if deviceCred.IsDevice == nil {
			continue
		}
		credentials.ManagedObjects = append(credentials.ManagedObjects, deviceCred)
	}

	return credentials
}
