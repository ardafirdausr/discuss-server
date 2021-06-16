package entity

import "time"

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	SenderID  string    `json:"sender_id"`
	DiscussID string    `json:"discuss_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateMessage struct {
	Content   string    `json:"content"`
	SenderID  string    `json:"sender_id"`
	DiscussID string    `json:"discuss_id"`
	CreatedAt time.Time `json:"created_at"`
}
