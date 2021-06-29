package entity

type MessageSent struct {
	Event string  `json:"event"`
	Data  Message `json:"data"`
}

func NewMessageSent(data Message) *MessageSent {
	return &MessageSent{
		Event: "message.sent",
		Data:  data,
	}
}
