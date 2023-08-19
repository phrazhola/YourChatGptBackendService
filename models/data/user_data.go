package data

import (
	models "YourChatGptBackendService/models/api"
)

type UserData struct {
	ID        string        `json:"id"` // id is user_id value
	UserID    string        `json:"user_id"`
	Chats     []models.Chat `json:"chats"`
	Total     int           `json:"total"`
	Partition string        `json:"partition"` // partition is id value
}

type User struct {
	ID        string `json:"id"` // id value is encoded user email
	UserID    string `json:"user_id"`
	Partition string `json:"partition"` // partition value is the first letter of id
}
