package managedObjects

import (
	"github.com/tarent/gomulocity/generic"
	"time"
)

const (
	managedObjectPath = "/inventory/managedObjects"
)

type (
	ManagedObjectCollection struct {
		Next           string                   `json:"next"`
		Self           string                   `json:"self"`
		ManagedObjects []ManagedObjects         `json:"managedObjects"`
		Prev           string                   `json:"prev"`
		Statistics     generic.PagingStatistics `json:"statistics"`
	}

	ManagedObjects struct {
		AdditionParents         AdditionParents         `json:"additionParents"`
		AssetParents            AdditionParents         `json:"assetParents"`
		C8YActiveAlarmsStatus   C8YActiveAlarmsStatus   `json:"c8y_ActiveAlarmsStatus"`
		C8YAvailability         C8YAvailability         `json:"c8y_Availability"`
		C8YConnection           C8YConnection           `json:"c8y_Connection"`
		C8YDataPoint            C8YDataPoint            `json:"c8y_DataPoint"`
		C8YFirmware             C8YFirmware             `json:"c8y_Firmware"`
		C8YHardware             C8YHardware             `json:"c8y_Hardware"`
		C8YIsDevice             interface{}             `json:"c8y_IsDevice"`
		C8YIsSensorPhone        interface{}             `json:"c8y_IsSensorPhone"`
		C8YRequiredAvailability C8YRequiredAvailability `json:"c8y_RequiredAvailability"`
		C8YSupportedOperations  []string                `json:"c8y_SupportedOperations"`
		ChildAdditions          ChildAdditions          `json:"childAdditions"`
		ChildAssets             ChildAssets             `json:"childAssets"`
		ChildDevices            ChildDevices            `json:"childDevices"`
		ComCumulocityModelAgent ComCumulocityModelAgent `json:"com_cumulocity_model_Agent"`
		ID                      string                  `json:"id"`
		LastUpdated             time.Time               `json:"lastUpdated"`
		Name                    string                  `json:"name"`
		Owner                   string                  `json:"owner"`
		Self                    string                  `json:"self"`
		Type                    string                  `json:"type"`
		CreationTime            time.Time               `json:"creationTime"`
		DeviceParents           DeviceParents           `json:"deviceParents"`
	}

	AssetParents struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self"`
	}
	AdditionParents struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self"`
	}
	C8YActiveAlarmsStatus struct {
		Critical int `json:"critical"`
		Major    int `json:"major"`
	}
	C8YAvailability struct {
		LastMessage time.Time `json:"lastMessage"`
		Status      string    `json:"status"`
	}
	C8YConnection struct {
		Status string `json:"status"`
	}
	C8YDataPoint struct {
	}
	C8YFirmware struct {
		Version string `json:"version"`
	}
	C8YHardware struct {
		Model        string `json:"model"`
		SerialNumber string `json:"serialNumber"`
	}
	C8YRequiredAvailability struct {
		ResponseInterval int `json:"responseInterval"`
	}
	ChildAdditions struct {
		References []struct {
			ManagedObject struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Self string `json:"self"`
			} `json:"managedObject"`
			Self string `json:"self"`
		} `json:"references"`
		Self string `json:"self"`
	}
	ChildAssets struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self"`
	}
	ChildDevices struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self"`
	}
	ComCumulocityModelAgent struct {
	}
	DeviceParents struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self"`
	}
)
