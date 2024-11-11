package inventory

import (
	"fmt"
	"log"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// connectToOpenStack establishes a connection to OpenStack using environment variables
func connectToOpenStack() (*gophercloud.ServiceClient, error) {
	// Set up authentication options from environment variables
	authOptions, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenStack auth options: %w", err)
	}

	// Attempt to authenticate with OpenStack
	provider, err := openstack.AuthenticatedClient(authOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with OpenStack: %w", err)
	}

	// Connect to the OpenStack Compute (Nova) service
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"), // Using region from environment variables
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Compute V2 client: %w", err)
	}

	log.Println("Successfully connected to OpenStack Compute service")
	return client, nil
}

// GetHosts retrieves server data from OpenStack, organizing it into groups
func (oi *OpenStackInventory) GetHosts() error {
	// Retrieve all pages of servers
	pages, err := servers.List(oi.Conn, servers.ListOpts{}).AllPages()
	if err != nil {
		return fmt.Errorf("failed to retrieve servers list: %w", err)
	}

	// Extract servers from retrieved pages
	allServers, err := servers.ExtractServers(pages)
	if err != nil {
		return fmt.Errorf("failed to extract servers: %w", err)
	}

	// Process each server and organize it into the inventory structure
	for _, server := range allServers {
		for networkName, addresses := range server.Addresses {
			for _, address := range addresses.([]interface{}) {
				addr := address.(map[string]interface{})
				if addr["version"] == 4 { // Only include IPv4 addresses
					// Create a unique hostname with network suffix
					hostName := fmt.Sprintf("%s_%s", server.Name, networkName)

					// Add server variables to the inventory
					serverVars := map[string]interface{}{
						"ansible_host":     addr["addr"],
						"openstack_id":     server.ID,
						"openstack_name":   server.Name,
						"openstack_status": server.Status,
						"network_name":     networkName,
					}

					// Add the server to the groups based on metadata
					metadata := make(map[string]string)
					for k, v := range server.Metadata {
						metadata[k] = v
					}
					oi.AddHostToGroups(hostName, metadata)

					// Store the server variables in the _meta hostvars section
					oi.Inventory["_meta"].(map[string]interface{})["hostvars"].(map[string]interface{})[hostName] = serverVars
				}
			}
		}
	}
	log.Println("Successfully organized servers into inventory groups")
	return nil
}
