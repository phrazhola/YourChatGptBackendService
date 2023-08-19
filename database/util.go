package database

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func ConnectDatabase(connectionString string, databaseID string, containerID string) (*azcosmos.ContainerClient, error) {
	if connectionString == "" {
		connectionString = "YOUR_STRING"
	}

	if databaseID == "" {
		databaseID = "YOUR_ID"
	}

	if containerID == "" {
		containerID = "YOUR_CONTAINER_ID"
	}

	// Create a new Cosmos DB client
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Cosmos DB client: %v", err)
	}

	// Get a reference to the database
	database, err := client.NewDatabase(databaseID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get Cosmos DB: %v", err)
	}

	// Get a reference to the container
	container, err := database.NewContainer(containerID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get DB container: %v", err)
	}

	return container, err
}