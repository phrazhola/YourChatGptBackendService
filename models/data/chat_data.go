package data

import "time"

type ChatData struct {
	ID          string      `json:"id"` // id is chat_id value
	ChatID      string      `json:"chat_id"`
	Title       string      `json:"title"`
	CreateTime  time.Time   `json:"create_time"`
	UpdateTime  time.Time   `json:"update_time"`
	ChatHistory ChatHistory `json:"chat_history"`
	BackRef     BackRef     `json:"back_ref"`
	Partition   string      `json:"partition"` // partition value is user_id
}

type BackRef struct {
	UserID string `json:"user_id"`
}

type ChatHistory struct {
	MessagesMapping    map[string]MessageThreadData `json:"messages_mapping"`
	LastUpdatedMessage string                       `json:"last_updated_msg"` // value is the message id
}
