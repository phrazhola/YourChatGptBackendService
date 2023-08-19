package models

import "time"

type SendMessageRequest struct {
	UserID  string  `json:"user_id"`
	ChatID  string  `json:"chat_id"`
	Message Message `json:"message"`
}

type SendMessageResponse struct {
	MessageID    string    `json:"message_id"`
	Role         string    `json:"role"`
	CreateTime   time.Time `json:"create_time"`
	Content      string    `json:"content"`
	FinishReason string    `json:"finish_reason"`
	ParentID     string    `json:"parent_id"`
	Children     []string  `json:"children"`
}
