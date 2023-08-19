package database

import (
	"context"
	"encoding/json"
	"fmt"

	"YourChatGptBackendService/models/data"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

// email is encoded email string
func CreateUser(email string, container *azcosmos.ContainerClient) (string, error) {
	if container == nil {
		ctn, err := ConnectDatabase("", "", "")
		if err != nil {
			return "", err
		}

		container = ctn
	}

	userID := uuid.New().String()
	partition := string(email[0]);

	user := data.User{
		ID: email,
		UserID: userID,
		Partition: partition,
	}

	userToCreate, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	_, err = container.CreateItem(context.Background(), azcosmos.NewPartitionKeyString(partition), userToCreate, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create user: %v", err)
	}

	fmt.Println("New user created")

	return userID, nil
}

