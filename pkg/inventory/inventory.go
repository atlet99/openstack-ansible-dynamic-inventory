package inventory

import (
	"fmt"
	"log"
	"openstack-ansible-dynamic-inventory/pkg/utils"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type OpenStackInventory struct {
	Inventory        map[string]interface{}
	Conn             *gophercloud.ServiceClient
	EnvironmentTag   string
	EnvironmentValue string
	BaseGroupName    string
}

// NewOpenStackInventory initializes OpenStack inventory with required environment variables
func NewOpenStackInventory() (*OpenStackInventory, error) {
	envTag := os.Getenv("ENVIRONMENT_TAG")
	envValue := os.Getenv("ENVIRONMENT_VALUE")
	baseGroup := os.Getenv("BASE_GROUP_NAME")

	// Check if essential environment variables are set
	if envTag == "" || envValue == "" || baseGroup == "" {
		return nil, fmt.Errorf("missing required environment variables: ENVIRONMENT_TAG, ENVIRONMENT_VALUE, or BASE_GROUP_NAME")
	}

	client, err := connectToOpenStack()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OpenStack: %w", err)
	}

	inventory := map[string]interface{}{
		"_meta": map[string]interface{}{
			"hostvars": map[string]interface{}{},
		},
		baseGroup: map[string]interface{}{
			"hosts": []string{},
			"vars": map[string]interface{}{
				"environment_tag":   envTag,
				"environment_value": envValue,
			},
		},
	}

	return &OpenStackInventory{
		Inventory:        inventory,
		Conn:             client,
		EnvironmentTag:   envTag,
		EnvironmentValue: envValue,
		BaseGroupName:    baseGroup,
	}, nil
}

// shouldIncludeServer checks if a server should be included based on metadata
func (oi *OpenStackInventory) shouldIncludeServer(server servers.Server) bool {
	return server.Metadata[oi.EnvironmentTag] == oi.EnvironmentValue
}

// AddHostToGroups adds the host to the base group and creates metadata-based groups as needed
func (oi *OpenStackInventory) AddHostToGroups(hostname string, metadata map[string]string) {
	if oi.Inventory[oi.BaseGroupName] != nil {
		hosts := oi.Inventory[oi.BaseGroupName].(map[string]interface{})["hosts"].([]string)
		oi.Inventory[oi.BaseGroupName].(map[string]interface{})["hosts"] = append(hosts, hostname)
	}

	// Create metadata-based groups dynamically
	for key, value := range metadata {
		if key == oi.EnvironmentTag {
			continue
		}
		groupName := fmt.Sprintf("%s_%s", key, value)
		if _, exists := oi.Inventory[groupName]; !exists {
			oi.Inventory[groupName] = map[string]interface{}{
				"hosts": []string{},
				"vars": map[string]interface{}{
					"group_tag":   key,
					"group_value": value,
				},
			}
		}
		oi.Inventory[groupName].(map[string]interface{})["hosts"] = append(oi.Inventory[groupName].(map[string]interface{})["hosts"].([]string), hostname)
	}
}

// GetInventory processes servers from OpenStack and builds the inventory structure
func (oi *OpenStackInventory) GetInventory() (string, error) {
	pages, err := servers.List(oi.Conn, servers.ListOpts{}).AllPages()
	if err != nil {
		return "", fmt.Errorf("failed to list servers: %w", err)
	}

	allServers, err := servers.ExtractServers(pages)
	if err != nil {
		return "", fmt.Errorf("failed to extract servers: %w", err)
	}

	for _, server := range allServers {
		if !oi.shouldIncludeServer(server) {
			continue
		}

		// Get preferred IPv4 address
		ip, err := oi.getPreferredIPv4(server.Addresses)
		if err != nil {
			log.Printf("Warning: No suitable IP found for server %s, skipping...", server.Name)
			continue
		}

		// Add the host to the inventory and groups
		oi.AddHostToGroups(server.Name, server.Metadata)
		oi.Inventory["_meta"].(map[string]interface{})["hostvars"].(map[string]interface{})[server.Name] = map[string]interface{}{
			"ansible_host":       ip,
			"openstack_id":       server.ID,
			"openstack_name":     server.Name,
			"openstack_metadata": server.Metadata,
		}
	}

	return utils.JSONFormat(oi.Inventory)
}

// getPreferredIPv4 returns the first available IPv4 address from the server's networks
func (oi *OpenStackInventory) getPreferredIPv4(addresses map[string]interface{}) (string, error) {
	for _, addrs := range addresses {
		for _, addr := range addrs.([]interface{}) {
			a := addr.(map[string]interface{})
			if a["version"] == 4 {
				return a["addr"].(string), nil
			}
		}
	}
	return "", fmt.Errorf("no IPv4 address found")
}
