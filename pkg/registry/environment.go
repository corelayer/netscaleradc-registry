/*
 * Copyright 2023 CoreLayer BV
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package registry

import (
	"fmt"

	"github.com/corelayer/netscaleradc-nitro-go/pkg/nitro"
)

type Environment struct {
	Name        string      `json:"name" yaml:"name"`               // Target environment name, such as "Production"
	Type        string      `json:"type" yaml:"type"`               // Target type: "StandAlone", "HighAvailabilityPair", "Cluster"
	Address     Node        `json:"address" yaml:"address"`         // Connection details for the shared Address (Address) of the environment
	Nodes       []Node      `json:"nodes" yaml:"nodes"`             // Connection details for the individual Nodes of each node
	Credentials Credentials `json:"credentials" yaml:"credentials"` // Credentials
	Settings    Settings    `json:"settings" yaml:"settings"`       // Connections settings
}

func (e *Environment) GetNodeNames() []string {
	var output []string
	for _, n := range e.Nodes {
		output = append(output, n.Name)
	}
	output = append(output, e.Address.Name)

	return output
}

// GetAllNitroClients Get a map of Client for every node in the environment (Nodes/Address)
func (e *Environment) GetAllNitroClients() (map[string]*nitro.Client, error) {
	clients := make(map[string]*nitro.Client)
	if len(e.Nodes) != 0 {
		for _, n := range e.Nodes {
			c, err := nitro.NewClient(n.GetNitroSettings(e.Settings, e.Credentials))

			if err != nil {
				return nil, fmt.Errorf("could not create client for environment %s, node %s: %w", e.Name, n.Name, err)
			}

			clients[n.Name] = c
		}
	}

	emptyNode := Node{}
	if e.Address != emptyNode {
		c, err := nitro.NewClient(e.Address.GetNitroSettings(e.Settings, e.Credentials))

		if err != nil {
			return nil, fmt.Errorf("could not create client for environment %s, Address %s: %w", e.Name, e.Address.Name, err)
		}
		clients["Address"] = c
	}

	return clients, nil
}

// GetPrimaryNodeName Get the nitro name of the primary node in the environment
func (e *Environment) GetPrimaryNodeName() (string, error) {
	var (
		err     error
		clients map[string]*nitro.Client
	)
	clients, err = e.GetAllNitroClients()
	if err != nil {
		return "", err
	}

	// Return nitro for Address if defined, as it always points to the primary node
	if _, exists := clients["Address"]; exists {
		return "Address", nil
	}

	// Return error if there are no individual nodes defined
	if len(e.Nodes) == 0 {
		return "", fmt.Errorf("invalid number of nodes defined for the environment %s (%d)", e.Name, len(e.Nodes))
	}

	// Return nitro for Nodes of the only node in a Standalone NetScaler environment
	if e.Type == "Standalone" {
		if len(e.Nodes) != 1 {
			return "", fmt.Errorf("invalid number of nodes defined for the environment %s (%d)", e.Name, len(e.Nodes))
		}
		return e.Nodes[0].Name, nil
	}

	// TODO Cluster Environment??

	// Return name of the primary node by checking the state
	for _, n := range e.Nodes {
		if _, err := clients[n.Name].IsPrimaryNode(); err == nil {
			return n.Name, nil
		} else {
			break
		}
	}
	return "", fmt.Errorf("could not detect primary node for environment %s: %w", e.Name, err)
}
