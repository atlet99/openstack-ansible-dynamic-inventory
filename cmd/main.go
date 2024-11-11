package main

import (
	"flag"
	"fmt"
	"log"
	"openstack-ansible-dynamic-inventory/pkg/inventory"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env if available
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables if available.")
	}

	// Define command-line flags
	listFlag := flag.Bool("list", false, "List all inventory")
	hostFlag := flag.String("host", "", "Get specific host details")
	flag.Parse()

	// Initialize OpenStack Inventory
	inv, err := inventory.NewOpenStackInventory()
	if err != nil {
		log.Fatalf("Failed to create inventory: %v", err)
	}

	// Handle flags
	if *listFlag {
		output, err := inv.GetInventory()
		if err != nil {
			log.Fatalf("Error getting inventory: %v", err)
		}
		fmt.Println(output)
	} else if *hostFlag != "" {
		fmt.Println("{}") // Returns empty JSON for host-specific request, similar to the Python version
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s [--list | --host <hostname>]\n", os.Args[0])
		os.Exit(1)
	}
}
