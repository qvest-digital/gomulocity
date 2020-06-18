package inventory

import "testing"

func TestFragmentsFilter(t *testing.T) {
	// given
	filter := "<filter>"

	object := ManageObject{
		C8yDashboard: struct {
			Children struct {
				ID struct {
					Config struct {
						Datapoints []Datapoints
					}
				}
			}
		}{Children: struct {
			ID struct {
				Config struct {
					Datapoints []Datapoints
				}
			}
		}{ID: struct {
			Config struct {
				Datapoints []Datapoints
			}
		}{Config: struct{ Datapoints []Datapoints }{Datapoints: []Datapoints{
			{
				Fragment: "<filter>",
			},
			{
				Fragment: "<invalidFilter>",
			},
		}}}}},
	}
}
