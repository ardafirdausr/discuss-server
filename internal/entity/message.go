package entity

import "time"

type Message struct {
	ID           interface{} `json:"id" bson:"_id"`
	ContentType  string      `json:"content_type" bson:"contentType"`
	Content      string      `json:"content" bson:"content"`
	ReceiverType string      `json:"receiver_type" bson:"receiverType"`
	ReceiverID   interface{} `json:"receiver_id" bson:"receiverId"`
	Sender       User        `json:"sender" bson:"senderId"`
	CreatedAt    time.Time   `json:"created_at" bson:"createdAt"`
}

type CreateMessage struct {
	ContentType  string      `json:"content_type"`
	Content      string      `json:"content"`
	ReceiverType string      `json:"receiver_type"`
	ReceiverID   interface{} `json:"receiver_id"`
	Sender       User        `json:"-"`
	CreatedAt    time.Time   `json:"-"`
}
