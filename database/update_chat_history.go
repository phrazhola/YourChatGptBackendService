package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"YourChatGptBackendService/models/api"
	"YourChatGptBackendService/models/data"

	"github.com/google/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func UpdateChatHistory(completionMsg data.MessageThreadData, completionRequest models.SendMessageRequest) {
	container, err := ConnectDatabase("", "", "")
	if err != nil {
		_ = fmt.Errorf("Failed to connect to the database: %v", err)
	}

	chatData, err := GetChatData(completionRequest.UserID, completionRequest.ChatID, container)
	if err != nil {
		_ = fmt.Errorf("Failed to get chat data: %v", err)
	}

	timeNow := time.Now()
	chatHistory := chatData.ChatHistory

	// construct message data for prompt msg
	promptMsg := data.MessageThreadData{
		MessageID: uuid.New().String(),
		Role: completionRequest.Message.Role,
		// create time of propmt msg should be earlier than completion msg, set the delay to be 1 sec
		CreateTime: timeNow.Add(-time.Second),
		Content: completionRequest.Message.Content,
		ParentID: chatHistory.LastUpdatedMessage,
	}

	// construct message data for completion msg
	completionMsg.CreateTime = timeNow
	completionMsg.MessageID = uuid.New().String()
	completionMsg.ParentID = promptMsg.MessageID
	completionMsg.Children = []string{}

	promptMsg.Children = []string{completionMsg.MessageID}

	// update last updated msg
	if chatHistory.LastUpdatedMessage != "" {
		lastUpdatedMsg := chatHistory.MessagesMapping[chatHistory.LastUpdatedMessage]
		lastUpdatedMsg.Children = append(lastUpdatedMsg.Children, promptMsg.MessageID)
		chatHistory.MessagesMapping[chatHistory.LastUpdatedMessage] = lastUpdatedMsg
	}

	// add prompt msg and completion msg to the message mappings
	chatHistory.MessagesMapping[promptMsg.MessageID] = promptMsg
	chatHistory.MessagesMapping[completionMsg.MessageID] = completionMsg
	chatHistory.LastUpdatedMessage = completionMsg.MessageID

	chatData.UpdateTime = timeNow
	chatData.ChatHistory = chatHistory

	chatDataToUpdate, err := json.Marshal(chatData)
	if err != nil {
		_ = fmt.Errorf("Error: %v", err)
	}

	_, err = container.ReplaceItem(context.Background(), azcosmos.NewPartitionKeyString(completionRequest.UserID), completionRequest.ChatID, chatDataToUpdate, nil)
	if err != nil {
		_ = fmt.Errorf("Failed to replace item with id: %v", completionRequest.ChatID)
	}
}
