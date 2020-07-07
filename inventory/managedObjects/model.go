package managedObjects

import (
	"github.com/tarent/gomulocity/generic"
	"time"
)

const (
	MANAGED_OBJECT_ACCEPT       = "application/vnd.com.nsn.cumulocity.managedObject+json"
	MANAGED_OBJECT_CONTENT_TYPE = "application/vnd.com.nsn.cumulocity.managedObject+json"

	managedObjectPath = "/inventory/managedObjects"
)

type Update struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type UpdateResponse struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Self             string       `json:"self"`
	Type             string       `json:"type"`
	LastUpdated      time.Time    `json:"lastUpdated"`
	StrongTypedClass struct{}     `json:"com_othercompany_StrongTypedClass"`
	ChildDevices     ChildDevices `json:"childDevices"`
}

type Reference struct {
	ManagedObject ManagedObject `json:"managedObject"`
	Self          string        `json:"self"`
}

type ReferenceCollection struct {
	Next       string      `json:"next"`
	Self       string      `json:"self"`
	References []Reference `json:"references"`
}

type NewManagedObject struct {
	ID           string    `json:"id"`
	Self         string    `json:"self"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creationDate"`
	LastUpdated  time.Time `json:"lastUpdated"`
	BinarySwitch struct {
		State string `json:"state"`
	} `json:"com_cumulocity_model_BinarySwitch"`
}

type CreateManagedObject struct {
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creationDate"`
}

type (
	ManagedObjectCollection struct {
		Next           string                   `json:"next"`
		Self           string                   `json:"self"`
		ManagedObjects []ManagedObject          `json:"managedObjects"`
		Prev           string                   `json:"prev"`
		Statistics     generic.PagingStatistics `json:"statistics"`
	}

	ManagedObject struct {
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
		ID                      string                  `json:"id"`
		LastUpdated             time.Time               `json:"lastUpdated"`
		Name                    string                  `json:"name"`
		Owner                   string                  `json:"owner"`
		Self                    string                  `json:"self"`
		Type                    string                  `json:"type"`
		CreationTime            time.Time               `json:"creationTime"`
		DeviceParents           DeviceParents           `json:"deviceParents"`
		C8YStatus               C8YStatus               `json:"c8y_Status"`
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
	DeviceParents struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self"`
	}

	C8YStatus struct {
		Details struct {
			Active              int `json:"active"`
			AggregatedResources struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"aggregatedResources"`
			Desired  int `json:"desired"`
			Restarts int `json:"restarts"`
		} `json:"details"`
		Instances struct {
			DeviceSimulatorScopeManagementDeployment77678578B4Vkn66 struct {
				CPUInMillis int `json:"cpuInMillis"`
				LastUpdated struct {
					Date struct {
						Date time.Time `json:"$date"`
					} `json:"date"`
					Offset int `json:"offset"`
				} `json:"lastUpdated"`
				MemoryInBytes int `json:"memoryInBytes"`
				Restarts      int `json:"restarts"`
			} `json:"device-simulator-scope-management-deployment-77678578b4-vkn66"`
		} `json:"instances"`
		LastUpdated struct {
			Date struct {
				Date time.Time `json:"$date"`
			} `json:"date"`
			Offset int `json:"offset"`
		} `json:"lastUpdated"`
		Status string `json:"status"`
	}
)
