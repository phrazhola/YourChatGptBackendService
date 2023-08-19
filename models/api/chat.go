package models

import "time"

type Chat struct {
	ChatID      string    `json:"chat_id"`
	Title       string    `json:"title"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}