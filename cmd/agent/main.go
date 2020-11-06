package main

import (
	"fmt"
	"os"

	"github.com/mattrax/Mattrax/mdm/agent"
)

// Version contains the version of the agent being used
const Version = "v1.0.0-dev"

// AgentConfigPath contains the path to store the agent configuration
const AgentConfigPath = "./agent.json"

// AgentConfig contains the agents enrollment information
type AgentConfig struct {
	agent.EnrollResponse
	PrivateKey string `json:"privateKey"`
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Please specify an email as the first argument!")
		return
	}

	if err := enroll(os.Args[1]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
