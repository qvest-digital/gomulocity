package managedObjects

import (
	"github.com/tarent/gomulocity/generic"
	"time"
)

type NewManagedObject struct {
	Type         string    `json:"type,omitempty"`
	Name         string    `json:"name,omitempty"`
	CreationDate time.Time `json:"creationDate"`
}

type ManagedObjectUpdate struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

type (
	ManagedObjectCollection struct {
		Self           string                    `json:"self"`
		ManagedObjects []ManagedObject           `json:"managedObjects"`
		Statistics     *generic.PagingStatistics `json:"statistics,omitempty"`
		Prev           string                    `json:"prev,omitempty"`
		Next           string                    `json:"next,omitempty"`
	}

	ManagedObject struct {
		ID           string    `json:"id"`
		Type         string    `json:"type"`
		Name         string    `json:"name"`
		CreationTime time.Time `json:"creationTime"`
		LastUpdated  time.Time `json:"lastUpdated"`
		Self         string    `json:"self"`
		Owner        string    `json:"owner"`

		AdditionParents AdditionParents `json:"additionParents"`
		AssetParents    AdditionParents `json:"assetParents"`
		DeviceParents   DeviceParents   `json:"deviceParents"`

		ChildAdditions ChildAdditions `json:"childAdditions"`
		ChildAssets    ChildAssets    `json:"childAssets"`
		ChildDevices   ChildDevices   `json:"childDevices"`

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

type ReferenceType string

const (
	CHILD_DEVICES ReferenceType = "childDevices"
	CHILD_ASSETS  ReferenceType = "childAssets"
)

type Source struct {
	Id string `json:"id"`
}

type NewManagedObjectReference struct {
	ManagedObject Source `json:"managedObject"`
}

type ManagedObjectReference struct {
	ManagedObject ManagedObject `json:"managedObject"`
	Self          string        `json:"self"`
}

type ManagedObjectReferenceCollection struct {
	Self       string                    `json:"self"`
	References []ManagedObjectReference  `json:"references"`
	Statistics *generic.PagingStatistics `json:"statistics,omitempty"`
	Prev       string                    `json:"prev,omitempty"`
	Next       string                    `json:"next,omitempty"`
}
