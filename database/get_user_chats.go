package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"YourChatGptBackendService/models/api"
	"YourChatGptBackendService/models/data"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func GetUserChats(email string) (models.GetChatsResponse, error) {
	emptyRes := models.GetChatsResponse{}

	container, err := ConnectDatabase("", "", "")
	if err != nil {
		return emptyRes, fmt.Errorf("Failed to connect to the database: %v", err)
	}

	// If the user doesn't exist yet in the database, user id == ""
	userID, err := GetUser(email, container)
	if err != nil {
		return emptyRes, fmt.Errorf("Failed to fetch user: %v", err)
	}

	if userID != "" {
		itemResponse, err := container.ReadItem(context.Background(), azcosmos.NewPartitionKeyString(userID), userID, nil)
		if err != nil {
			return emptyRes, fmt.Errorf("Failed to fetch item response: %v", err)
		}
		var chatsResponse models.GetChatsResponse
		var itemResponseBody map[string]interface{}
		err = json.Unmarshal(itemResponse.Value, &itemResponseBody)

		chats, ok := itemResponseBody["chats"].([]interface{})

		// retrieve chat items from chats list
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
				chatsResponse.Chats = append(chatsResponse.Chats, chatItem)
			}
		} else {
			return emptyRes, fmt.Errorf("Failed to parse chats: %v", err)
		}

		chatsResponse.Total = int(itemResponseBody["total"].(float64))
		chatsResponse.UserID = userID

		fmt.Println("Get chats response")

		return chatsResponse, nil
	}

	fmt.Println("User doesn't have any chats yet, return empty response")

	return emptyRes, nil
}

func GetUserChatsData(userID string) (data.UserData, error) {
	emptyRes := data.UserData{}

	container, err := ConnectDatabase("", "", "")
	if err != nil {
		return emptyRes, fmt.Errorf("Failed to connect to the database: %v", err)
	}

	itemResponse, err := container.ReadItem(context.Background(), azcosmos.NewPartitionKeyString(userID), userID, nil)
	if err != nil {
		return emptyRes, fmt.Errorf("Failed to fetch item response: %v", err)
	}

	var userDataResponse data.UserData
	var itemResponseBody map[string]interface{}
	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)

	chats, ok := itemResponseBody["chats"].([]interface{})

	// retrieve chat items from chats list
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
			userDataResponse.Chats = append(userDataResponse.Chats, chatItem)
		}
	} else {
		return emptyRes, fmt.Errorf("Failed to parse chats: %v", err)
	}

	userDataResponse.Total = int(itemResponseBody["total"].(float64))
	userDataResponse.UserID = userID
	userDataResponse.ID = userID
	userDataResponse.Partition = userID

	fmt.Println("Get chats response")

	return userDataResponse, nil
}
