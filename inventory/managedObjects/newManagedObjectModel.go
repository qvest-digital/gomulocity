package managedObjects

type NewManagedObject struct {
	Self         string `json:"self"`
	ID           string `json:"id"`
	LastUpdated  string `json:"lastUpdated"`
	Name         string `json:"name"`
	BinarySwitch struct {
		State string `json:"state"`
	} `json:"com_cumulocity_model_BinarySwitch"`
}

