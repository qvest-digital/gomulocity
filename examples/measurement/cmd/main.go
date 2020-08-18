package main

import (
	"github.com/tarent/gomulocity"
	"github.com/tarent/gomulocity/examples"
	exampleMeasurement "github.com/tarent/gomulocity/examples/measurement"
	"github.com/tarent/gomulocity/measurement"
	"sync"
	"time"
)

func main() {
	// Initializes a new gomulocity client
	client := gomulocity.NewGomulocity(examples.AgentConfig.BaseURL, examples.AgentConfig.Username, examples.AgentConfig.Password,
		examples.AgentConfig.BootstrapUsername, examples.AgentConfig.BootstrapPassword)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	// ID of the target managed object
	sourceID := "2278202"

	sender := exampleMeasurement.NewMeasurementSender(client, 15*time.Second, wg, measurement.Source{Id: sourceID})

	go sender.Errors()
	go sender.Fill()
	go sender.Send()
	wg.Wait()
}
