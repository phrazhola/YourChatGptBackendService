package database

import (
	"fmt"

	"YourChatGptBackendService/models/api"
)

func GetChatHistory(userID string, chatID string) (models.GetChatHistoryResponse, error) {
	container, err := ConnectDatabase("", "", "")
	if err != nil {
		return models.GetChatHistoryResponse{}, fmt.Errorf("Failed to connect to the database: %v", err)
	}

	chatData, err := GetChatData(userID, chatID, container)
	if err != nil {
		return models.GetChatHistoryResponse{}, fmt.Errorf("Failed to get chat data: %v", err)
	}

	msgMapping := make(map[string]models.MessageThread)
	mappingData := chatData.ChatHistory.MessagesMapping

	for key, value := range mappingData {
		msgMapping[key] = models.MessageThread{
			MessageID: value.MessageID,
			Role: value.Role,
			CreateTime: value.CreateTime,
			Content: value.Content,
			FinishReason: value.FinishReason,
			ParentID: value.ParentID,
			Children: value.Children,
		}
	}

	response := models.GetChatHistoryResponse{
		ChatID: chatID,
		Title: chatData.Title,
		CreateTime: chatData.CreateTime,
		UpdateTime: chatData.UpdateTime,
		MessagesMapping: msgMapping,
	}

	return response, nil
}
