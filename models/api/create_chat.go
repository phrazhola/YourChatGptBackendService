package models

type CreateChatRequest struct {
	Email  string `json:"eeid"` // encoded email
	Title string `json:"title"`
}

type CreateChatResponse struct {
	UserID  string `json:"user_id"`
	ChatID  string `json:"chat_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
