package entity

import "time"

type Message struct {
	ID        interface{} `json:"id"`
	Content   string      `json:"content"`
	SenderID  interface{} `json:"sender_id"`
	DiscussID interface{} `json:"discuss_id"`
	CreatedAt time.Time   `json:"created_at"`
}

type CreateMessage struct {
	Content   string      `json:"content"`
	SenderID  interface{} `json:"sender_id"`
	DiscussID interface{} `json:"discuss_id"`
	CreatedAt time.Time   `json:"created_at"`
}
