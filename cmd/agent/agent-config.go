package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mattrax/Mattrax/mdm/agent"
)

// AgentConfigDefaultPath contains the path to store the agent configuration
const AgentConfigDefaultPath = "./agent.json"

// AgentConfig contains the agents enrollment information
type AgentConfig struct {
	agent.EnrollResponse
	PrivateKey string `json:"privateKey"`
}

func LoadConfig(path string) (AgentConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return AgentConfig{}, fmt.Errorf("error loading configuration: %w", err)
	}

	var config AgentConfig
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return AgentConfig{}, fmt.Errorf("error parsing configuration: %w", err)
	}

	f.Close()
	return config, nil
}

func SaveConfig() {
	// TODO
}
