package examples

import (
	"log"
	"os"
)

type Config struct {
	BaseURL           string
	Username          string
	Password          string
	BootstrapUsername string
	BootstrapPassword string
}

var AgentConfig *Config

func init() {
	baseURL := os.Getenv("C8Y_AGENT_BASE_URL")
	if len(baseURL) == 0 {
		log.Fatal("Environment variable 'C8Y_AGENT_BASE_URL' must be set")
	}
	username := os.Getenv("C8Y_AGENT_USERNAME")
	if len(username) == 0 {
		log.Fatal("Environment variable 'C8Y_AGENT_USERNAME' must be set")
	}
	password := os.Getenv("C8Y_AGENT_PASSWORD")
	if len(password) == 0 {
		log.Fatal("Environment variable 'C8Y_AGENT_PASSWORD' must be set")
	}
	bootstrapUsername := os.Getenv("C8Y_AGENT_BOOTSTRAP_USERNAME")
	if len(bootstrapUsername) == 0 {
		log.Fatal("Environment variable 'C8Y_AGENT_BOOTSTRAP_USERNAME' must be set")
	}
	bootstrapPassword := os.Getenv("C8Y_AGENT_BOOTSTRAP_PASSWORD")
	if len(bootstrapPassword) == 0 {
		log.Fatal("Environment variable 'C8Y_AGENT_BOOTSTRAP_PASSWORD' must be set")
	}

	AgentConfig = &Config{
		BaseURL:           baseURL,
		Username:          username,
		Password:          password,
		BootstrapUsername: bootstrapUsername,
		BootstrapPassword: bootstrapPassword,
	}
}
