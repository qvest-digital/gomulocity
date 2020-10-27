package main

import (
	"github.com/tarent/gomulocity"
	"github.com/tarent/gomulocity/examples"
	"github.com/tarent/gomulocity/examples/selfregistration"
	"log"
	"time"
)

func main() {
	// Initializes a new gomulocity client
	examples.InitConfig()

	client := gomulocity.NewGomulocity(examples.AgentConfig.BaseURL, examples.AgentConfig.Username, examples.AgentConfig.Password,
		examples.AgentConfig.BootstrapUsername, examples.AgentConfig.BootstrapPassword)

	timer := 15 * time.Second
	deviceCredentials, genericErr := selfregistration.SelfRegistration(client, timer)
	if genericErr != nil {
		log.Fatal(genericErr)
	}

	managedObject, genericErr := selfregistration.CreateManagedObjectForCredentials(client, deviceCredentials.ID, "GolangDevice")
	if genericErr != nil {
		log.Fatal(genericErr)
	}

	if err := selfregistration.StoreManagedObjectID(managedObject.Id); err != nil {
		log.Println(err.Error())
	}

	externalID, genericErr := selfregistration.CreateExternalID(client, deviceCredentials.ID, "c8y_Serial", managedObject.Id)
	if genericErr != nil {
		log.Fatal(genericErr)
	}
	log.Println(externalID)
}
