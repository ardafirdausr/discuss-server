package usecase

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/entity"
)

type MessageUsecase struct {
	messageRepo internal.MessageRepository
	pubsub      *internal.PubSub
}

func NewMessageUsecase(messageRepo internal.MessageRepository) *MessageUsecase {
	return &MessageUsecase{messageRepo: messageRepo}
}

func (muc MessageUsecase) SendMessage(pubsub internal.PubSub, param entity.CreateMessage) (*entity.Message, error) {
	message, err := muc.messageRepo.Create(param)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	msgChannel := fmt.Sprintf("%s: %v", message.ReceiverType, message.ReceiverID)
	msgPayload, err := json.Marshal(message)
	if err == nil {
		pubsub.Publish(msgChannel, msgPayload)
		log.Println(err.Error())
	}

	return message, nil
}

func (muc MessageUsecase) ListenMessage(pubsub internal.PubSub) internal.SubscribeListener {
	return func(channel string, message interface{}) error {
		return nil
	}
}
