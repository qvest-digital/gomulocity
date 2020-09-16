package measurement

import (
	"log"
	"time"

	"github.com/tarent/gomulocity/measurement"
)

var Example1NewMeasurements = measurement.NewMeasurement{
	Time:            timeToPointer(time.Now().Format(time.RFC3339)),
	MeasurementType: "P",
	Metrics: map[string]interface{}{
		"P": struct {
			P struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			}
		}{
			P: struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			}{Unit: "W", Value: 71},
		},
	},
}

var Example2NewMeasurements = measurement.NewMeasurement{
	Time:            timeToPointer(time.Now().Format(time.RFC3339)),
	MeasurementType: "P",
	Metrics: map[string]interface{}{
		"P": struct {
			P1 struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			}
			P2 struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			}
			P3 struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			}
		}{
			P1: struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			}{Unit: "W", Value: 77},
			P2: struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			}{Unit: "W", Value: 43},
			P3: struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			}{Unit: "W", Value: 12},
		},
	},
}

var ExampleCollection = []measurement.NewMeasurement{
	{
		Time:            timeToPointer(time.Now().Format(time.RFC3339)),
		MeasurementType: "VoltageMeasurement",
		Metrics: map[string]interface{}{
			"VoltageMeasurement": struct {
				Voltage struct {
					Unit  string  `json:"unit"`
					Value float64 `json:"value"`
				} `json:"voltage"`
			}{
				struct {
					Unit  string  `json:"unit"`
					Value float64 `json:"value"`
				}{Unit: "V", Value: 227.32},
			},
		},
	},
	{
		Time:            timeToPointer(time.Now().Format(time.RFC3339)),
		MeasurementType: "c8y_FrequencyMeasurement",
		Metrics: map[string]interface{}{
			"c8y_FrequencyMeasurement": struct {
				Frequency struct {
					Unit  string
					Value float64
				}
			}{
				struct {
					Unit  string
					Value float64
				}{Unit: "Hz", Value: 37.71},
			},
		},
	},
	{
		Time:            timeToPointer(time.Now().Format(time.RFC3339)),
		MeasurementType: "P",
		Metrics: map[string]interface{}{
			"P": struct {
				P1 struct {
					Unit  string `json:"unit"`
					Value int    `json:"value"`
				}
				P2 struct {
					Unit  string `json:"unit"`
					Value int    `json:"value"`
				}
				P3 struct {
					Unit  string `json:"unit"`
					Value int    `json:"value"`
				}
			}{
				P1: struct {
					Unit  string `json:"unit"`
					Value int    `json:"value"`
				}{Unit: "W", Value: 77},
				P2: struct {
					Unit  string `json:"unit"`
					Value int    `json:"value"`
				}{Unit: "W", Value: 43},
				P3: struct {
					Unit  string `json:"unit"`
					Value int    `json:"value"`
				}{Unit: "W", Value: 12},
			},
		},
	},
}

func timeToPointer(timeString string) *time.Time {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		log.Fatal(err)
	}
	return &t
}
