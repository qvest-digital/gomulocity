package models

type ManagedObject struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	IsDevice interface{} `json:"c8y_IsDevice"`

	C8yDashboard struct {
		Children struct {
			ID struct {
				Config struct {
					Datapoints []Datapoints `json:"datapoints"`
				} `json:"config"`
			}
		} `json:"children"`
	} `json:"c8y_Dashboard"`
}

type NewManagedObject struct {
	Self         string `json:"self"`
	ID           string `json:"id"`
	LastUpdated  string `json:"lastUpdated"`
	Name         string `json:"name"`
	BinarySwitch struct {
		State string `json:"state"`
	} `json:"com_cumulocity_model_BinarySwitch"`
}
