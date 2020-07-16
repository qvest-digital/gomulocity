package inventory

import (
	"github.com/tarent/gomulocity/generic"
	"time"
)

type NewManagedObject struct {
	Type         string    `json:"type,omitempty"`
	Name         string    `json:"name,omitempty"`
	CreationTime *time.Time `json:"creationTime,omitempty"`
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
		Id           string    `json:"id"`
		Type         string    `json:"type,omitempty"`
		Name         string    `json:"name,omitempty"`
		CreationTime time.Time `json:"creationTime"`
		LastUpdated  time.Time `json:"lastUpdated"`
		Self         string    `json:"self"`
		Owner        string    `json:"owner"`

		AdditionParents AdditionParents `json:"additionParents,omitempty"`
		AssetParents    AssetParents    `json:"assetParents,omitempty"`
		DeviceParents   DeviceParents   `json:"deviceParents,omitempty"`

		ChildAdditions ChildAdditions `json:"childAdditions,omitempty"`
		ChildAssets    ChildAssets    `json:"childAssets,omitempty"`
		ChildDevices   ChildDevices   `json:"childDevices,omitempty"`

		C8YActiveAlarmsStatus   *C8YActiveAlarmsStatus   `json:"c8y_ActiveAlarmsStatus,omitempty"`
		C8YAvailability         *C8YAvailability         `json:"c8y_Availability,omitempty"`
		C8YConnection           *C8YConnection           `json:"c8y_Connection,omitempty"`
		C8YDataPoint            *C8YDataPoint            `json:"c8y_DataPoint,omitempty"`
		C8YFirmware             *C8YFirmware             `json:"c8y_Firmware,omitempty"`
		C8YHardware             *C8YHardware             `json:"c8y_Hardware,omitempty"`
		C8YIsDevice             *interface{}             `json:"c8y_IsDevice,omitempty"`
		C8YIsSensorPhone        *interface{}             `json:"c8y_IsSensorPhone,omitempty"`
		C8YRequiredAvailability *C8YRequiredAvailability `json:"c8y_RequiredAvailability,omitempty"`
		C8YSupportedOperations  *[]string                `json:"c8y_SupportedOperations,omitempty"`
		C8YStatus               *C8YStatus               `json:"c8y_Status,omitempty"`
	}

	AssetParents struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self,omitempty"`
	}
	AdditionParents struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self,omitempty"`
	}
	C8YActiveAlarmsStatus struct {
		Critical int `json:"critical,omitempty"`
		Major    int `json:"major,omitempty"`
	}
	C8YAvailability struct {
		LastMessage *time.Time `json:"lastMessage,omitempty"`
		Status      string     `json:"status,omitempty"`
	}
	C8YConnection struct {
		Status string `json:"status,omitempty"`
	}
	C8YDataPoint struct {
	}
	C8YFirmware struct {
		Version string `json:"version,omitempty"`
	}
	C8YHardware struct {
		Model        string `json:"model,omitempty"`
		SerialNumber string `json:"serialNumber,omitempty"`
	}
	C8YRequiredAvailability struct {
		ResponseInterval int `json:"responseInterval,omitempty"`
	}
	ChildAdditions struct {
		References []struct {
			ManagedObject struct {
				Id   string `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
				Self string `json:"self,omitempty"`
			} `json:"managedObject,omitempty"`
			Self string `json:"self,omitempty"`
		} `json:"references"`
		Self string `json:"self,omitempty"`
	}
	ChildAssets struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self,omitempty"`
	}
	ChildDevices struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self,omitempty"`
	}
	DeviceParents struct {
		References []interface{} `json:"references"`
		Self       string        `json:"self,omitempty"`
	}

	C8YStatus struct {
		Details struct {
			Active              int `json:"active,omitempty"`
			AggregatedResources struct {
				CPU    string `json:"cpu,omitempty"`
				Memory string `json:"memory,omitempty"`
			} `json:"aggregatedResources,omitempty"`
			Desired  int `json:"desired,omitempty"`
			Restarts int `json:"restarts,omitempty"`
		} `json:"details,omitempty"`
		Instances struct {
			DeviceSimulatorScopeManagementDeployment77678578B4Vkn66 struct {
				CPUInMillis int `json:"cpuInMillis,omitempty"`
				LastUpdated struct {
					Date struct {
						Date *time.Time `json:"date,omitempty"`
					} `json:"date,omitempty"`
					Offset int `json:"offset,omitempty"`
				} `json:"lastUpdated,omitempty"`
				MemoryInBytes int `json:"memoryInBytes,omitempty"`
				Restarts      int `json:"restarts,omitempty"`
			} `json:"device-simulator-scope-management-deployment-77678578b4-vkn66,omitempty"`
		} `json:"instances,omitempty"`
		LastUpdated struct {
			Date struct {
				Date *time.Time `json:"date,omitempty"`
			} `json:"date,omitempty"`
			Offset int `json:"offset,omitempty"`
		} `json:"lastUpdated,omitempty"`
		Status string `json:"status,omitempty"`
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
