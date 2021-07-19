package entity

import "time"

type MessageContentType string
type MessageReceiverType string

var (
	MessageContentText  MessageContentType = "message.content.text"
	MessageContentEvent MessageContentType = "message.content.event"
	MessageContentImage MessageContentType = "message.content.image"

	MessageReceiverUser       MessageReceiverType = "message.receiver.user"
	MessageReceiverDiscussion MessageReceiverType = "message.receiver.discussion"
)

type Message struct {
	ID           interface{}         `json:"id"`
	ContentType  MessageContentType  `json:"content_type"`
	Content      string              `json:"content"`
	ReceiverType MessageReceiverType `json:"receiver_type"`
	ReceiverID   interface{}         `json:"receiver_id"`
	Sender       User                `json:"sender"`
	CreatedAt    time.Time           `json:"created_at"`
}

type CreateMessage struct {
	ContentType  MessageContentType  `json:"content_type" bson:"contentType"`
	Content      string              `json:"content" bson:"content"`
	ReceiverType MessageReceiverType `json:"receiver_type" bson:"receiverType"`
	ReceiverID   interface{}         `json:"receiver_id" bson:"receiverId"`
	Sender       User                `json:"-"`
	CreatedAt    time.Time           `json:"-"`
}
