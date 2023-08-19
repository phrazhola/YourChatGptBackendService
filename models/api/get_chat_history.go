package models

import (
	"time"
)

type GetChatHistoryRequest struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

type GetChatHistoryResponse struct {
	ChatID          string                   `json:"chat_id"`
	Title           string                   `json:"title"`
	CreateTime      time.Time                `json:"create_time"`
	UpdateTime      time.Time                `json:"update_time"`
	MessagesMapping map[string]MessageThread `json:"messages_mapping"`
}
