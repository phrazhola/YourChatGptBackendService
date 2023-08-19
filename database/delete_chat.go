package database

import (
	"context"
	"encoding/json"
	"fmt"
	"YourChatGptBackendService/models/api"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func DeleteChat(userID string, chatID string) (error) {
	container, err := ConnectDatabase("", "", "")
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %v", err)
	}

	_, err = container.DeleteItem(context.Background(), azcosmos.NewPartitionKeyString(userID), chatID, nil)
	if err != nil {
		return fmt.Errorf("Failed to delete item: %v", err)
	}

	userChatsData, err := GetUserChatsData(userID)
	if err != nil {
		return fmt.Errorf("Failed to get user chats data: %v", err)
	}
	chats := userChatsData.Chats

	filteredChats := []models.Chat{}

	for _, chat := range chats {
		if chat.ChatID != chatID {
			filteredChats = append(filteredChats, chat)
		}
	}

	userChatsData.Chats = filteredChats
	userChatsData.Total--

	userDataToReplace, err := json.Marshal(userChatsData)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	// update the user chats data in database
	_, err = container.ReplaceItem(context.Background(), azcosmos.NewPartitionKeyString(userID), userID, userDataToReplace, nil)
	if err != nil {
		return fmt.Errorf("Failed to replace item with id: %v", userID)
	}

	fmt.Println("User chat item updated")
	
	return nil
}
