package main

import (
	"github.com/tarent/gomulocity"
	"github.com/tarent/gomulocity/examples"
	"github.com/tarent/gomulocity/examples/selfregistration"
	"log"
)

func main() {
	// Initializes a new gomulocity client
	client := gomulocity.NewGomulocity(examples.AgentConfig.BaseURL, examples.AgentConfig.Username, examples.AgentConfig.Password,
		examples.AgentConfig.BootstrapUsername, examples.AgentConfig.BootstrapPassword)

	deviceCredentialsID, genericErr := selfregistration.SelfRegistration(client)
	if genericErr != nil {
		log.Fatal(genericErr)
	}

	managedObject, genericErr := selfregistration.CreateManagedObjectForCredentials(client, deviceCredentialsID, "GolangDevice")
	if genericErr != nil {
		log.Fatal(genericErr)
	}

	if err := selfregistration.StoreManagedObjectID(managedObject.Id); err != nil {
		log.Println(err.Error())
	}

	externalID, genericErr := selfregistration.CreateExternalID(client, deviceCredentialsID, "c8y_Serial", managedObject.Id)
	if genericErr != nil {
		log.Fatal(genericErr)
	}
	log.Println(externalID)
}
