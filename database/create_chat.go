package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"YourChatGptBackendService/models/api"
	"YourChatGptBackendService/models/data"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

func CreateChat(email string, title string) (string, string, error) {
	container, err := ConnectDatabase("", "", "")
	if err != nil {
		return "", "", fmt.Errorf("Failed to connect to the database: %v", err)
	}

	createTime := time.Now()
	updateTime := time.Now()

	// Generate a new chat ID
	chatID := uuid.New().String()
	chat := models.Chat{
		ChatID: chatID,
		Title: title,
		CreateTime: createTime,
		UpdateTime: updateTime,
	}

	// If the user doesn't exist yet in the database, user id == ""
	userID, err := GetUser(email, container)
	if err != nil {
		return "", "", fmt.Errorf("Failed to fetch user: %v", err)
	}

	// Create a new chat data
	ref := data.BackRef{
		UserID: userID,
	}
	chatData := data.ChatData{
		ID: chatID,
		ChatID:     chatID,
		Title:      title,
		CreateTime: createTime,
		UpdateTime: updateTime,
		ChatHistory: data.ChatHistory{
			MessagesMapping: map[string]data.MessageThreadData{},
			LastUpdatedMessage: "",
		},
		BackRef: ref,
		Partition: userID,
	}

	if userID != "" {
		itemResponse, err := container.ReadItem(context.Background(), azcosmos.NewPartitionKeyString(userID), userID, nil)
		if err != nil {
			return "", "", fmt.Errorf("Failed to fetch item response: %v", err)
		}
		var userData data.UserData
		var itemResponseBody map[string]interface{}
		err = json.Unmarshal(itemResponse.Value, &itemResponseBody)

		chats, ok := itemResponseBody["chats"].([]interface{})

		// retrieve chat items in the chat list
		if ok {
			for _, itemData := range chats {
				item := itemData.(map[string]interface{})
				chatItem := models.Chat{
					ChatID: item["chat_id"].(string),
					Title:  item["title"].(string),
					CreateTime: func() time.Time {
						t, _ := time.Parse(time.RFC3339, item["create_time"].(string))
						return t
					}(),
					UpdateTime: func() time.Time {
						t, _ := time.Parse(time.RFC3339, item["update_time"].(string))
						return t
					}(),
				}
				userData.Chats = append(userData.Chats, chatItem)
			}
		} else {
			return "", "", fmt.Errorf("Failed to parse chats: %v", err)
		}

		userData.ID = userID
		userData.UserID = userID
		userData.Total = int(itemResponseBody["total"].(float64)) + 1
		userData.Partition = userID

		userData.Chats = append(userData.Chats, chat)

		userDataToReplace, err := json.Marshal(userData)
		if err != nil {
			return "", "", fmt.Errorf("Error: %v", err)
		}

		// update the user chats data in database
		itemResponse, err = container.ReplaceItem(context.Background(), azcosmos.NewPartitionKeyString(userID), userID, userDataToReplace, nil)
		if err != nil {
			return "", "", fmt.Errorf("Failed to replace item with id: %v", userID)
		}

		fmt.Println("User chat item updated")

		// create new chat data in database
		err = CreateNewChatData(chatData, container, userID)
		if err != nil {
			return "", "", fmt.Errorf("Failed to create new chat data: %v", err)
		}

		return userID, chatID, nil
	} else {
		// User ID does not exist, create a new user chats data
		userID, err := CreateUser(email, container)
		if err != nil {
			return "", "", fmt.Errorf("Error: %v", err)
		}

		chatData.BackRef.UserID = userID
		chatData.Partition = userID
		
		userData := data.UserData{
			ID: userID,
			UserID: userID,
			Chats:  []models.Chat{chat},
			Total:  1,
			Partition: userID,
		}

		userDataToCreate, err := json.Marshal(userData)
		if err != nil {
			return "", "", fmt.Errorf("Error: %v", err)
		}

		// create a new user chats data in the database
		_, err = container.CreateItem(context.TODO(), azcosmos.NewPartitionKeyString(userID), userDataToCreate, nil)
		if err != nil {
			return "", "", fmt.Errorf("Failed to create user chat item: %v", err)
		}

		fmt.Println("New user chat item created")

		// create new chat data in database
		err = CreateNewChatData(chatData, container, userID)
		if err != nil {
			return "", "", fmt.Errorf("Failed to create new chat data: %v", err)
		}

		return userID, chatID, nil
	}
}

func CreateNewChatData(chatData data.ChatData, container *azcosmos.ContainerClient, userID string) (error) {
	chatDataToCreate, err := json.Marshal(chatData)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	_, err = container.CreateItem(context.Background(), azcosmos.NewPartitionKeyString(userID), chatDataToCreate, nil)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	fmt.Println("New chat item created")

	return nil
}
