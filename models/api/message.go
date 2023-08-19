package models

import "time"

type Message struct {
	Role         string    `json:"role"`
	Content      string    `json:"content"`
}

type MessageThread struct {
	MessageID    string    `json:"message_id"`
	Role         string    `json:"role"`
	CreateTime   time.Time `json:"create_time"`
	Content      string    `json:"content"`
	FinishReason string    `json:"finish_reason"`
	ParentID     string    `json:"parent_id"`
	Children     []string  `json:"children"`
}
