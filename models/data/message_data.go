package data

import "time"

type MessageThreadData struct {
	MessageID    string    `json:"message_id"`
	Role         string    `json:"role"`
	CreateTime   time.Time `json:"create_time"`
	Content      string    `json:"content"`
	FinishReason string    `json:"finish_reason"`
	ParentID     string    `json:"parent_id"`
	Children     []string  `json:"children"`
}
