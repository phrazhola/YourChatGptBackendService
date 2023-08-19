package models

type GetChatsRequest struct {
	UserID string `json:"user_id"`
}

type GetChatsResponse struct {
	UserID string `json:"user_id"`
	Chats  []Chat `json:"chats"`
	Total  int    `json:"total"`
}
