package models

type (
	InventoryStructure struct {
		ManagedObject []ManagedObject `json:"managedObjects"`
	}

	ManagedObject struct {
		ID             string `json:"id"`
		Owner          string `json:"owner"`
		C8YTemperature struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"c8y_Temperature,omitempty"`
	}
)

type NewManagedObject struct {
	Self         string `json:"self"`
	ID           string `json:"id"`
	LastUpdated  string `json:"lastUpdated"`
	Name         string `json:"name"`
	BinarySwitch struct {
		State string `json:"state"`
	} `json:"com_cumulocity_model_BinarySwitch"`
}
