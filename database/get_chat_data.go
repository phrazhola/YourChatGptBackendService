package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"YourChatGptBackendService/models/data"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func GetChatData(userID string, chatID string, container *azcosmos.ContainerClient) (data.ChatData, error) {
	emptyRes := data.ChatData{}

	itemResponse, err := container.ReadItem(context.Background(), azcosmos.NewPartitionKeyString(userID), chatID, nil)
	if err != nil {
		return emptyRes, fmt.Errorf("Failed to fetch item response: %v", err)
	}
	var chatData data.ChatData
	var chatHistory data.ChatHistory
	var itemResponseBody map[string]interface{}
	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		return emptyRes, fmt.Errorf("Failed to parse chat data: %v", err)
	}

	// retrieve chat history from response
	historyResponse, ok := itemResponseBody["chat_history"].(map[string]interface{})
	if !ok {
		return emptyRes, fmt.Errorf("Failed to parse chat history")
	}

	// retrieve messagesMapping from chat history
	msgMappingResponse := historyResponse["messages_mapping"].(map[string]interface{})
	msgMapping := make(map[string]data.MessageThreadData)

	if msgMappingResponse != nil && len(msgMappingResponse) > 0 {
		for key, value := range msgMappingResponse {
			msgItem := value.(map[string]interface{})
			childrenResponse := msgItem["children"].([]interface{})
			children := []string{}

			for _, child := range childrenResponse {
				children = append(children, child.(string))
			}

			msgData := data.MessageThreadData{
				MessageID: msgItem["message_id"].(string),
				Role: msgItem["role"].(string),
				CreateTime: func() time.Time {
					t, _ := time.Parse(time.RFC3339, msgItem["create_time"].(string))
					return t
				}(),
				Content: msgItem["content"].(string),
				FinishReason: msgItem["finish_reason"].(string),
				ParentID: msgItem["parent_id"].(string),
				Children: children,
			}

			msgMapping[key] = msgData
		}
	}

	chatHistory = data.ChatHistory{
		MessagesMapping: msgMapping,
		LastUpdatedMessage: historyResponse["last_updated_msg"].(string),
	}

	backRefResponse := itemResponseBody["back_ref"].(map[string]interface{})
	backRef := data.BackRef{
		UserID: backRefResponse["user_id"].(string),
	}

	chatData = data.ChatData{
		ID: itemResponseBody["id"].(string),
		ChatID: itemResponseBody["chat_id"].(string),
		Title: itemResponseBody["title"].(string),
		CreateTime: func() time.Time {
			t, _ := time.Parse(time.RFC3339, itemResponseBody["create_time"].(string))
			return t
		}(),
		UpdateTime: func() time.Time {
			t, _ := time.Parse(time.RFC3339, itemResponseBody["update_time"].(string))
			return t
		}(),
		ChatHistory: chatHistory,
		BackRef: backRef,
		Partition: itemResponseBody["partition"].(string),
	}

	return chatData, nil
}
